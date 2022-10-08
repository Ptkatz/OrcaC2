package socketgo

import (
	"errors"
	"fmt"
	"Orca_Puppet/pkg/go-engine/rbuffergo"
	"net"
	"time"
)

const (
	LU_CLIENT int = 1
	LU_SERVER int = 2
)

type LuConfig struct {
	Type                 int
	Ip                   string
	Port                 int
	RecvBufferSize       int
	SendBufferSize       int
	SocketRecvBufferSize int
	SocketSendBufferSize int
	IsKeepAlive          bool
	KeepInterval         int
	Linger               int
	IsNonBlocking        bool
	IsNoDelay            bool
	IsEncrypt            bool
	UserData             interface{}
}

type LuSocket struct {
	config   *LuConfig
	listener *net.TCPListener
	son      map[string]luConn
	self     *luConn
}

type luConn struct {
	recvBuffer *rbuffergo.RBuffergo
	sendBuffer *rbuffergo.RBuffergo
	conn       *net.TCPConn
}

func New(config *LuConfig) (*LuSocket, error) {
	s := &LuSocket{}

	if config.RecvBufferSize == 0 || config.SendBufferSize == 0 || config.SocketRecvBufferSize == 0 || config.SocketSendBufferSize == 0 {
		return nil, errors.New("need set BufferSize")
	}
	if config.Type == LU_SERVER && config.Ip == "" {
		return nil, errors.New("need set Ip")
	}
	if config.Port == 0 {
		return nil, errors.New("need set Port")
	}
	s.config = config

	if config.Type == LU_SERVER {
		err := server(s)
		if err != nil {
			return nil, err
		}
	} else if config.Type == LU_CLIENT {
		err := client(s)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("need set Type")
	}

	return s, nil
}

func server(s *LuSocket) error {

	s.son = make(map[string]luConn)

	addr := ""
	if s.config.Ip == "" {
		addr = fmt.Sprintf(":%d", s.config.Port)
	} else {
		addr = fmt.Sprintf("%s:%d", s.config.Ip, s.config.Port)
	}

	listenAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		return err
	}

	s.listener = listener

	go accept(s)

	return nil
}

func accept(s *LuSocket) {

	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		conn.SetLinger(s.config.Linger)
		conn.SetNoDelay(s.config.IsNoDelay)
		conn.SetKeepAlive(s.config.IsKeepAlive)
		conn.SetKeepAlivePeriod(time.Duration(s.config.KeepInterval))
		conn.SetReadBuffer(s.config.SocketRecvBufferSize)
		conn.SetWriteBuffer(s.config.SocketSendBufferSize)

		lc := &luConn{}
		lc.recvBuffer = rbuffergo.New(s.config.RecvBufferSize, true)
		lc.sendBuffer = rbuffergo.New(s.config.SendBufferSize, true)
		lc.conn = conn

		go handleConn(s, lc)
	}
}

func handleConn(s *LuSocket, lc *luConn) {

}

func client(s *LuSocket) error {
	return nil
}
