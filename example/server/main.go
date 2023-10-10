package main

import (
	"bytes"
	"crypto/md5"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/smserror"
	"github.com/zhiyin2021/zysms/utils"
)

const (
	userS     string = "010000"      // "900001" //010000,pwd:tDd34J443e3
	passwordS string = "tDd34J443e3" //"888888"
)

func main() {
	sms := zysms.New(codec.CMPP30, nil)
	sms.OnConnect = func(c *zysms.Conn) {
		c.Logger().Println("server: connect")
	}
	sms.OnDisconnect = func(c *zysms.Conn) {
		c.Logger().Println("server: disconnect")
	}
	sms.OnError = func(c *zysms.Conn, err error) {
		c.Logger().Errorln("server: error: ", err)
	}
	sms.OnRecv = func(p *zysms.Packet) error {
		var err error
		switch req := p.Req.(type) {
		case *cmpp.SubmitReq:
			p.Resp, err = handleSubmit(p, req)
		case *cmpp.ConnReq:
			p.Resp, err = handleLogin(p, req)
		default:
			p.Conn.Logger().Infof("event %T", p)
		}
		return err
	}
	go func() {
		_, err := sms.Listen(":7890")
		if err != nil {
			log.Println("cmpp ListenAndServ error:", err)
			return
		}
	}()
	logrus.Println("server started")
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}
func handleLogin(p *zysms.Packet, req *cmpp.ConnReq) (codec.PDU, error) {
	resp := req.GetResponse().(*cmpp.ConnResp)
	if req.SrcAddr != utils.OctetString(userS, 6) {
		resp.Status = uint32(smserror.ErrnoConnInvalidSrcAddr)
		return resp, smserror.ConnRspStatusErrMap[smserror.ErrnoConnInvalidSrcAddr]
	}

	tm := req.Timestamp
	authSrc := md5.Sum(bytes.Join([][]byte{[]byte(utils.OctetString(userS, 6)),
		make([]byte, 9),
		[]byte(passwordS),
		[]byte(utils.Timestamp2Str(tm))},
		nil))

	if req.AuthSrc != string(authSrc[:]) {
		// conn.Logger().Errorln("handleLogin error: ", cmpp.ConnRspStatusErrMap[smserror.ErrnoConnAuthFailed])
		resp.Status = uint32(smserror.ErrnoConnAuthFailed)
		return resp, smserror.ConnRspStatusErrMap[smserror.ErrnoConnAuthFailed]
	}

	authIsmg := md5.Sum(bytes.Join([][]byte{{byte(resp.Status)},
		authSrc[:],
		[]byte(passwordS)},
		nil))
	resp.AuthIsmg = string(authIsmg[:])
	return resp, nil
}

func handleSubmit(p *zysms.Packet, req *cmpp.SubmitReq) (codec.PDU, error) {
	resp := req.GetResponse().(*cmpp.SubmitResp)
	resp.MsgId = 12878564852733378560
	msg := req.Message.GetMessage()
	for i, d := range req.DestTerminalId {
		p.Conn.Logger().Infof("handleSubmit: handle submit from %s ok!seqId[%d], msgid[%d], srcId[%s], destTerminalId[%s],=>%s\n",
			req.MsgSrc, req.SequenceNumber, resp.MsgId+uint64(i), req.SrcId, d, msg)
	}
	return resp, nil
}
