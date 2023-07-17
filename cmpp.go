package zysms

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/smserror"
)

type cmppConn struct {
	net.Conn
	State enum.State
	Typ   codec.Version
	// for SeqId generator goroutine
	// SeqId  <-chan uint32
	// done   chan<- struct{}
	_seqId   uint32
	stop     func()
	ctx      context.Context
	counter  int32
	logger   *logrus.Entry
	checkVer bool
}

// New returns an abstract structure for successfully
// established underlying net.Conn.
func newCmppConn(conn net.Conn, typ codec.Version, checkVer bool) *Conn {
	c := &cmppConn{
		Conn:     conn,
		Typ:      typ,
		_seqId:   0,
		logger:   logrus.WithFields(logrus.Fields{"r": conn.RemoteAddr()}),
		checkVer: checkVer,
	}
	c.ctx, c.stop = context.WithCancel(context.Background())

	tc := c.Conn.(*net.TCPConn)
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(1 * time.Minute) // 1min
	return &Conn{smsConn: c, UUID: uuid.New().String()}
}

func (c *cmppConn) Auth(uid string, pwd string) error {
	// Login to the server.
	req := cmpp.NewConnReq(c.Typ).(*cmpp.ConnReq)
	req.SrcAddr = uid
	req.Secret = pwd
	req.Version = c.Typ

	err := c.SendPDU(req)
	if err != nil {
		c.logger.Errorf("cmpp.auth send error: %v", err)
		return err
	}
	p, err := c.RecvPDU()
	if err != nil {
		c.logger.Errorf("cmpp.auth recv error: %v", err)
		return err
	}
	var status uint8

	if rsp, ok := p.(*cmpp.ConnResp); ok {
		if c.checkVer && rsp.Version != c.Typ {
			return smserror.ErrVersionNotMatch
		}
		status = uint8(rsp.Status)
	} else {
		return smserror.ErrRespNotMatch
	}

	if status != 0 {
		if status <= cmpp.ErrnoConnOthers { //ErrnoConnOthers = 5
			err = cmpp.ConnRspStatusErrMap[status]
		} else {
			err = cmpp.ConnRspStatusErrMap[cmpp.ErrnoConnOthers]
		}
		return err
	}
	c.SetState(enum.CONN_AUTHOK)
	return nil
}
func (c *cmppConn) Close() {
	if c != nil {
		if c.State == enum.CONN_CLOSED {
			return
		}
		c.Conn.Close() // close the underlying net.Conn
		c.State = enum.CONN_CLOSED
		c.stop()
	}
}

func (c *cmppConn) SetState(state enum.State) {
	c.State = state
}
func (c *cmppConn) seqId() uint32 {
	sid := atomic.AddUint32(&c._seqId, 1)
	return sid
}

// SendPkt pack the cmpp packet structure and send it to the other peer.
func (c *cmppConn) SendPDU(pdu codec.PDU) error {
	if c.State == enum.CONN_CLOSED {
		return smserror.ErrConnIsClosed
	}
	if pdu == nil {
		return smserror.ErrPktIsNil
	}
	c.Logger().Infof("send pkt:%T , %s", pdu, c.Typ)
	if p, ok := pdu.(*cmpp.SubmitReq); ok {
		multiMsg, _ := p.Message.Split()
		p.TpUdhi = 0
		if len(multiMsg) > 1 {
			p.TpUdhi = 1
		}
		for _, msg := range multiMsg {
			// p.MsgLength = byte(len(content))
			// p.MsgContent = content
			p.Message = *msg
			p.AssignSequenceNumber()
			buf := codec.NewWriter()
			pdu.Marshal(buf)
			_, err := c.Conn.Write(buf.Bytes()) //block write
			if err != nil {
				return err
			}
		}
	} else {
		buf := codec.NewWriter()
		pdu.Marshal(buf)
		_, err := c.Conn.Write(buf.Bytes()) //block write
		if err != nil {
			return err
		}
	}
	return nil
}

