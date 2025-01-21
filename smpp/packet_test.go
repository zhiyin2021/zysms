package smpp

import (
	"encoding/hex"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/zhiyin2021/zycli/tools/logger"
	"github.com/zhiyin2021/zysms/codec"
	"golang.org/x/exp/rand"
)

func TestPacket(t *testing.T) {
	data := []string{
		`0000006d000000040000000000000110000201313131313636360001013836313331333131353839363800000000000001000800385c0600a000360039003600380037003100a075284f5c0020004d006900630072006f0073006f0066007400205e1062375b8951684ee37801`,
		`0000001b80000004000000000000cf766e334d71386f6534796500`,
		`000000c100000004000000000000c2ea0005004b50696e7461720001013632383132323038343136393600000000003234313132323033353132393130342b00010000007c284b70696e74617229205974682e4e414e44412e6e6d72207669727475616c206163636f756e742042524920616e6461313034373732313138313632363336352c6261796172207461676968616e2052703736322c323032207365676572612075746b206d656e6768696e646172692062756e67612064656e64612e`,
		`000000730000000400000000000007370002013131313136363600010138363133373939333130303139000000000000010008003e4f7f75289a8c8bc17801002000350038003500350031003400208fdb884c0020004d006900630072006f0073006f0066007400208eab4efd9a8c8bc13002`,
		`000000950000000500000000000bbfb500000036323835373137393632353732000000000400000000010000006769643a536164586d3037345a49207375623a30303120646c7672643a303031207375626d697420646174653a3234313131393131323720646f6e6520646174653a3234313131393131323720737461743a44454c49565244206572723a30303020746578743a20`,
	}
	bufs := [][]byte{}
	for _, d := range data {
		buf, _ := hex.DecodeString(d)
		bufs = append(bufs, buf)

	}
	rand.Seed(uint64(time.Now().UnixNano()))
	total := 0
	go func() {
		for {
			sleep(100) // 1000毫秒以内生成随机数暂停
			num := rand.Intn(100) % 5
			log.Println(num)
			buf := bufs[num]
			total++
			go parseTest(buf)
		}
	}()
	var ss string
	for {
		fmt.Scan(&ss)
		log.Println("total", total)
	}
}

func parseTest(buf []byte) {
	var headerBuf [16]byte
	copy(headerBuf[:], buf[:16])
	header := ParseHeader(headerBuf)
	if pdu, err := CreatePDUFromCmdID(header.CommandID); err == nil {
		// wr := codec.WriterPool.Get(buf)
		// defer codec.WriterPool.Put(wr)
		wr := codec.NewWriter()
		wr.Write(buf)
		// reader := codec.ReaderPool.Get(wr.Bytes())
		// defer codec.ReaderPool.Put(reader)
		reader := codec.NewReader(wr.Bytes())
		pdu.Unmarshal(reader)
		// log.Printf("%#v", pdu)
	} else {
		logger.Errorf("read.CreatePDUFromCmdID %d,%v", header.CommandID, err)
	}
}

// 随机暂停
func sleep(ms int) {
	num := rand.Intn(ms)
	time.Sleep(time.Duration(num) * time.Millisecond)
}
