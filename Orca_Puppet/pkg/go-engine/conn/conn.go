package conn

import (
	"errors"
	"Orca_Puppet/pkg/go-engine/common"
	"io"
	"strings"
	"syscall"
)

type Conn interface {
	io.ReadWriteCloser

	Name() string

	Info() string

	Dial(dst string) (Conn, error)

	Listen(dst string) (Conn, error)
	Accept() (Conn, error)
}

func NewConn(proto string) (Conn, error) {
	proto = strings.ToLower(proto)
	if proto == "tcp" {
		return &TcpConn{}, nil
	} else if proto == "udp" {
		return &UdpConn{}, nil
	} else if proto == "rudp" {
		return &RudpConn{}, nil
	} else if proto == "ricmp" {
		return &RicmpConn{id: common.UniqueId()}, nil
	} else if proto == "kcp" {
		return &KcpConn{}, nil
	} else if proto == "quic" {
		return &QuicConn{}, nil
	} else if proto == "rhttp" {
		return &RhttpConn{}, nil
	}
	return nil, errors.New("undefined proto " + proto)
}

func SupportReliableProtos() []string {
	ret := make([]string, 0)
	ret = append(ret, "tcp")
	ret = append(ret, "rudp")
	ret = append(ret, "ricmp")
	ret = append(ret, "kcp")
	ret = append(ret, "quic")
	ret = append(ret, "rhttp")
	return ret
}

func SupportProtos() []string {
	ret := make([]string, 0)
	ret = append(ret, SupportReliableProtos()...)
	ret = append(ret, "udp")
	return ret
}

func HasReliableProto(proto string) bool {
	return common.HasString(SupportReliableProtos(), proto)
}

func HasProto(proto string) bool {
	return common.HasString(SupportProtos(), proto)
}

var gControlOnConnSetup func(network, address string, c syscall.RawConn) error

func RegisterDialerController(fn func(network, address string, c syscall.RawConn) error) {
	gControlOnConnSetup = fn
}