// func (c *cmppConn) splitSubmitContent(req *cmpp.CmppSubmitReq) [][]byte {
// 	cLen := 140
// 	if req.MsgFmt == 0 {
// 		cLen = 160
// 	}
// 	cLen -= 6 // 减去7字节的消息头
// 	if len(req.MsgContent) <= cLen {
// 		return [][]byte{req.MsgContent}
// 	}
// 	count := len(req.MsgContent) / cLen
// 	if len(req.MsgContent)%cLen > 0 {
// 		count++
// 	}
// 	contentList := make([][]byte, count)
// 	idx := uint16(time.Now().UnixMilli() % 0xff)
// 	// 0x06 数据头长度
// 	// 0x00 信息标识
// 	// 0x04 后续信息头长度
// 	// 0x00,0x00 信息序列号
// 	// 0x00 总条数
// 	// 0x01 当前条数
// 	dhi := []byte{0x05, 0x00, 0x04, byte(idx), byte(count), 0x01}
// 	for i := 0; i < count; i++ {
// 		dhi[5] = byte(i + 1)
// 		if i == count-1 {
// 			contentList[i] = append(dhi, req.MsgContent[i*cLen:]...)
// 		} else {
// 			contentList[i] = append(dhi, req.MsgContent[i*cLen:(i+1)*cLen]...)
// 		}
// 	}
// 	return contentList
// }

func (c *cmppConn) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
func (c *cmppConn) Logger() *logrus.Entry {
	return c.logger
}
func (c *cmppConn) SetReadDeadline(timeout time.Duration) {
	if timeout > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(timeout))
	}
}

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (c *cmppConn) RecvPDU() (codec.PDU, error) {
	if c.State == enum.CONN_CLOSED {
		return nil, smserror.ErrConnIsClosed
	}

	pdu, err := cmpp.Parse(c.Conn, c.Typ)
	if err != nil {
		return nil, err
	}

	switch p := pdu.(type) {
	case *cmpp.ActiveTestReq: // 当收到心跳请求,内部直接回复心跳,并递归继续获取数据
		resp := p.GetResponse()
		c.SendPDU(resp)
		return c.RecvPDU()
	case *cmpp.ActiveTestResp: // 当收到心跳回复,内部直接处理,并递归继续获取数据
		atomic.AddInt32(&c.counter, -1)
		return c.RecvPDU()
	case *cmpp.ConnResp: // 当收到登录回复,内部先校验版本
		if c.checkVer && p.Version != c.Typ {
			return nil, fmt.Errorf("cmpp version not match [ local: %d != remote: %d ]", c.Typ, p.Version)
		}
	case *cmpp.ConnReq: // 当收到登录回复,内部先校验版本
		// 服务端自适应版本
		c.Typ = p.Version
		c.logger = logrus.WithFields(logrus.Fields{"r": c.RemoteAddr(), "v": c.Typ})
	}
	return pdu, nil

}

type cmppListener struct {
	net.Listener
}

func newCmppListener(l net.Listener) *cmppListener {
	return &cmppListener{l}
}

func (l *cmppListener) accept() (*Conn, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	conn := newCmppConn(c, cmpp.V30, false)
	conn.SetState(enum.CONN_CONNECTED)
	conn.smsConn.(*cmppConn).startActiveTest()
	return conn, nil
}

func (c *cmppConn) startActiveTest() {
	go func() {
		t := time.NewTicker(10 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-c.ctx.Done():
				// once conn close, the goroutine should exit
				return
			case <-t.C:
				// send a active test packet to peer, increase the active test counter
				p := cmpp.NewActiveTestReq(c.Typ)
				err := c.SendPDU(p)
				if err != nil {
					c.logger.Errorf("cmpp.active send error: %v", err)
				} else {
					atomic.AddInt32(&c.counter, 1)
				}
			}
		}
	}()
}
func (l *cmppListener) Close() error {
	return l.Listener.Close()
}
