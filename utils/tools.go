package utils

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var ErrInvalidUtf8Rune = errors.New("not invalid utf-8 runes")

func Timestamp2() (string, uint32) {
	s := time.Now().Format("0102150405")
	i, _ := strconv.Atoi(s)
	return s, uint32(i)
}

// TimeStamp2Str converts a timestamp(MMDDHHMMSS) int to a string(10 bytes).
func Timestamp2Str(t uint32) string {
	return fmt.Sprintf("%010d", t)
}

func FormatTime(time time.Time) string {
	s := time.Format("060102150405")
	return s + "032+"
}

// MsgFmt 通过消息内容判断，设置编码格式。
// 如果是纯拉丁字符采用0：ASCII串
// 如果含多字节字符，这采用8：UCS-2编码
func MsgFmt(content string) uint8 {
	if len(content) < 2 {
		return 0
	}
	all7bits := len(content) == len([]rune(content))
	if all7bits {
		return 0
	} else {
		return 8
	}
}

func Utf8ToUcs2(in string) ([]byte, error) {
	if !utf8.ValidString(in) {
		return nil, ErrInvalidUtf8Rune
	}

	r := bytes.NewReader([]byte(in))
	t := transform.NewReader(r, unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewEncoder()) // UTF-16 bigendian, no-bom
	out, err := io.ReadAll(t)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func Ucs2ToUtf8(in []byte) ([]byte, error) {
	r := bytes.NewReader(in)
	t := transform.NewReader(r, unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder()) // UTF-16 bigendian, no-bom
	out, err := io.ReadAll(t)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func Utf8ToGB18030(in string) (string, error) {
	if !utf8.ValidString(in) {
		return "", ErrInvalidUtf8Rune
	}

	r := bytes.NewReader([]byte(in))
	t := transform.NewReader(r, simplifiedchinese.GB18030.NewEncoder())
	out, err := io.ReadAll(t)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func GB18030ToUtf8(in string) (string, error) {
	r := bytes.NewReader([]byte(in))
	t := transform.NewReader(r, simplifiedchinese.GB18030.NewDecoder())
	out, err := io.ReadAll(t)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func OctetString(s string, fixedLength int) string {
	length := len(s)
	if length == fixedLength {
		return s
	}

	if length > fixedLength {
		return s[length-fixedLength:]
	}

	return strings.Join([]string{s, string(make([]byte, fixedLength-length))}, "")
}

func TrimOctetString(in []byte) []byte {
	i := bytes.IndexByte(in, 0)
	if i == -1 {
		return in
	} else {
		return in[:i]
	}
}

func TrimStr(in []byte) string {
	return string(TrimOctetString(in))
}

var NumTable = [16]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func Uint64HexString(i uint64) string {
	var sb strings.Builder
	sb.Grow(16)

	for shift := 60; shift >= 0; shift -= 4 {
		sb.WriteByte(NumTable[(i>>shift)&0x0f])
	}
	return sb.String()
}

func Uint32HexString(i uint32) string {
	var sb strings.Builder
	sb.Grow(8)

	for shift := 28; shift >= 0; shift -= 4 {
		sb.WriteByte(NumTable[(i>>shift)&0x0f])
	}
	return sb.String()
}

func Bytes2StringSlice(in []byte, pl int) (ret []string) {
	if len(in) <= pl {
		return []string{TrimStr(in)}
	} else {
		part := len(in) / pl
		ret = make([]string, part)
		for i := 0; i < part; i++ {
			ret[i] = TrimStr(in[i*pl : (i+1)*pl])
		}
	}
	return
}

func MsgSlices(fmt uint8, content string) (slices [][]byte) {
	var msgBytes []byte
	// 含中文
	if fmt == 8 {
		msgBytes, _ = Utf8ToUcs2(content)
		slices = ToTPUDHISlices(msgBytes, 140)
	} else {
		// 纯英文
		msgBytes = []byte(content)
		slices = ToTPUDHISlices(msgBytes, 160)
	}
	return
}

// ToTPUDHISlices 拆分为长短信切片
// 纯ASCII内容的拆分 pkgLen = 160
// 含中文内容的拆分   pkgLen = 140
func ToTPUDHISlices(content []byte, pkgLen int) (rt [][]byte) {
	if len(content) < pkgLen {
		return [][]byte{content}
	}

	headLen := 6
	bodyLen := pkgLen - headLen
	parts := len(content) / bodyLen
	tailLen := len(content) % bodyLen
	if tailLen != 0 {
		parts++
	}
	// 分片消息组的标识，用于收集组装消息
	groupId := byte(time.Now().UnixNano() & 0xff)
	var part []byte
	for i := 0; i < parts; i++ {
		if i != parts-1 {
			part = make([]byte, pkgLen)
		} else {
			// 最后一片
			part = make([]byte, 6+tailLen)
		}
		part[0], part[1], part[2] = 0x05, 0x00, 0x03
		part[3] = groupId
		part[4], part[5] = byte(parts), byte(i+1)
		if i != parts-1 {
			copy(part[6:pkgLen], content[bodyLen*i:bodyLen*(i+1)])
		} else {
			copy(part[6:], content[0:tailLen])
		}
		rt = append(rt, part)
	}
	return rt
}

func RandomNum(min, max int) int {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rnd.Intn(max-min) + min
}

// Charsets
const (
	Uppercase  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Lowercase  = "abcdefghijklmnopqrstuvwxyz"
	Alphabetic = Uppercase + Lowercase
)

func Md5(data string) string {
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func RandomStr(length uint8, charsets ...string) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	charset := strings.Join(charsets, "")
	if charset == "" {
		charset = Alphabetic
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rnd.Int63()%int64(len(charset))]
	}
	return string(b)
}

// DiceCheck 投概率骰子，得到结果比给定数字大则返回true，否则返回false
func DiceCheck(prob float64) bool {
	return float64(rand.Intn(10000))/10000.0 > prob
}

// StructCopy 采用反射拷贝结构体属性
func StructCopy(from, to any) {
	fromValue := reflect.ValueOf(from)
	toValue := reflect.ValueOf(to)

	// 必须是指针类型
	if fromValue.Kind() != reflect.Ptr || toValue.Kind() != reflect.Ptr {
		return
	}
	// 均不可为空
	if fromValue.IsNil() || toValue.IsNil() {
		return
	}

	// 获取到来源数据
	fromElem := fromValue.Elem()
	// 需要的数据
	toElem := toValue.Elem()

	for i := 0; i < toElem.NumField(); i++ {
		toField := toElem.Type().Field(i)

		// 看看来源的结构体中是否有这个属性
		fromFieldName, ok := fromElem.Type().FieldByName(toField.Name)
		// 存在相同的属性名称并且类型一致
		if ok && fromFieldName.Type == toField.Type {
			toElem.Field(i).Set(fromElem.FieldByName(toField.Name))
		}
	}
}

func MapItem[T any](m map[string]string, key string, def T) (ret T) {
	ret = def
	if val, ok := m[key]; ok {
		v := reflect.ValueOf(&ret).Elem()
		switch v.Kind() {
		case reflect.String:
			v.SetString(val)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if n, err := strconv.ParseInt(val, 0, 0); err == nil {

				v.SetInt(n)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if n, err := strconv.ParseUint(val, 0, 0); err == nil {
				v.SetUint(n)
			}
		}
	}
	return
}
