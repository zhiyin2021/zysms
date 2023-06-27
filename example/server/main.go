package main

import (
	"bytes"
	"crypto/md5"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/event"
	"github.com/zhiyin2021/zysms/proto"
	"github.com/zhiyin2021/zysms/utils"
)

const (
	userS     string = "010000"      // "900001" //010000,pwd:tDd34J443e3
	passwordS string = "tDd34J443e3" //"888888"
)

func main() {
	sms := zysms.New(proto.CMPP3)
	sms.OnConnect = func(c *zysms.Conn) {
		c.Logger.Println("server: connect")
	}
	sms.OnDisconnect = func(c *zysms.Conn) {
		c.Logger.Println("server: disconnect")
	}
	sms.OnError = func(c *zysms.Conn, err error) {
		c.Logger.Errorln("server: error: ", err)
	}
	sms.Handle(event.SmsEventAuthReq, func(c *zysms.Conn, p proto.Packer) error {
		pkt := p.(*cmpp.CmppConnReq)
		resp, err := handleLogin(pkt)
		c.SendPkt(resp, pkt.SeqId())
		return err
	})
	sms.Handle(event.SmsEventSubmitReq, func(c *zysms.Conn, p proto.Packer) error {
		pkt := p.(*cmpp.Cmpp3SubmitReq)
		resp, err := handleSubmit(pkt)
		c.SendPkt(resp, pkt.SeqId())
		return err
	})
	err := sms.Listen(":7890")
	if err != nil {
		log.Println("cmpp ListenAndServ error:", err)
	}
}
func handleLogin(req *cmpp.CmppConnReq) (proto.Packer, error) {
	resp := &cmpp.Cmpp3ConnRsp{
		Version: cmpp.V30,
	}

	if req.SrcAddr != utils.OctetString(userS, 6) {
		resp.Status = uint32(cmpp.ErrnoConnInvalidSrcAddr)
		return resp, cmpp.ConnRspStatusErrMap[cmpp.ErrnoConnInvalidSrcAddr]
	}

	tm := req.Timestamp
	authSrc := md5.Sum(bytes.Join([][]byte{[]byte(utils.OctetString(userS, 6)),
		make([]byte, 9),
		[]byte(passwordS),
		[]byte(utils.Timestamp2Str(tm))},
		nil))

	if req.AuthSrc != string(authSrc[:]) {
		// conn.Logger().Errorln("handleLogin error: ", cmpp.ConnRspStatusErrMap[cmpp.ErrnoConnAuthFailed])
		resp.Status = uint32(cmpp.ErrnoConnAuthFailed)
		return resp, cmpp.ConnRspStatusErrMap[cmpp.ErrnoConnAuthFailed]
	}

	authIsmg := md5.Sum(bytes.Join([][]byte{{byte(resp.Status)},
		authSrc[:],
		[]byte(passwordS)},
		nil))
	resp.AuthIsmg = string(authIsmg[:])
	return resp, nil
}

func handleSubmit(req *cmpp.Cmpp3SubmitReq) (proto.Packer, error) {
	resp := &cmpp.Cmpp3SubmitRsp{
		MsgId: 12878564852733378560,
	}
	for i, d := range req.DestTerminalId {
		logrus.Printf("handleSubmit: handle submit from %s ok!seqId[%d], msgid[%d], srcId[%s], destTerminalId[%s]\n",
			req.MsgSrc, req.SeqId(), resp.MsgId+uint64(i), req.SrcId, d)
	}
	return resp, nil
}
