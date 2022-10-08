package common

// 客户端连接封装

import (
	"Orca_Puppet/pkg/psexec/encoder"
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"net"
	"runtime/debug"
)

// 会话结构
type Client struct {
	IsSigningRequired bool
	IsAuthenticated   bool
	debug             bool
	securityMode      uint16
	messageId         uint64
	sessionId         uint64
	conn              net.Conn
	dialect           uint16
	options           *ClientOptions
	trees             map[string]uint32
}

// SMB连接参数
type ClientOptions struct {
	Host        string
	Port        int
	Workstation string
	Domain      string
	User        string
	Password    string
	Hash        string
}

func (c *Client) Debug(msg string, err error) {
	if c.debug {
		log.Println("[ DEBUG ] ", msg)
		if err != nil {
			debug.PrintStack()
		}
	}
}

func (c *Client) Send(req interface{}) (res []byte, err error) {
	buf, err := encoder.Marshal(req)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}

	b := new(bytes.Buffer)
	if err = binary.Write(b, binary.BigEndian, uint32(len(buf))); err != nil {
		c.Debug("", err)
		return
	}
	c.Debug("Raw:\n"+hex.Dump(append(b.Bytes(), buf...)), nil)
	rw := bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))
	if _, err = rw.Write(append(b.Bytes(), buf...)); err != nil {
		c.Debug("", err)
		return
	}
	rw.Flush()

	var size uint32
	if err = binary.Read(rw, binary.BigEndian, &size); err != nil {
		c.Debug("", err)
		return
	}
	if size > 0x00FFFFFF {
		return nil, errors.New("Invalid NetBIOS Session message")
	}

	data := make([]byte, size)
	l, err := io.ReadFull(rw, data)
	if err != nil {
		c.Debug("", err)
		return nil, err
	}
	if uint32(l) != size {
		return nil, errors.New("Message size invalid")
	}

	//protID := data[0:4]
	//switch string(protID) {
	//default:
	//	return nil, errors.New("Protocol Not Implemented")
	//case ProtocolSMB:
	//}

	c.messageId++
	return data, nil
}

func (c *Client) WithDebug(debug bool) *Client {
	c.debug = debug
	return c
}

func (c *Client) WithSecurityMode(securityMode uint16) *Client {
	c.securityMode = securityMode
	return c
}

func (c *Client) GetSecurityMode() uint16 {
	return c.securityMode
}

func (c *Client) GetMessageId() uint64 {
	return c.messageId
}

func (c *Client) WithSessionId(sessionId uint64) *Client {
	c.sessionId = sessionId
	return c
}

func (c *Client) GetSessionId() uint64 {
	return c.sessionId
}

func (c *Client) GetConn() net.Conn {
	return c.conn
}

func (c *Client) WithConn(conn net.Conn) *Client {
	c.conn = conn
	return c
}

func (c *Client) WithDialect(dialect uint16) *Client {
	c.dialect = dialect
	return c
}

func (c *Client) WithOptions(clientOptions *ClientOptions) *Client {
	c.options = clientOptions
	return c
}

func (c *Client) GetOptions() *ClientOptions {
	return c.options
}

func (c *Client) WithTrees(trees map[string]uint32) *Client {
	c.trees = trees
	return c
}

func (c *Client) GetTrees() map[string]uint32 {
	return c.trees
}
