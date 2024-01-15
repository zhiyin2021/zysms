package codec

import (
	"fmt"

	"github.com/zhiyin2021/zysms/codec/gsm7"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
)

var (
	// ErrNotImplSplitterInterface indicates that encoding does not support Splitter interface
	ErrNotImplSplitterInterface = fmt.Errorf("encoding not implementing Splitter interface")
	// ErrNotImplDecode indicates that encoding does not support Decode method
	ErrNotImplDecode = fmt.Errorf("decode is not implemented in this Encoding")
	// ErrNotImplEncode indicates that encoding does not support Encode method
	ErrNotImplEncode = fmt.Errorf("encode is not implemented in this Encoding")

	ErrAsciiInvalidCharacter = fmt.Errorf("invalid ascii character")

	// ErrInvalidByte means that a given byte is outside of the GSM 7-bit encoding range.
	//
	// This can only happen during decoding.
	ErrAsciiInvalidByte = fmt.Errorf("invalid ascii byte")
)

const (
	// GSM7BITCoding is gsm-7bit coding
	GSM7BITCoding byte = 0x00
	// ASCIICoding is ascii coding
	ASCIICoding byte = 0x01
	// BINARY8BIT1Coding is 8-bit binary coding
	BINARY8BIT1Coding byte = 0x02
	// LATIN1Coding is iso-8859-1 coding
	LATIN1Coding byte = 0x03
	// BINARY8BIT2Coding is 8-bit binary coding
	BINARY8BIT2Coding byte = 0x04
	// CYRILLICCoding is iso-8859-5 coding
	CYRILLICCoding byte = 0x06
	// HEBREWCoding is iso-8859-8 coding
	HEBREWCoding byte = 0x07
	// UCS2Coding is UCS2 coding
	UCS2Coding    byte = 0x08
	GB18030Coding byte = 0x0F
)

// EncDec wraps encoder and decoder interface.
type EncDec interface {
	Encode(str string) ([]byte, error)
	Decode([]byte) (string, error)
}

// Encoding interface.
type Encoding interface {
	EncDec
	DataCoding() byte
	EncodeSplit(text string) ([][]byte, error)
}

func encode(str string, encoder *encoding.Encoder) ([]byte, error) {
	return encoder.Bytes([]byte(str))
}

func decode(data []byte, decoder *encoding.Decoder) (st string, err error) {
	tmp, err := decoder.Bytes(data)
	if err == nil {
		st = string(tmp)
	}
	return
}

type gsm7bit struct {
	packed bool
	lang   gsm7.Lang
}

func (c *gsm7bit) Encode(str string) ([]byte, error) {
	return encode(str, gsm7.GSM7(c.packed, c.lang).NewEncoder())
}

func (c *gsm7bit) Decode(data []byte) (string, error) {
	return decode(data, gsm7.GSM7(c.packed, c.lang).NewDecoder())
}

func (c *gsm7bit) DataCoding() byte { return GSM7BITCoding }

func (c *gsm7bit) EncodeSplit(text string) (allSeg [][]byte, err error) {

	tmp, err := c.Encode(text)
	if err != nil {
		return nil, err
	}
	octetLimit := SM_GSM_LONG_LEN
	if !c.packed {
		octetLimit = SM_GSM_LONG_PACKLEN
		if len(tmp) <= SM_GSM_MSG_PACKLEN {
			return [][]byte{tmp}, nil
		}
	} else {
		if len(tmp) <= SM_GSM_MSG_LEN {
			return [][]byte{tmp}, nil
		}
	}

	allSeg = [][]byte{}
	runeSlice := tmp
	fr, to := 0, int(octetLimit)

	for fr < len(runeSlice) {
		if to > len(runeSlice) {
			to = len(runeSlice)
		}
		// seg, err := c.Encode(string(runeSlice[fr:to]))
		// if err != nil {
		// 	return nil, err
		// }
		seg := runeSlice[fr:to]
		allSeg = append(allSeg, seg)
		fr, to = to, to+int(octetLimit)
	}

	return
}

type ascii struct{}

func (*ascii) Encode(str string) ([]byte, error) {
	if HasWidthChar(str) {
		return nil, ErrAsciiInvalidCharacter
	}
	return []byte(str), nil
}

