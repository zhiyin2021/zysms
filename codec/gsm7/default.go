package gsm7

import (
	"bytes"
	"errors"
	"fmt"
	"math"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

type Lang byte

// UDH contains 0x24 0x01 {lang} //扩展表
// UDH contains 0x25 0x01 {lang} //主表

const (
	LangSpanish    Lang = 0x02 //"spanish" //西班牙
	LangPortuguese Lang = 0x03 // "portuguese" //葡萄牙
	LangTurkish    Lang = 0x01 //"turkish"    //土耳其
)

var (
	// ErrInvalidCharacter means a given character can not be represented in GSM 7-bit encoding.
	//
	// This can only happen during encoding.
	ErrInvalidCharacter = errors.New("invalid gsm7 character")

	// ErrInvalidByte means that a given byte is outside of the GSM 7-bit encoding range.
	//
	// This can only happen during decoding.
	ErrInvalidByte = errors.New("invalid gsm7 byte")

	zr0            = rune(0)
	esc            = rune(0x1B)
	escapeSequence = byte(0x1B)
	gsmDefault     = []rune{ // 7-bit GSM charset
		'@', '£', '$', '¥', 'è', 'é', 'ù', 'ì', 'ò', 'Ç', '\n', 'Ø', 'ø', '\r', 'Å', 'å',
		'Δ', '_', 'Φ', 'Γ', 'Λ', 'Ω', 'Π', 'Ψ', 'Σ', 'Θ', 'Ξ', esc, 'Æ', 'æ', 'ß', 'É',
		' ', '!', '"', '#', '¤', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
		'¡', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'Ä', 'Ö', 'Ñ', 'Ü', '§',
		'¿', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'ä', 'ö', 'ñ', 'ü', 'à',
	}

	gsmDefaultExt = []rune{ // 7-bit GSM charset
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, '\f', zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, zr0, '^', zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, '{', '}', zr0, zr0, zr0, zr0, zr0, '\\',
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, '[', '~', ']', zr0,
		'|', zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, zr0, zr0, '€', zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
		zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0, zr0,
	}

	gsmDeDefault    map[byte]rune
	gsmDeDefaultExt map[byte]rune

	gsmEnDefault    map[rune]byte
	gsmEnDefaultExt map[rune]byte
)

func errInvalidByte(n int, buf []byte) error {
	return fmt.Errorf("invalid gsm7 byte at %d, %x", n, buf)
}
func init() {
	gsmDeDefault = make(map[byte]rune)
	gsmDeDefaultExt = make(map[byte]rune)
	gsmEnDefault = make(map[rune]byte)
	gsmEnDefaultExt = make(map[rune]byte)
	for i, v := range gsmDefault {
		gsmDeDefault[byte(i)] = v
		gsmEnDefault[v] = byte(i)
	}
	for i, v := range gsmDefaultExt {
		if v != zr0 {
			gsmDeDefaultExt[byte(i)] = v
			gsmEnDefaultExt[v] = byte(i)
		}
	}
}

// GSM7 returns a GSM 7-bit Bit Encoding.
//
// Set the packed flag to true if you wish to convert septets to octets,
// this should be false for most SMPP providers.
func GSM7(packed bool, lang Lang) encoding.Encoding {
	return gsm7Encoding{packed: packed, lang: lang}
}

type gsm7Encoding struct {
	lang   Lang
	packed bool
}

func (g gsm7Encoding) NewDecoder() *encoding.Decoder {
	data, dateExt := gsmDeDefault, gsmDeDefaultExt
	switch g.lang {
	case LangSpanish:
		data, dateExt = gsmDeSpanish, gsmDeSpanishExt
	case LangPortuguese:
		data, dateExt = gsmDePortuguese, gsmDePortugueseExt
	case LangTurkish:
		data, dateExt = gsmDeTurkish, gsmDeTurkishExt
	}
	return &encoding.Decoder{Transformer: &gsm7Decoder{
		data:    data,
		dataExt: dateExt,
		packed:  g.packed,
	}}
}

func (g gsm7Encoding) NewEncoder() *encoding.Encoder {
	data, dateExt := gsmEnDefault, gsmEnDefaultExt
	switch g.lang {
	case LangSpanish:
		data, dateExt = gsmEnSpanish, gsmEnSpanishExt
	case LangPortuguese:
		data, dateExt = gsmEnPortuguese, gsmEnPortugueseExt
	case LangTurkish:
		data, dateExt = gsmEnTurkish, gsmEnTurkishExt
	}
	return &encoding.Encoder{Transformer: &gsm7Encoder{
		data:    data,
		dataExt: dateExt,
		packed:  g.packed,
	}}
}

func (g gsm7Encoding) String() string {
	if g.packed {
		return "GSM 7-bit (Packed)"
	}
	return "GSM 7-bit (Unpacked)"
}

type gsm7Decoder struct {
	packed  bool
	data    map[byte]rune
	dataExt map[byte]rune
}

func (g *gsm7Decoder) Reset() { /* not needed */ }

func unpack(src []byte, packed bool) (septets []byte) {
	septets = src
	if packed {
		septets = make([]byte, 0, len(src))
		count := 0
		for remain := len(src) - count; remain > 0; {
			// Unpack by converting octets into septets.
			switch {
			case remain >= 7:
				septets = append(septets, src[count+0]&0x7F<<0)
				septets = append(septets, (src[count+1]&0x3F<<1)|(src[count+0]&0x80>>7))
				septets = append(septets, (src[count+2]&0x1F<<2)|(src[count+1]&0xC0>>6))
				septets = append(septets, (src[count+3]&0x0F<<3)|(src[count+2]&0xE0>>5))
				septets = append(septets, (src[count+4]&0x07<<4)|(src[count+3]&0xF0>>4))
				septets = append(septets, (src[count+5]&0x03<<5)|(src[count+4]&0xF8>>3))
				septets = append(septets, (src[count+6]&0x01<<6)|(src[count+5]&0xFC>>2))
				if src[count+6] > 0 {
					septets = append(septets, src[count+6]&0xFE>>1)
				}
				count += 7
			case remain >= 6:
				septets = append(septets, src[count+0]&0x7F<<0)
				septets = append(septets, (src[count+1]&0x3F<<1)|(src[count+0]&0x80>>7))
				septets = append(septets, (src[count+2]&0x1F<<2)|(src[count+1]&0xC0>>6))
				septets = append(septets, (src[count+3]&0x0F<<3)|(src[count+2]&0xE0>>5))
				septets = append(septets, (src[count+4]&0x07<<4)|(src[count+3]&0xF0>>4))
				septets = append(septets, (src[count+5]&0x03<<5)|(src[count+4]&0xF8>>3))
				count += 6
			case remain >= 5:
				septets = append(septets, src[count+0]&0x7F<<0)
				septets = append(septets, (src[count+1]&0x3F<<1)|(src[count+0]&0x80>>7))
				septets = append(septets, (src[count+2]&0x1F<<2)|(src[count+1]&0xC0>>6))
				septets = append(septets, (src[count+3]&0x0F<<3)|(src[count+2]&0xE0>>5))
				septets = append(septets, (src[count+4]&0x07<<4)|(src[count+3]&0xF0>>4))
				count += 5
			case remain >= 4:
				septets = append(septets, src[count+0]&0x7F<<0)
				septets = append(septets, (src[count+1]&0x3F<<1)|(src[count+0]&0x80>>7))
				septets = append(septets, (src[count+2]&0x1F<<2)|(src[count+1]&0xC0>>6))
				septets = append(septets, (src[count+3]&0x0F<<3)|(src[count+2]&0xE0>>5))
				count += 4
			case remain >= 3:
				septets = append(septets, src[count+0]&0x7F<<0)
				septets = append(septets, (src[count+1]&0x3F<<1)|(src[count+0]&0x80>>7))
				septets = append(septets, (src[count+2]&0x1F<<2)|(src[count+1]&0xC0>>6))
				count += 3
			case remain >= 2:
				septets = append(septets, src[count+0]&0x7F<<0)
				septets = append(septets, (src[count+1]&0x3F<<1)|(src[count+0]&0x80>>7))
				count += 2
			case remain >= 1:
				septets = append(septets, src[count+0]&0x7F<<0)
				count++
			default:
				return
			}
			remain = len(src) - count
		}
	}
	return
}

func (g *gsm7Decoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if len(src) == 0 {
		return 0, 0, nil
	}

	septets := unpack(src, g.packed)

	nSeptet := 0
	builder := bytes.NewBufferString("")
	for nSeptet < len(septets) {
		b := septets[nSeptet]
		if b == escapeSequence {
			nSeptet++
			if nSeptet >= len(septets) {
				err = errInvalidByte(nSeptet, septets)
				continue
			}
			e := septets[nSeptet]
			if r, ok := g.dataExt[e]; ok {
				builder.WriteRune(r)
			} else {
				err = errInvalidByte(nSeptet, septets)
				continue
			}
		} else if r, ok := g.data[b]; ok {
			builder.WriteRune(r)
		} else {
			err = errInvalidByte(nSeptet, septets)
			continue
		}
		nSeptet++
	}
	text := builder.Bytes()
	nDst = len(text)

	if len(dst) < nDst {
		return 0, 0, transform.ErrShortDst
	}

	copy(dst, text)
	return
}

type gsm7Encoder struct {
	packed  bool
	data    map[rune]byte
	dataExt map[rune]byte
}

func (g *gsm7Encoder) Reset() {
	/* no needed */
}

func pack(dst []byte, septets []byte) (nDst int) {
	nSeptet := 0
	for remain := len(septets); remain > 0; {
		// Pack by converting septets into octets.
		switch {
		case remain >= 8:
			dst[nDst+0] = (septets[nSeptet+0] & 0x7F >> 0) | (septets[nSeptet+1] & 0x01 << 7)
			dst[nDst+1] = (septets[nSeptet+1] & 0x7E >> 1) | (septets[nSeptet+2] & 0x03 << 6)
			dst[nDst+2] = (septets[nSeptet+2] & 0x7C >> 2) | (septets[nSeptet+3] & 0x07 << 5)
			dst[nDst+3] = (septets[nSeptet+3] & 0x78 >> 3) | (septets[nSeptet+4] & 0x0F << 4)
			dst[nDst+4] = (septets[nSeptet+4] & 0x70 >> 4) | (septets[nSeptet+5] & 0x1F << 3)
			dst[nDst+5] = (septets[nSeptet+5] & 0x60 >> 5) | (septets[nSeptet+6] & 0x3F << 2)
			dst[nDst+6] = (septets[nSeptet+6] & 0x40 >> 6) | (septets[nSeptet+7] & 0x7F << 1)
			nSeptet += 8
			nDst += 7
		case remain >= 7:
			dst[nDst+0] = (septets[nSeptet+0] & 0x7F >> 0) | (septets[nSeptet+1] & 0x01 << 7)
			dst[nDst+1] = (septets[nSeptet+1] & 0x7E >> 1) | (septets[nSeptet+2] & 0x03 << 6)
			dst[nDst+2] = (septets[nSeptet+2] & 0x7C >> 2) | (septets[nSeptet+3] & 0x07 << 5)
			dst[nDst+3] = (septets[nSeptet+3] & 0x78 >> 3) | (septets[nSeptet+4] & 0x0F << 4)
			dst[nDst+4] = (septets[nSeptet+4] & 0x70 >> 4) | (septets[nSeptet+5] & 0x1F << 3)
			dst[nDst+5] = (septets[nSeptet+5] & 0x60 >> 5) | (septets[nSeptet+6] & 0x3F << 2)
			dst[nDst+6] = septets[nSeptet+6] & 0x40 >> 6
			nSeptet += 7
			nDst += 7
		case remain >= 6:
			dst[nDst+0] = (septets[nSeptet+0] & 0x7F >> 0) | (septets[nSeptet+1] & 0x01 << 7)
			dst[nDst+1] = (septets[nSeptet+1] & 0x7E >> 1) | (septets[nSeptet+2] & 0x03 << 6)
			dst[nDst+2] = (septets[nSeptet+2] & 0x7C >> 2) | (septets[nSeptet+3] & 0x07 << 5)
			dst[nDst+3] = (septets[nSeptet+3] & 0x78 >> 3) | (septets[nSeptet+4] & 0x0F << 4)
			dst[nDst+4] = (septets[nSeptet+4] & 0x70 >> 4) | (septets[nSeptet+5] & 0x1F << 3)
			dst[nDst+5] = septets[nSeptet+5] & 0x60 >> 5
			nSeptet += 6
			nDst += 6
		case remain >= 5:
			dst[nDst+0] = (septets[nSeptet+0] & 0x7F >> 0) | (septets[nSeptet+1] & 0x01 << 7)
			dst[nDst+1] = (septets[nSeptet+1] & 0x7E >> 1) | (septets[nSeptet+2] & 0x03 << 6)
			dst[nDst+2] = (septets[nSeptet+2] & 0x7C >> 2) | (septets[nSeptet+3] & 0x07 << 5)
			dst[nDst+3] = (septets[nSeptet+3] & 0x78 >> 3) | (septets[nSeptet+4] & 0x0F << 4)
			dst[nDst+4] = septets[nSeptet+4] & 0x70 >> 4
			nSeptet += 5
			nDst += 5
		case remain >= 4:
			dst[nDst+0] = (septets[nSeptet+0] & 0x7F >> 0) | (septets[nSeptet+1] & 0x01 << 7)
			dst[nDst+1] = (septets[nSeptet+1] & 0x7E >> 1) | (septets[nSeptet+2] & 0x03 << 6)
			dst[nDst+2] = (septets[nSeptet+2] & 0x7C >> 2) | (septets[nSeptet+3] & 0x07 << 5)
			dst[nDst+3] = septets[nSeptet+3] & 0x78 >> 3
			nSeptet += 4
			nDst += 4
		case remain >= 3:
			dst[nDst+0] = (septets[nSeptet+0] & 0x7F >> 0) | (septets[nSeptet+1] & 0x01 << 7)
			dst[nDst+1] = (septets[nSeptet+1] & 0x7E >> 1) | (septets[nSeptet+2] & 0x03 << 6)
			dst[nDst+2] = septets[nSeptet+2] & 0x7C >> 2
			nSeptet += 3
			nDst += 3
		case remain >= 2:
			dst[nDst+0] = (septets[nSeptet+0] & 0x7F >> 0) | (septets[nSeptet+1] & 0x01 << 7)
			dst[nDst+1] = septets[nSeptet+1] & 0x7E >> 1
			nSeptet += 2
			nDst += 2
		case remain >= 1:
			dst[nDst+0] = septets[nSeptet+0] & 0x7F >> 0
			nSeptet++
			nDst++
		default:
			return
		}
		remain = len(septets) - nSeptet
	}
	return
}

func (g *gsm7Encoder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	if len(src) == 0 {
		return 0, 0, nil
	}

	text := string(src) // work with []rune (a.k.a string) instead of []byte
	septets := make([]byte, 0, len(text))
	for _, r := range text {
		if v, ok := g.data[r]; ok {
			septets = append(septets, v)
		} else if v, ok := g.dataExt[r]; ok {
			septets = append(septets, escapeSequence, v)
		} else {
			return 0, 0, ErrInvalidCharacter
		}
		nSrc++
	}

	nDst = len(septets)
	if g.packed {
		nDst = int(math.Ceil(float64(len(septets)) * 7 / 8))
	}
	if len(dst) < nDst {
		return 0, 0, transform.ErrShortDst
	}

	if !g.packed {
		copy(dst, septets)
		return nDst, nSrc, nil
	}

	nDst = pack(dst, septets)
	return
}
