package zysms

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/zhiyin2021/zysms/cmpp"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/smserror"
)

type cmpp_action struct {
	*sms_conn
}

func newCmpp(conn *sms_conn) sms_action {
	return &cmpp_action{conn}
}
func (c *cmpp_action) login(uid string, pwd string) error {
	// Login to the server.
	req := cmpp.NewConnReq(c.Typ).(*cmpp.ConnReq)
	req.SrcAddr = uid
	req.Secret = pwd
	req.Version = c.Typ

	if err := c.SendPDU(req); err != nil {
		return err
	}
	p, err := c.recv()
	if err != nil {
		return err
	}
	var status uint8

	if rsp, ok := p.(*cmpp.ConnResp); ok {
		if c.checkVer && rsp.Version != c.Typ {
			return smserror.ErrVersionNotMatch
		}
		status = uint8(rsp.Status)
	} else {
		return smserror.ErrRespNotMatch
	}

	if status != 0 {
		return smserror.NewSmsErr(int(status), "cmpp.login.error")
	}
	c.IsAuth = true
	return nil
}
func (c *cmpp_action) logout() {
	c.SendPDU(cmpp.NewCancelReq(c.Typ))
}

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (c *cmpp_action) recv() (codec.PDU, error) {
	if c.Connected == enum.CONN_DISCONNECTED {
		return nil, smserror.ErrConnIsClosed
	}
	pdu, err := cmpp.Parse(c.Conn, c.Typ, c.logger)
	if err != nil {
		return nil, err
	}

	switch p := pdu.(type) {
	case *cmpp.ActiveTestReq: // 当收到心跳请求,内部直接回复心跳,并递归继续获取数据
		resp := p.GetResponse()
		c.SendPDU(resp)
		// if !c.activePeer {
		// 	c.activePeer = true
		// }
	case *cmpp.ActiveTestResp: // 当收到心跳回复,内部直接处理,并递归继续获取数据
		atomic.StoreInt32(&c.counter, 0)
		c.activeTestResp(p.GetSequenceNumber())
	case *cmpp.ConnResp: // 当收到登录回复,内部先校验版本
		if c.checkVer && p.Version != c.Typ {
			return nil, fmt.Errorf("cmpp version not match [ local: %d != remote: %d ]", c.Typ, p.Version)
		}
	case *cmpp.CancelReq:
		resp := p.GetResponse()
		c.SendPDU(resp)
		time.Sleep(100 * time.Millisecond)
		return nil, smserror.ErrConnIsClosed
	case *cmpp.ConnReq:
		switch p.Version {
		case cmpp.V20, cmpp.V30, cmpp.V21:
			// 服务端自适应版本
			c.Typ = p.Version
			fallthrough
		case 0:
			c.logger = c.logger.With("v", c.Protocol.String(), "v1", c.Typ)
		default:
			return nil, fmt.Errorf("cmpp version not support [ %d ]", p.Version)
		}
	}
	return pdu, nil

}

func (c *cmpp_action) active_test() error {
	p := cmpp.NewActiveTestReq(c.Typ)
	c.activeTestReq(p.GetSequenceNumber())
	return c.SendPDU(p)
}
