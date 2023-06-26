package main

import (
	"bytes"
	"crypto/md5"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms"
	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/proto"
	"github.com/zhiyin2021/zysms/utils"
)

const (
	userS     string = "010000"      // "900001" //010000,pwd:tDd34J443e3
	passwordS string = "tDd34J443e3" //"888888"
)

func main() {
	ln, err := zysms.Listen(":7890", zysms.CMPP3)
	if err != nil {
		log.Println("cmpp ListenAndServ error:", err)
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("cmpp Accept error:", err)
			continue
		}
		go handleRun(conn)
	}
}
func handleRun(conn zysms.SmsConn) {
	defer conn.Close()
	log.Println("cmpp Accept a new connection:", conn.RemoteAddr())

	for {
		pkt, err := conn.RecvPkt(0)
		if err != nil {
			logrus.Printf("handle.recv err: %v", err)
			if err.Error() == "EOF" {
				return
			}
			if _, ok := err.(enum.SmsError); ok {
				continue
			}
			return
		}
		var resp proto.Packer

		switch p := pkt.(type) {
		case *cmpp.CmppConnReq:
			{
				resp, err = handleLogin(p)
				if err != nil {
					conn.Logger().Errorf("handleLogin error: %v", err)
				}
			}
		case *cmpp.Cmpp3SubmitReq:
			{
				resp, err = handleSubmit(p)
				if err != nil {
					conn.Logger().Errorf("handleSubmit error: %v", err)
				}
			}
		default:
			continue
		}
		err1 := conn.SendPkt(resp, pkt.SeqId())
		if err1 != nil {
			conn.Logger().Errorf("sendPkt error: %v", err)
			return
		}
		if err != nil {
			return
		}
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
