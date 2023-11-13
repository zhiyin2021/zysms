package zysms

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zhiyin2021/zysms/codec"
	"github.com/zhiyin2021/zysms/enum"
	"github.com/zhiyin2021/zysms/smpp"
	"github.com/zhiyin2021/zysms/smserror"
)

type smpp_action struct {
	*sms_conn
}

func newSmpp(conn *sms_conn) *smpp_action {
	return &smpp_action{conn}
}

func (c *smpp_action) login(uid string, pwd string) error {
	// Login to the server.
	req := smpp.NewBindRequest(smpp.Transceiver)
	req.SystemID = uid
	req.Password = pwd
	req.InterfaceVersion = c.Typ
	if c.extParam != nil && c.extParam["system_type"] != "" {
		req.SystemType = c.extParam["system_type"]
	}
	err := c.SendPDU(req)
	if err != nil {
		return err
	}
	pdu, err := c.recv()
	if err != nil {
		return err
	}
	if header, ok := pdu.GetHeader().(*smpp.Header); ok {
		status := header.CommandStatus
		if status != smpp.ESME_ROK {
			return smserror.NewSmsErr(int(status), "smpp.login.error") //fmt.Errorf("login error: %v", status)
		}
		c.IsAuth = true
	}
	return nil
}

func (c *smpp_action) logout() {
	c.SendPDU(smpp.NewUnbind())
}

// RecvAndUnpackPkt receives smpp byte stream, and unpack it to some smpp packet structure.
func (c *smpp_action) recv() (codec.PDU, error) {
	if c.Connected == enum.CONN_DISCONNECTED {
		return nil, smserror.ErrConnIsClosed
	}
	pdu, err := smpp.Parse(c.Conn)
	if err != nil {
		return nil, err
	}
	switch p := pdu.(type) {
	case *smpp.EnquireLink: // 当收到心跳请求,内部直接回复心跳,并递归继续获取数据
		if c.autoActiveResp {
			resp := p.GetResponse()
			c.SendPDU(resp)
		}
	case *smpp.EnquireLinkResp: // 当收到心跳回复,内部直接处理,并递归继续获取数据
		atomic.AddInt32(&c.counter, -1)
		// c.activeLast = time.Now()
	case *smpp.BindResp: // 当收到登录回复,内部先校验版本
		if p.CommandStatus != smpp.ESME_ROK {
			return nil, smserror.NewSmsErr(int(p.CommandStatus), "smpp.login.error")
		}
	case *smpp.Unbind:
		resp := p.GetResponse()
		c.SendPDU(resp)
		time.Sleep(100 * time.Millisecond)
		return nil, smserror.ErrConnIsClosed
	case *smpp.BindRequest:
		switch p.InterfaceVersion {
		case smpp.V33, smpp.V34:
			// 服务端自适应版本
			c.Typ = p.InterfaceVersion
			fallthrough
		case 0:
			c.logger = c.logger.WithFields(logrus.Fields{"v": c.Protocol.String(), "v1": c.Typ})
		default:
			return nil, fmt.Errorf("smpp version not support [ %d ]", p.InterfaceVersion)
		}
	}
	return pdu, nil
}

func (c *smpp_action) active_test() error {
	p := smpp.NewEnquireLink()
	return c.SendPDU(p)
}
