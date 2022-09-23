package conn

import (
	"context"
	"crypto/tls"
	"errors"
	"Orca_Server/pkg/go-engine/common"
	"github.com/lucas-clemente/quic-go"
	"github.com/xtaci/smux"
	"net"
)

type QuicConn struct {
	qsession quic.Connection
	session  *smux.Session
	qsteam   quic.Stream
	stream   *smux.Stream
	listener quic.Listener
	info     string
}

func (c *QuicConn) Name() string {
	return "quic"
}

func (c *QuicConn) Read(p []byte) (n int, err error) {
	if c.stream != nil {
		return c.stream.Read(p)
	}
	return 0, errors.New("empty conn")
}

func (c *QuicConn) Write(p []byte) (n int, err error) {
	if c.stream != nil {
		return c.stream.Write(p)
	}
	return 0, errors.New("empty conn")
}

func (c *QuicConn) Close() error {
	if c.stream != nil {
		return c.stream.Close()
	} else if c.listener != nil {
		return c.listener.Close()
	}
	return nil
}

func (c *QuicConn) Info() string {
	if c.info != "" {
		return c.info
	}
	if c.session != nil {
		c.info = c.qsession.LocalAddr().String() + "<--quic-->" + c.qsession.RemoteAddr().String()
	} else if c.listener != nil {
		c.info = "kcp--" + c.listener.Addr().String()
	} else {
		c.info = "empty kcp conn"
	}
	return c.info
}

func (c *QuicConn) Dial(dst string) (Conn, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"QuicConn"},
	}

	var lc net.ListenConfig
	if gControlOnConnSetup != nil {
		lc.Control = gControlOnConnSetup
	}

	laddr := &net.UDPAddr{}
	pconn, err := lc.ListenPacket(context.Background(), "udp", laddr.String())
	if err != nil {
		return nil, err
	}

	udpAddr, err := net.ResolveUDPAddr("udp", dst)
	if err != nil {
		return nil, err
	}

	session, err := quic.Dial(pconn, udpAddr, dst, tlsConf, nil)
	if err != nil {
		return nil, err
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		return nil, err
	}

	ss, err := smux.Client(stream, nil)
	if err != nil {
		return nil, err
	}

	st, err := ss.OpenStream()
	if err != nil {
		return nil, err
	}

	return &QuicConn{qsession: session, session: ss, qsteam: stream, stream: st}, nil
}

func (c *QuicConn) Listen(dst string) (Conn, error) {
	config, err := common.GenerateTLSConfig("QuicConn")
	if err != nil {
		return nil, err
	}

	listener, err := quic.ListenAddr(dst, config, nil)
	if err != nil {
		return nil, err
	}

	return &QuicConn{listener: listener}, nil
}

func (c *QuicConn) Accept() (Conn, error) {
	session, err := c.listener.Accept(context.Background())
	if err != nil {
		return nil, err
	}

	stream, err := session.AcceptStream(context.Background())
	if err != nil {
		return nil, err
	}

	ss, err := smux.Server(stream, nil)
	if err != nil {
		return nil, err
	}

	st, err := ss.AcceptStream()
	if err != nil {
		return nil, err
	}

	return &QuicConn{qsession: session, session: ss, qsteam: stream, stream: st}, nil
}