func (*ascii) Decode(data []byte) (string, error) {
	if HasWidthChar(string(data)) {
		return "nil", ErrAsciiInvalidByte
	}
	return string(data), nil
}

func (*ascii) DataCoding() byte { return ASCIICoding }

func (c *ascii) EncodeSplit(text string) (allSeg [][]byte, err error) {
	octetLimit := SM_GSM_LONG_PACKLEN

	allSeg = [][]byte{}
	runeSlice, err := c.Encode(text)
	if err != nil {
		return nil, err
	}
	if len(runeSlice) <= SM_GSM_MSG_PACKLEN {
		return [][]byte{runeSlice}, nil
	}
	fr, to := 0, octetLimit
	for fr < len(runeSlice) {
		if to > len(runeSlice) {
			to = len(runeSlice)
		}
		seg, err := c.Encode(string(runeSlice[fr:to]))
		if err != nil {
			return nil, err
		}
		allSeg = append(allSeg, seg)
		fr, to = to, to+int(octetLimit)
	}
	return
}

type iso88591 struct{}

func (*iso88591) Encode(str string) ([]byte, error) {
	return encode(str, charmap.ISO8859_1.NewEncoder())
}

func (*iso88591) Decode(data []byte) (string, error) {
	return decode(data, charmap.ISO8859_1.NewDecoder())
}

func (*iso88591) DataCoding() byte { return LATIN1Coding }

func (c *iso88591) EncodeSplit(text string) (allSeg [][]byte, err error) {
	octetLimit := SM_GSM_LONG_LEN

	allSeg = [][]byte{}
	runeSlice, err := c.Encode(text)
	if err != nil {
		return nil, err
	}
	if len(runeSlice) <= SM_GSM_MSG_LEN {
		return [][]byte{runeSlice}, nil
	}

	fr, to := 0, octetLimit
	for fr < len(runeSlice) {
		if to > len(runeSlice) {
			to = len(runeSlice)
		}
		seg, err := c.Encode(string(runeSlice[fr:to]))
		if err != nil {
			return nil, err
		}
		allSeg = append(allSeg, seg)
		fr, to = to, to+int(octetLimit)
	}
	return
}

type binary8bit1 struct{}

func (*binary8bit1) Encode(msg string) ([]byte, error) {
	return []byte(msg), nil
}

func (*binary8bit1) Decode(msg []byte) (string, error) {
	return string(msg), nil
}

func (*binary8bit1) DataCoding() byte { return BINARY8BIT1Coding }

func (c *binary8bit1) EncodeSplit(text string) (allSeg [][]byte, err error) {

	octetLimit := SM_GSM_LONG_LEN
	allSeg = [][]byte{}
	runeSlice, err := c.Encode(text)
	if err != nil {
		return nil, err
	}
	if len(runeSlice) <= SM_GSM_MSG_LEN {
		return [][]byte{runeSlice}, nil
	}

	fr, to := 0, int(octetLimit)
	for fr < len(runeSlice) {
		if to > len(runeSlice) {
			to = len(runeSlice)
		}
		seg, err := c.Encode(string(runeSlice[fr:to]))
		if err != nil {
			return nil, err
		}
		allSeg = append(allSeg, seg)
		fr, to = to, to+int(octetLimit)
	}

	return
}

type binary8bit2 struct{}

func (*binary8bit2) Encode(msg string) ([]byte, error) {
	return []byte(msg), nil
}

func (*binary8bit2) Decode(msg []byte) (string, error) {
	return string(msg), nil
}

func (*binary8bit2) DataCoding() byte { return BINARY8BIT2Coding }

func (c *binary8bit2) EncodeSplit(text string) (allSeg [][]byte, err error) {
	octetLimit := SM_GSM_LONG_LEN
	allSeg = [][]byte{}
	runeSlice, err := c.Encode(text)
	if err != nil {
		return nil, err
	}
	if len(runeSlice) <= SM_GSM_MSG_LEN {
		return [][]byte{runeSlice}, nil
	}

	fr, to := 0, int(octetLimit)
	for fr < len(runeSlice) {
		if to > len(runeSlice) {
			to = len(runeSlice)
		}
		seg, err := c.Encode(string(runeSlice[fr:to]))
		if err != nil {
			return nil, err
		}
		allSeg = append(allSeg, seg)
		fr, to = to, to+int(octetLimit)
	}
	return
}

