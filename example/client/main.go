package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/proto"
	"github.com/zhiyin2021/zysms/utils"
)

const (
	user           string        = "010000"      // "900001" //010000,pwd:tDd34J443e3
	password       string        = "tDd34J443e3" //"888888"
	connectTimeout time.Duration = time.Second * 2
)

func startAClient(idx int) {
	// defer wg.Done()
	sms := zysms.New(proto.CMPP2)
	sms.OnConnect = func(c *zysms.Conn) {
		log.Printf("client %d: connect ok", idx)
	}
	sms.OnDisconnect = func(c *zysms.Conn) {
		log.Printf("client %d: disconnect", idx)
	}
	sms.OnError = func(c *zysms.Conn, err error) {
		log.Printf("client %d: err %s", idx, err)
	}
	sms.OnEvent = func(p *zysms.Packet) error {
		switch req := p.Req.(type) {
		case *cmpp.Cmpp2ConnRsp:
			log.Printf("client %d: receive a cmpp connect2 response: %v.", idx, req.Status)
		case *cmpp.Cmpp2SubmitRsp:
			log.Printf("client %d: receive a cmpp submit2 response: %v.", idx, req.MsgId)
		case *cmpp.Cmpp2DeliverReq:
			log.Printf("client %d: receive a cmpp deliver2 request: %v.", idx, req.MsgId)
		case *cmpp.Cmpp3ConnRsp:
			log.Printf("client %d: receive a cmpp connect3 response: %v.", idx, req.Status)
		case *cmpp.Cmpp3SubmitRsp:
			log.Printf("client %d: receive a cmpp submit3 response: %v.", idx, req.MsgId)
		case *cmpp.Cmpp3DeliverReq:
			log.Printf("client %d: receive a cmpp deliver3 request: %v.", idx, req.MsgId)
		default:
			log.Printf("client %d => %d: unknown event: %v", p.Req.Event(), idx, p)
		}
		return nil
	}
	c, err := sms.Dial(":7890", user, password, connectTimeout)
	if err != nil {
		logrus.Printf("client %d: connect error: %s.", idx, err)
		return
	}

	log.Printf("client %d: connect and auth ok", idx)

	t := time.NewTicker(time.Second * 1)
	defer t.Stop()
	for i := 0; i < 1; i++ {
		go func() {
			//submit a message
			cont, err := utils.Utf8ToUcs2("测试 cmpp submit【百度网盘】")
			if err != nil {
				fmt.Printf("client %d: utf8 to ucs2 transform err: %s.", idx, err)
				return
			}
			p := &cmpp.Cmpp2SubmitReq{
				PkTotal:            1,
				PkNumber:           1,
				RegisteredDelivery: 1,
				MsgLevel:           1,
				ServiceId:          "test",
				FeeUserType:        2,
				FeeTerminalId:      "13500002696",
				// FeeTerminalType:    0,
				MsgFmt:         8,
				MsgSrc:         "900001",
				FeeType:        "02",
				FeeCode:        "10",
				ValidTime:      "151105131555101+",
				AtTime:         "",
				SrcId:          "900001",
				DestUsrTl:      1,
				DestTerminalId: []string{"+8613500002696", "8613500002697", "13500002698"},
				// DestTerminalType:   0,
				MsgLength:  uint8(len(cont)),
				MsgContent: string(cont),
			}
			err = c.SendPkt(p, 0)
			if err != nil {
				log.Printf("client %d: send a cmpp submit request error: %s.", idx, err)
				return
			} else {
				log.Printf("client %d: send a cmpp3 submit request ok", idx)
			}
		}()
	}

}

var wg sync.WaitGroup

func main() {
	log.Println("Client example start!")
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go startAClient(i + 1)
	}
	wg.Wait()
	log.Println("Client example ends!")
}
