package zysms

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/sgip"
	"github.com/zhiyin2021/zysms/smserror"
)

type sgip_action struct {
	*sms_conn
}

func newSgip(conn *sms_conn) *sgip_action {
	return &sgip_action{conn}
}

func (c *sgip_action) login(uid string, pwd string) error {
	// Login to the server.
	req := sgip.NewBindReq(c.Typ, c.nodeId).(*sgip.BindReq)
	req.LoginName = uid
	req.LoginPassword = pwd
	req.LoginType = 1
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

	if rsp, ok := p.(*sgip.BindResp); ok {
		if c.checkVer && rsp.Version != c.Typ {
			return smserror.ErrVersionNotMatch
		}
		status = uint8(rsp.Status)
	} else {
		return smserror.ErrRespNotMatch
	}

	if status != 0 {
		return smserror.NewSmsErr(int(status), "sgip.login.error")
	}
	c.IsAuth = true
	return nil
}
func (c *sgip_action) logout() {
	c.SendPDU(sgip.NewUnbindReq(c.Typ, c.nodeId))
}

// RecvAndUnpackPkt receives sgip byte stream, and unpack it to some sgip packet structure.
func (c *sgip_action) recv() (codec.PDU, error) {
	if c.Connected == enum.CONN_DISCONNECTED {
		return nil, smserror.ErrConnIsClosed
	}

	pdu, err := sgip.Parse(c.Conn, c.Typ, c.nodeId)
	if err != nil {
		return nil, err
	}

	switch p := pdu.(type) {

	case *sgip.ReportReq: // 当收到心跳请求,内部直接回复心跳,并递归继续获取数据
		resp := p.GetResponse().(*sgip.ReportResp)
		resp.Status = 0
		c.SendPDU(resp)
	case *sgip.ReportResp: // 当收到心跳回复,内部直接处理,并递归继续获取数据
		atomic.StoreInt32(&c.counter, 0)
	case *sgip.BindResp: // 当收到登录回复,内部先校验版本
		if c.checkVer && p.Version != c.Typ {
			return nil, fmt.Errorf("sgip version not match [ local: %d != remote: %d ]", c.Typ, p.Version)
		}
	case *sgip.UnbindReq: // 当收到退出请求,内部直接回复退出
		resp := p.GetResponse()
		c.SendPDU(resp)
		time.Sleep(100 * time.Millisecond)
		return nil, smserror.ErrConnIsClosed
	case *sgip.BindReq:
		switch p.Version {
		case sgip.V12:
			// 服务端自适应版本
			c.Typ = p.Version
			c.logger = c.logger.With("v", c.Protocol.String(), "v1", c.Typ)
		default:
			return nil, fmt.Errorf("cmpp version not support [ %d ]", p.Version)
		}
	}
	return pdu, nil

}

func (c *sgip_action) active_test() error {
	p := sgip.NewReportReq(c.Typ, c.nodeId).(*sgip.ReportReq)
	return c.SendPDU(p)
}
