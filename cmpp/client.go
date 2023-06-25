package cmpp

import (
	"errors"
	"net"
	"time"

	"github.com/zhiyin2021/zysms/cmpp/codec"
	"github.com/zhiyin2021/zysms/proto"
)

var ErrNotCompleted = errors.New("data not being handled completed")
var ErrRespNotMatch = errors.New("the response is not matched with the request")

// Client stands for one client-side instance, just like a session.
// It may connect to the server, send & recv cmpp packets and terminate the connection.
type Client struct {
	conn *Conn
	typ  codec.Version
}

// New establishes a new cmpp client.
func NewClient(typ codec.Version) *Client {
	return &Client{
		typ: typ,
	}
}

// Connect connect to the cmpp server in block mode.
// It sends login packet, receive and parse connect response packet.
func (cli *Client) Connect(servAddr, user, password string, timeout time.Duration) error {
	var err error
	conn, err := net.DialTimeout("tcp", servAddr, timeout)
	if err != nil {
		return err
	}
	cli.conn = newConn(conn, cli.typ)
	defer func() {
		if err != nil {
			if cli.conn != nil {
				cli.conn.Close()
			}
		}
	}()
	cli.conn.SetState(CONN_CONNECTED)

	// Login to the server.
	req := &codec.CmppConnReq{
		SrcAddr: user,
		Secret:  password,
		Version: cli.typ,
	}

	_, err = cli.SendReq(req)
	if err != nil {
		return err
	}

	p, err := cli.conn.RecvAndUnpackPkt(timeout)
	if err != nil {
		return err
	}

	var ok bool
	var status uint8
	if cli.typ == codec.V20 || cli.typ == codec.V21 {
		var rsp *codec.Cmpp2ConnRsp
		rsp, ok = p.(*codec.Cmpp2ConnRsp)
		if !ok {
			err = ErrRespNotMatch
			return err
		}
		status = rsp.Status
	} else {
		var rsp *codec.Cmpp3ConnRsp
		rsp, ok = p.(*codec.Cmpp3ConnRsp)
		if !ok {
			err = ErrRespNotMatch
			return err
		}
		status = uint8(rsp.Status)
	}

	if status != 0 {
		if status <= codec.ErrnoConnOthers { //ErrnoConnOthers = 5
			err = codec.ConnRspStatusErrMap[status]
		} else {
			err = codec.ConnRspStatusErrMap[codec.ErrnoConnOthers]
		}
		return err
	}

	cli.conn.SetState(CONN_AUTHOK)
	return nil
}

func (cli *Client) Disconnect() {
	if cli.conn != nil {
		cli.conn.Close()
	}
}

// SendReq pack the cmpp request packet structure and send it to the other peer.
func (cli *Client) SendReq(packet proto.Packer) (uint32, error) {
	seq := cli.conn.SeqId()
	return seq, cli.conn.SendPkt(packet, seq)
}

// SendRsp pack the cmpp response packet structure and send it to the other peer.
func (cli *Client) SendRsp(packet proto.Packer, seqId uint32) error {
	return cli.conn.SendPkt(packet, seqId)
}

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (cli *Client) RecvAndUnpackPkt(timeout time.Duration) (interface{}, error) {
	return cli.conn.RecvAndUnpackPkt(timeout)
}
