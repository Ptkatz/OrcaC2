package conn

import (
	"context"
	"errors"
	"github.com/xtaci/kcp-go"
	"github.com/xtaci/smux"
	"net"
)

type KcpConn struct {
	session  *smux.Session
	stream   *smux.Stream
	listener *kcp.Listener
	info     string
}

func (c *KcpConn) Name() string {
	return "kcp"
}

func (c *KcpConn) Read(p []byte) (n int, err error) {
	if c.stream != nil {
		return c.stream.Read(p)
	}
	return 0, errors.New("empty conn")
}

func (c *KcpConn) Write(p []byte) (n int, err error) {
	if c.stream != nil {
		return c.stream.Write(p)
	}
	return 0, errors.New("empty conn")
}

func (c *KcpConn) Close() error {
	if c.session != nil {
		return c.session.Close()
	} else if c.listener != nil {
		return c.listener.Close()
	}
	return nil
}

func (c *KcpConn) Info() string {
	if c.info != "" {
		return c.info
	}
	if c.session != nil {
		c.info = c.session.LocalAddr().String() + "<--kcp-->" + c.session.RemoteAddr().String()
	} else if c.listener != nil {
		c.info = "kcp--" + c.listener.Addr().String()
	} else {
		c.info = "empty kcp conn"
	}
	return c.info
}

func (c *KcpConn) Dial(dst string) (Conn, error) {
	var lc net.ListenConfig
	if gControlOnConnSetup != nil {
		lc.Control = gControlOnConnSetup
	}

	laddr := &net.UDPAddr{}
	pconn, err := lc.ListenPacket(context.Background(), "udp", laddr.String())
	if err != nil {
		return nil, err
	}

	conn, err := kcp.NewConn(dst, nil, 0, 0, pconn.(*net.UDPConn))
	if err != nil {
		return nil, err
	}

	c.setParam(conn)

	session, err := smux.Client(conn, nil)
	if err != nil {
		return nil, err
	}

	stream, err := session.OpenStream()
	if err != nil {
		return nil, err
	}

	return &KcpConn{session: session, stream: stream}, nil
}

func (c *KcpConn) Listen(dst string) (Conn, error) {
	listener, err := kcp.Listen(dst)
	if err != nil {
		return nil, err
	}

	listener.(*kcp.Listener).SetReadBuffer(4 * 1024 * 1024)
	listener.(*kcp.Listener).SetWriteBuffer(4 * 1024 * 1024)
	listener.(*kcp.Listener).SetDSCP(46)

	return &KcpConn{listener: listener.(*kcp.Listener)}, nil
}

func (c *KcpConn) Accept() (Conn, error) {
	conn, err := c.listener.Accept()
	if err != nil {
		return nil, err
	}

	c.setParam(conn.(*kcp.UDPSession))

	session, err := smux.Server(conn, nil)
	if err != nil {
		return nil, err
	}

	stream, err := session.AcceptStream()
	if err != nil {
		return nil, err
	}

	return &KcpConn{session: session, stream: stream}, nil
}

func (c *KcpConn) setParam(conn *kcp.UDPSession) {
	conn.SetStreamMode(true)
	conn.SetWindowSize(10000, 10000)
	conn.SetReadBuffer(16 * 1024 * 1024)
	conn.SetWriteBuffer(16 * 1024 * 1024)
	conn.SetNoDelay(0, 100, 1, 1)
	conn.SetMtu(500)
	conn.SetACKNoDelay(false)
}
