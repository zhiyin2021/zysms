package zysms

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/smgp"
	"github.com/zhiyin2021/zysms/smserror"
)

type smgp_action struct {
	*sms_conn
}

func newSmgp(conn *sms_conn) *smgp_action {
	return &smgp_action{conn}
}

func (c *smgp_action) login(uid string, pwd string) error {
	// Login to the server.
	req := smgp.NewLoginReq(c.Typ).(*smgp.LoginReq)
	req.ClientID = uid
	req.Secret = pwd
	req.Version = c.Typ

	err := c.SendPDU(req)
	if err != nil {
		return err
	}
	p, err := c.recv()
	if err != nil {
		return err
	}
	var status uint8

	if rsp, ok := p.(*smgp.LoginResp); ok {
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
func (c *smgp_action) logout() {
	c.SendPDU(smgp.NewExitReq(c.Typ))
}

// RecvAndUnpackPkt receives smgp byte stream, and unpack it to some smgp packet structure.
func (c *smgp_action) recv() (codec.PDU, error) {
	if c.Connected == enum.CONN_DISCONNECTED {
		return nil, smserror.ErrConnIsClosed
	}

	pdu, err := smgp.Parse(c.Conn, c.Typ)
	if err != nil {
		return nil, err
	}

	switch p := pdu.(type) {
	case *smgp.ActiveTestReq: // 当收到心跳请求,内部直接回复心跳,并递归继续获取数据
		if c.autoActiveResp {
			resp := p.GetResponse()
			c.SendPDU(resp)
		}
	case *smgp.ActiveTestResp: // 当收到心跳回复,内部直接处理,并递归继续获取数据
		atomic.StoreInt32(&c.counter, 0)
		c.activeTestResp(p.GetSequenceNumber())
	case *smgp.LoginResp: // 当收到登录回复,内部先校验版本
		if c.checkVer && p.Version != c.Typ {
			return nil, fmt.Errorf("smgp version not match [ local: %d != remote: %d ]", c.Typ, p.Version)
		}
	case *smgp.ExitReq:
		resp := p.GetResponse()
		c.SendPDU(resp)
		time.Sleep(100 * time.Millisecond)
		return nil, smserror.ErrConnIsClosed
	case *smgp.LoginReq:
		switch p.Version {
		case smgp.V20, smgp.V30:
			// 服务端自适应版本
			c.Typ = p.Version
			fallthrough
		case 0:
			c.logger = c.logger.WithFields(logrus.Fields{"v": c.Protocol.String(), "v1": c.Typ})
		default:
			return nil, fmt.Errorf("smgp version not support [ %d ]", p.Version)
		}
	}
	return pdu, nil

}

func (c *smgp_action) active_test() error {
	p := smgp.NewActiveTestReq(c.Typ)
	c.activeTestReq(p.GetSequenceNumber())
	return c.SendPDU(p)
}
