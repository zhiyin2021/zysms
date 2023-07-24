package main

import (
	"log"
	"time"

	"github.com/zhiyin2021/zysms/cmpp"
)

func main() {
	tm := time.Now()
	for i := 0; i < 1000000; i++ {
		// cmpp.CreatePDUFromCmdID(cmpp.CMPP_ACTIVE_TEST, cmpp.V30)
	}
	s1 := time.Since(tm)
	header := cmpp.Header{CommandID: cmpp.CMPP_ACTIVE_TEST}
	tm = time.Now()
	for i := 0; i < 1000000; i++ {
		cmpp.CreatePDUHeader(header, cmpp.V30)
	}
	s2 := time.Since(tm)
	log.Println("time1:", s1.Microseconds(), s1.Milliseconds(), "time2:", s2.Microseconds(), s2.Milliseconds())
}
