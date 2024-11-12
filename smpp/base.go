package smpp

import (
	"io"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/smserror"
)

type base struct {
	Header
	OptionalParameters codec.OptionalFields
}

func newBase(commandId codec.CommandId, seqId int32) (v base) {
	v.OptionalParameters = make(codec.OptionalFields)
	v.CommandID = commandId
	if seqId > 0 {
		v.SequenceNumber = seqId
	} else {
		v.AssignSequenceNumber()
	}
	return
}

// GetHeader returns pdu header.
func (c *base) GetHeader() codec.Header {
	return &c.Header
}

func (c *base) unmarshal(b *codec.BytesReader, bodyReader func(*codec.BytesReader) error) (err error) {
	fullLen := b.Len()

	if err = c.Header.Unmarshal(b); err == nil {

		// try to unmarshal body
		if bodyReader != nil {
			err = bodyReader(b)
		}

		if err == nil {
			// command length
			cmdLength := int(c.CommandLength)

			// got - total read byte(s)
			got := fullLen - b.Len()
			if got > cmdLength {
				err = smserror.ErrInvalidPDU
				return
			}

			// body < command_length, still have optional parameters ?
			if got < cmdLength {
				optParam := b.ReadN(cmdLength - got)
				if err = b.Err(); err == nil {
					err = c.unmarshalOptionalParam(optParam)
				}
				if err != nil {
					return
				}
			}

			// validate again
			if b.Len() != fullLen-cmdLength {
				err = smserror.ErrInvalidPDU
			}
		}
	}

	return
}

func (c *base) unmarshalOptionalParam(optParam []byte) (err error) {
	buf := codec.ReaderPool.Get().(*codec.BytesReader)
	defer codec.ReaderPool.Put(buf)
	buf.Init(optParam)
	for buf.Len() > 0 {
		var field codec.Field
		if err = field.Unmarshal(buf); err == nil {
			c.OptionalParameters[field.Tag] = field
		} else {
			return
		}
	}
	return
}

// Marshal to buffer.
func (c *base) marshal(b *codec.BytesWriter, bodyWriter func(*codec.BytesWriter)) {

	bodyBuf := codec.WriterPool.Get().(*codec.BytesWriter)
	defer codec.WriterPool.Put(bodyBuf)
	bodyBuf.Reset()
	// body
	if bodyWriter != nil {
		bodyWriter(bodyBuf)
	}

	// optional body
	for _, v := range c.OptionalParameters {
		v.Marshal(bodyBuf)
	}

	// write header
	c.CommandLength = uint32(PDU_HEADER_SIZE + bodyBuf.Len())
	c.Header.Marshal(b)

	// write body and its optional params
	b.WriteBytes(bodyBuf.Bytes())
}

// RegisterOptionalParam register optional param.
func (c *base) RegisterOptionalParam(tlv codec.Field) {
	c.OptionalParameters[tlv.Tag] = tlv
}

// IsOk is status ok.
func (c *base) IsOk() bool {
	return c.CommandStatus == ESME_ROK
}

// IsGNack is generic n-ack.
func (c *base) IsGNack() bool {
	return c.CommandID == GENERIC_NACK
}

// Parse PDU from reader.
func Parse(r io.Reader, logger *logrus.Entry) (pdu codec.PDU, err error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorln("smpp.parse.err", err)
			err = smserror.ErrInvalidPDU
		}
	}()
	var headerBytes [16]byte

	if _, err = io.ReadFull(r, headerBytes[:]); err != nil {
		return
	}

	header := ParseHeader(headerBytes)
	if header.CommandLength < 16 || header.CommandLength > MAX_PDU_LEN {
		err = smserror.ErrInvalidPDU
		return
	}

	// read pdu body
	bodyBytes := make([]byte, header.CommandLength-16)
	if len(bodyBytes) > 0 {
		if _, err = io.ReadFull(r, bodyBytes); err != nil {
			return
		}
	}
	if logger != nil {
		switch header.CommandID {
		case ENQUIRE_LINK, ENQUIRE_LINK_RESP:
		default:
			logger.Infof("recv[%s]%x%x", header, headerBytes, bodyBytes)
		}
	}
	// try to create pdu
	if pdu, err = CreatePDUFromCmdID(header.CommandID); err == nil {

		buf := codec.WriterPool.Get().(*codec.BytesWriter)
		defer codec.WriterPool.Put(buf)
		buf.Init(headerBytes[:])
		if len(bodyBytes) > 0 {
			buf.Write(bodyBytes)
		}
		reader := codec.ReaderPool.Get().(*codec.BytesReader)
		defer codec.ReaderPool.Put(reader)
		reader.Init(buf.Bytes())
		err = pdu.Unmarshal(reader)
	} else {
		logrus.Infof("read.CreatePDUFromCmdID %d,%v", header.CommandID, err)
	}
	return
}