type iso88595 struct{}

func (*iso88595) Encode(str string) ([]byte, error) {
	return encode(str, charmap.ISO8859_5.NewEncoder())
}

func (*iso88595) Decode(data []byte) (string, error) {
	return decode(data, charmap.ISO8859_5.NewDecoder())
}

func (*iso88595) DataCoding() byte { return CYRILLICCoding }

func (c *iso88595) EncodeSplit(text string) (allSeg [][]byte, err error) {
	octetLimit := SM_GSM_LONG_LEN
	allSeg = [][]byte{}
	runeSlice, err := c.Encode(text)
	if err != nil {
		return nil, err
	}
	if len(runeSlice) <= SM_GSM_MSG_LEN {
		return [][]byte{runeSlice}, nil
	}

	fr, to := 0, int(octetLimit)
	for fr < len(runeSlice) {
		if to > len(runeSlice) {
			to = len(runeSlice)
		}
		seg, err := c.Encode(string(runeSlice[fr:to]))
		if err != nil {
			return nil, err
		}
		allSeg = append(allSeg, seg)
		fr, to = to, to+int(octetLimit)
	}
	return
}

type iso88598 struct{}

func (*iso88598) Encode(str string) ([]byte, error) {
	return encode(str, charmap.ISO8859_8.NewEncoder())
}

func (*iso88598) Decode(data []byte) (string, error) {
	return decode(data, charmap.ISO8859_8.NewDecoder())
}

func (c *iso88598) EncodeSplit(text string) (allSeg [][]byte, err error) {
	octetLimit := SM_GSM_LONG_LEN
	allSeg = [][]byte{}
	runeSlice, _ := c.Encode(text)
	if len(runeSlice) <= SM_GSM_MSG_LEN {
		return [][]byte{[]byte(runeSlice)}, nil
	}

	fr, to := 0, int(octetLimit)
	for fr < len(runeSlice) {
		if to > len(runeSlice) {
			to = len(runeSlice)
		}
		seg, err := c.Encode(string(runeSlice[fr:to]))
		if err != nil {
			return nil, err
		}
		allSeg = append(allSeg, seg)
		fr, to = to, to+int(octetLimit)
	}
	return
}

func (*iso88598) DataCoding() byte { return HEBREWCoding }

type ucs2 struct{}

func (*ucs2) Encode(str string) ([]byte, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	return encode(str, tmp.NewEncoder())
}

func (*ucs2) Decode(data []byte) (string, error) {
	tmp := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	return decode(data, tmp.NewDecoder())
}

func (c *ucs2) EncodeSplit(text string) (allSeg [][]byte, err error) {
	allSeg = [][]byte{}
	if text == "" {
		return allSeg, nil
	}
	octetLimit := SM_GSM_LONG_LEN
	allSeg = [][]byte{}
	runeSlice, err := c.Encode(text)
	if err != nil {
		return nil, err
	}
	if len(runeSlice) <= SM_GSM_MSG_LEN {
		return [][]byte{runeSlice}, nil
	}

	size := len(runeSlice)
	count := size / octetLimit
	if size%octetLimit > 0 {
		count++
	}
	for i := 0; i < count; i++ {
		if (i+1)*octetLimit > size {
			allSeg = append(allSeg, runeSlice[i*octetLimit:])
		} else {
			allSeg = append(allSeg, runeSlice[i*octetLimit:(i+1)*octetLimit])
		}
	}

	return
}

func (*ucs2) DataCoding() byte { return UCS2Coding }

type gb18030 struct{}

func (*gb18030) Encode(str string) ([]byte, error) {
	return encode(str, simplifiedchinese.GB18030.NewEncoder())
}

func (*gb18030) Decode(data []byte) (string, error) {
	return decode(data, simplifiedchinese.GB18030.NewDecoder())
}

