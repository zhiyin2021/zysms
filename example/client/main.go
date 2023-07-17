package main

import (
	"log"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/codec"
)

// 11101100 01000110 11000100 0000000 00101101 11000000 00000000 00000011
const (
	user           string        = "010000"      // "900001" //010000,pwd:tDd34J443e3
	password       string        = "tDd34J443e3" //"888888"
	connectTimeout time.Duration = time.Second * 2
)

func startAClient(idx int) {
	// defer wg.Done()
	sms := zysms.New(codec.CMPP30)
	sms.OnConnect = func(c *zysms.Conn) {
		log.Printf("client %d: connect ok", idx)
	}
	sms.OnDisconnect = func(c *zysms.Conn) {
		log.Printf("client %d: disconnect", idx)
	}
	sms.OnError = func(c *zysms.Conn, err error) {
		log.Printf("client %d: err %s", idx, err)
	}
	sms.OnRecv = func(p *zysms.Packet) error {
		switch req := p.Req.(type) {
		case *cmpp.ConnResp:
			log.Printf("client %d: receive a cmpp connect response: %v.", idx, req.Status)
		case *cmpp.SubmitResp:
			log.Printf("client %d: receive a cmpp submit response: %v.", idx, req.MsgId)
		case *cmpp.DeliverReq:
			log.Printf("client %d: receive a cmpp deliver request: %v.", idx, req.MsgId)
		default:
			log.Printf("client %d: unknown event: %v", idx, p)
		}
		return nil
	}
	c, err := sms.Dial(":7890", user, password, connectTimeout, false)
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

			msg := "测试 abcdefghijklmnopqrstuvwxyz0123456789##测试 abcdefghijklmnopqrstuvwxyz0123456789##测试 abcdefghijklmnopqrstuvwxyz0123456789##测试 abcdefghijklmnopqrstuvwxyz0123456789##测试 abcdefghijklmnopqrstuvwxyz0123456789【百度网盘】"
			p := cmpp.NewSubmitReq(cmpp.V30).(*cmpp.SubmitReq)
			p.PkTotal = 1
			p.PkNumber = 1
			p.RegisteredDelivery = 1
			p.MsgLevel = 1
			p.ServiceId = "test"
			p.FeeUserType = 2
			p.FeeTerminalId = "13500002696"
			// FeeTerminalType:    0
			p.MsgFmt = 8
			p.MsgSrc = "900001"
			p.FeeType = "02"
			p.FeeCode = "10"
			p.ValidTime = "151105131555101+"
			p.AtTime = ""
			p.SrcId = "900001"
			p.DestUsrTl = 1
			p.DestTerminalId = []string{"+8613500002696", "8613500002697", "13500002698"}
			p.Message.SetMessage(msg, codec.UCS2)
			err = c.SendPDU(p)
			if err != nil {
				log.Printf("client %d: send a cmpp submit request error: %s.", idx, err)
				return
			} else {
				log.Printf("client %d: send a cmpp3 submit request ok", idx)
			}
		}()
		time.Sleep(time.Second * 5)
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
