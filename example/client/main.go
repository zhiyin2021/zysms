package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/zhiyin2021/zysms"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/utils"
)

const (
	user           string        = "010000"      // "900001" //010000,pwd:tDd34J443e3
	password       string        = "tDd34J443e3" //"888888"
	connectTimeout time.Duration = time.Second * 2
)

func startAClient(idx int) {
	c, err := zysms.Dial(":7890", zysms.CMPP3, user, password, connectTimeout)
	if err != nil {
		log.Printf("client %d: connect error: %s.", idx, err)
		return
	}
	defer wg.Done()
	defer c.Close()

	log.Printf("client %d: connect and auth ok", idx)

	t := time.NewTicker(time.Second * 5)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			//submit a message
			cont, err := utils.Utf8ToUcs2("测试gocmpp submit")
			if err != nil {
				fmt.Printf("client %d: utf8 to ucs2 transform err: %s.", idx, err)
				return
			}
			p := &cmpp.Cmpp3SubmitReq{
				PkTotal:            1,
				PkNumber:           1,
				RegisteredDelivery: 1,
				MsgLevel:           1,
				ServiceId:          "test",
				FeeUserType:        2,
				FeeTerminalId:      "13500002696",
				FeeTerminalType:    0,
				MsgFmt:             8,
				MsgSrc:             "900001",
				FeeType:            "02",
				FeeCode:            "10",
				ValidTime:          "151105131555101+",
				AtTime:             "",
				SrcId:              "900001",
				DestUsrTl:          1,
				DestTerminalId:     []string{"13500002696"},
				DestTerminalType:   0,
				MsgLength:          uint8(len(cont)),
				MsgContent:         string(cont),
			}

			err = c.SendPkt(p, 0)
			if err != nil {
				log.Printf("client %d: send a cmpp3 submit request error: %s.", idx, err)
			} else {
				log.Printf("client %d: send a cmpp3 submit request ok", idx)
			}
		default:
		}

		// recv packets
		i, err := c.RecvPkt(0)
		if err != nil {
			log.Printf("client %d: client read and unpack pkt error: %s.", idx, err)
			break
		}

		switch p := i.(type) {
		case *cmpp.Cmpp3SubmitRsp:
			log.Printf("client %d: receive a cmpp3 submit response: %d => %v.", idx, p.SeqId(), p)
		case *cmpp.CmppActiveTestReq:
			log.Printf("client %d: receive a cmpp active request: %v.", idx, p)
			// rsp := &cmpp.CmppActiveTestRsp{}
			// err := c.SendPkt(rsp, p.SeqId())
			// if err != nil {
			// 	log.Printf("client %d: send cmpp active response error: %s.", idx, err)
			// 	break
			// }
		case *cmpp.CmppActiveTestRsp:
			log.Printf("client %d: receive a cmpp activetest response: %v.", idx, p)

		case *cmpp.CmppTerminateReq:
			log.Printf("client %d: receive a cmpp terminate request: %v.", idx, p)
			rsp := &cmpp.CmppTerminateRsp{}
			err := c.SendPkt(rsp, p.SeqId())
			if err != nil {
				log.Printf("client %d: send cmpp terminate response error: %s.", idx, err)
				break
			}
		case *cmpp.CmppTerminateRsp:
			log.Printf("client %d: receive a cmpp terminate response: %v.", idx, p)
		}
	}
}

var wg sync.WaitGroup

func main() {
	log.Println("Client example start!")
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go startAClient(i + 1)
	}
	wg.Wait()
	log.Println("Client example ends!")
}