func (c *gb18030) EncodeSplit(text string) (allSeg [][]byte, err error) {

	octetLimit := SM_GSM_LONG_LEN
	allSeg = [][]byte{}
	runeSlice, _ := c.Encode(text)
	if len(runeSlice) <= SM_GSM_MSG_LEN {
		return [][]byte{runeSlice}, nil
	}

	hextetLim := int(octetLimit) // round down
	// hextet = 16 bits, the correct terms should be hexadectet
	fr, to := 0, 0
	i := 0
	count := len(runeSlice)
	for fr < count {
		to = 0
		for to < hextetLim && fr+to < count {
			if runeSlice[fr+to] > 0x7F {
				to += 2
				if to > hextetLim {
					to -= 2
					break
				}
			} else {
				to++
			}
		}
		to = fr + to
		if to > len(runeSlice) {
			to = len(runeSlice)
		}
		// seg, err := c.Encode(string(runeSlice[fr:to]))
		// if err != nil {
		// 	return nil, err
		// }
		seg := runeSlice[fr:to]
		allSeg = append(allSeg, seg)
		fr = to
		i++
		if i > 10 {
			return
		}
	}
	return
}

// 测试一段长短信内容,http://www.baidu.com,
func (*gb18030) DataCoding() byte { return GB18030Coding }

var (
	// GSM7BIT is gsm-7bit encoding.
	GSM7BIT Encoding = &gsm7bit{packed: false}
	// 西班牙
	GSM7Spanish Encoding = &gsm7bit{packed: false, lang: gsm7.LangSpanish}
	// 葡萄牙
	GSM7Portuguese Encoding = &gsm7bit{packed: false, lang: gsm7.LangPortuguese}
	// 土耳其
	GSM7Turkish Encoding = &gsm7bit{packed: false, lang: gsm7.LangTurkish}

	// GSM7BITPACKED is packed gsm-7bit encoding.
	// Most of SMSC(s) use unpack version.
	// Should be tested before using.
	GSM7BITPACKED Encoding = &gsm7bit{packed: true}

	// ASCII is ascii encoding.
	ASCII Encoding = &ascii{}

	// BINARY8BIT1 is binary 8-bit encoding.
	BINARY8BIT1 Encoding = &binary8bit1{}

	// LATIN1 encoding.
	LATIN1 Encoding = &iso88591{}

	// BINARY8BIT2 is binary 8-bit encoding.
	BINARY8BIT2 Encoding = &binary8bit2{}

	// CYRILLIC encoding.
	CYRILLIC Encoding = &iso88595{}

	// HEBREW encoding.
	HEBREW Encoding = &iso88598{}

	// UCS2 encoding.
	UCS2    Encoding = &ucs2{}
	GB18030 Encoding = &gb18030{}
)

var codingMap = map[byte]Encoding{
	GSM7BITCoding:     GSM7BIT,
	ASCIICoding:       ASCII,
	BINARY8BIT1Coding: BINARY8BIT1,
	LATIN1Coding:      LATIN1,
	BINARY8BIT2Coding: BINARY8BIT2,
	CYRILLICCoding:    CYRILLIC,
	HEBREWCoding:      HEBREW,
	UCS2Coding:        UCS2,
	GB18030Coding:     GB18030,
}

const (
	// 闪信
	FlashMsg = 0x10
)

// FromDataCoding returns encoding from DataCoding value.
func GetCodec(code byte) (enc Encoding) {
	if code&FlashMsg == FlashMsg {
		code = code ^ FlashMsg
	}
	enc = codingMap[code]
	if enc == nil {
		enc = UCS2
	}
	return
}

// Splitter extend encoding object by defining a split function
// that split a string into multiple segments
// Each segment string, when encoded, must be within a certain octet limit
// type Splitter interface {
// 	// ShouldSplit check if the encoded data of given text should be splitted under octetLimit
// 	ShouldSplit(text string, octetLimit uint) (should bool)
// 	EncodeSplit(text string, octetLimit uint) ([][]byte, error)
// }

// 判断字符串是否包含中文
func HasWidthChar(content string) bool {
	if content == "" {
		return false
	}
	for _, c := range content {
		if c > 0x7f {
			return true
		}
	}
	return false
}
