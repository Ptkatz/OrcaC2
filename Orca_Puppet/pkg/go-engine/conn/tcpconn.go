package conn

import (
	"context"
	"errors"
	"net"
)

type TcpConn struct {
	conn     *net.TCPConn
	listener *net.TCPListener
	cancel   context.CancelFunc
	info     string
}

func (c *TcpConn) Name() string {
	return "tcp"
}

func (c *TcpConn) Read(p []byte) (n int, err error) {
	if c.conn != nil {
		return c.conn.Read(p)
	}
	return 0, errors.New("empty conn")
}

func (c *TcpConn) Write(p []byte) (n int, err error) {
	if c.conn != nil {
		return c.conn.Write(p)
	}
	return 0, errors.New("empty conn")
}

func (c *TcpConn) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	if c.conn != nil {
		return c.conn.Close()
	} else if c.listener != nil {
		return c.listener.Close()
	}
	return nil
}

func (c *TcpConn) Info() string {
	if c.info != "" {
		return c.info
	}
	if c.conn != nil {
		c.info = c.conn.LocalAddr().String() + "<--tcp-->" + c.conn.RemoteAddr().String()
	} else if c.listener != nil {
		c.info = "tcp--" + c.listener.Addr().String()
	} else {
		c.info = "empty tcp conn"
	}
	return c.info
}

func (c *TcpConn) Dial(dst string) (Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp", dst)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	var d net.Dialer
	if gControlOnConnSetup != nil {
		d = net.Dialer{Control: gControlOnConnSetup}
	}
	conn, err := d.DialContext(ctx, "tcp", addr.String())
	if err != nil {
		return nil, err
	}
	c.cancel = nil
	return &TcpConn{conn: conn.(*net.TCPConn)}, nil
}

func (c *TcpConn) Listen(dst string) (Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp", dst)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &TcpConn{listener: listener}, nil
}

func (c *TcpConn) Accept() (Conn, error) {
	conn, err := c.listener.Accept()
	if err != nil {
		return nil, err
	}
	return &TcpConn{conn: conn.(*net.TCPConn)}, nil
}
