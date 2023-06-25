package main

import (
	"bytes"
	"crypto/md5"
	"log"
	"net"
	"time"

	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/cmpp/codec"
	"github.com/zhiyin2021/zysms/utils"
)

const (
	userS     string = "900001"
	passwordS string = "888888"
)

func handleLogin(r *cmpp.Response, p *cmpp.Packet, l *log.Logger) (bool, error) {
	req, ok := p.Packer.(*codec.CmppConnReq)
	if !ok {
		// not a connect request, ignore it,
		// go on to next handler
		return true, nil
	}

	l.Println("remote addr:", p.Conn.Conn.RemoteAddr().(*net.TCPAddr).IP.String())
	resp := r.Packer.(*codec.Cmpp3ConnRsp)

	// validate the user and password
	// set the status in the connect response.
	resp.Version = 0x30
	addr := req.SrcAddr
	if addr != utils.OctetString(userS, 6) {
		l.Println("handleLogin error:", codec.ConnRspStatusErrMap[codec.ErrnoConnInvalidSrcAddr])
		resp.Status = uint32(codec.ErrnoConnInvalidSrcAddr)
		return false, codec.ConnRspStatusErrMap[codec.ErrnoConnInvalidSrcAddr]
	}

	tm := req.Timestamp
	authSrc := md5.Sum(bytes.Join([][]byte{[]byte(utils.OctetString(userS, 6)),
		make([]byte, 9),
		[]byte(passwordS),
		[]byte(utils.Timestamp2Str(tm))},
		nil))

	if req.AuthSrc != string(authSrc[:]) {
		l.Println("handleLogin error: ", codec.ConnRspStatusErrMap[codec.ErrnoConnAuthFailed])
		resp.Status = uint32(codec.ErrnoConnAuthFailed)
		return false, codec.ConnRspStatusErrMap[codec.ErrnoConnAuthFailed]
	}

	authIsmg := md5.Sum(bytes.Join([][]byte{[]byte{byte(resp.Status)},
		authSrc[:],
		[]byte(passwordS)},
		nil))
	resp.AuthIsmg = string(authIsmg[:])
	l.Printf("handleLogin: %s login ok\n", addr)

	return false, nil
}

func handleSubmit(r *cmpp.Response, p *cmpp.Packet, l *log.Logger) (bool, error) {
	req, ok := p.Packer.(*codec.Cmpp3SubmitReq)
	if !ok {
		return true, nil // go on to next handler
	}

	resp := r.Packer.(*codec.Cmpp3SubmitRsp)
	resp.MsgId = 12878564852733378560 //0xb2, 0xb9, 0xda, 0x80, 0x00, 0x01, 0x00, 0x00
	for i, d := range req.DestTerminalId {
		l.Printf("handleSubmit: handle submit from %s ok!seqId[%d], msgid[%d], srcId[%s], destTerminalId[%s]\n",
			req.MsgSrc, req.SeqId, resp.MsgId+uint64(i), req.SrcId, d)
	}
	return true, nil
}

func main() {
	var handlers = []cmpp.Handler{
		cmpp.HandlerFunc(handleLogin),
		cmpp.HandlerFunc(handleSubmit),
	}

	err := cmpp.ListenAndServe(":7890",
		codec.V30,
		5*time.Second,
		3,
		nil,
		handlers...,
	)
	if err != nil {
		log.Println("cmpp ListenAndServ error:", err)
	}
}
