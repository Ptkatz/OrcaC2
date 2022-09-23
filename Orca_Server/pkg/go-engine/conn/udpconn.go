package conn

import (
	"context"
	"errors"
	"Orca_Server/pkg/go-engine/common"
	"Orca_Server/pkg/go-engine/group"
	"Orca_Server/pkg/go-engine/loggo"
	"net"
	"sync"
)

type UdpConn struct {
	info          string
	config        *UdpConfig
	dialer        *udpConnDialer
	listenersonny *udpConnListenerSonny
	listener      *udpConnListener
	cancel        context.CancelFunc
}

type udpConnDialer struct {
	conn *net.UDPConn
}

type udpConnListenerSonny struct {
	dstaddr    *net.UDPAddr
	fatherconn *net.UDPConn
	recvch     *common.Channel
	isclose    bool
}

type udpConnListener struct {
	listenerconn *net.UDPConn
	wg           *group.Group
	sonny        sync.Map
	accept       *common.Channel
}

type UdpConfig struct {
	MaxPacketSize       int
	RecvChanLen         int
	AcceptChanLen       int
	RecvChanPushTimeout int
}

func DefaultUdpConfig() *UdpConfig {
	return &UdpConfig{
		MaxPacketSize:       10240,
		RecvChanLen:         128,
		AcceptChanLen:       128,
		RecvChanPushTimeout: 100,
	}
}

func (c *UdpConn) Name() string {
	return "udp"
}

func (c *UdpConn) Read(p []byte) (n int, err error) {
	c.checkConfig()

	if c.dialer != nil {
		return c.dialer.conn.Read(p)
	} else if c.listener != nil {
		return 0, errors.New("listener can not be read")
	} else if c.listenersonny != nil {
		if c.listenersonny.isclose {
			return 0, errors.New("read closed conn")
		}
		b := <-c.listenersonny.recvch.Ch()
		if b == nil {
			return 0, errors.New("read closed conn")
		}
		data := b.([]byte)
		if len(data) > len(p) {
			return 0, errors.New("read buffer too small")
		}
		copy(p, data)
		return len(data), nil
	}
	return 0, errors.New("empty conn")
}

func (c *UdpConn) Write(p []byte) (n int, err error) {
	c.checkConfig()

	if c.dialer != nil {
		return c.dialer.conn.Write(p)
	} else if c.listener != nil {
		return 0, errors.New("listener can not be write")
	} else if c.listenersonny != nil {
		if c.listenersonny.isclose {
			return 0, errors.New("write closed conn")
		}
		return c.listenersonny.fatherconn.WriteToUDP(p, c.listenersonny.dstaddr)
	}
	return 0, errors.New("empty conn")
}

func (c *UdpConn) Close() error {
	c.checkConfig()

	if c.cancel != nil {
		c.cancel()
	}
	if c.dialer != nil {
		return c.dialer.conn.Close()
	} else if c.listener != nil {
		c.listener.wg.Stop()
		c.listener.wg.Wait()
		c.listener.sonny.Range(func(key, value interface{}) bool {
			u := value.(*UdpConn)
			u.Close()
			return true
		})
	} else if c.listenersonny != nil {
		c.listenersonny.recvch.Close()
		c.listenersonny.isclose = true
	}
	return nil
}

func (c *UdpConn) Info() string {
	c.checkConfig()

	if c.info != "" {
		return c.info
	}
	if c.dialer != nil {
		c.info = c.dialer.conn.LocalAddr().String() + "<--udp-->" + c.dialer.conn.RemoteAddr().String()
	} else if c.listener != nil {
		c.info = "udp--" + c.listener.listenerconn.LocalAddr().String()
	} else if c.listenersonny != nil {
		c.info = c.listenersonny.fatherconn.LocalAddr().String() + "<--udp-->" + c.listenersonny.dstaddr.String()
	} else {
		c.info = "empty udp conn"
	}
	return c.info
}

func (c *UdpConn) Dial(dst string) (Conn, error) {
	c.checkConfig()

	addr, err := net.ResolveUDPAddr("udp", dst)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	var d net.Dialer
	if gControlOnConnSetup != nil {
		d = net.Dialer{Control: gControlOnConnSetup}
	}
	conn, err := d.DialContext(ctx, "udp", addr.String())
	if err != nil {
		return nil, err
	}
	c.cancel = nil
	dialer := &udpConnDialer{conn: conn.(*net.UDPConn)}
	return &UdpConn{config: c.config, dialer: dialer}, nil
}

func (c *UdpConn) Listen(dst string) (Conn, error) {
	c.checkConfig()

	ipaddr, err := net.ResolveUDPAddr("udp", dst)
	if err != nil {
		return nil, err
	}

	listenerconn, err := net.ListenUDP("udp", ipaddr)
	if err != nil {
		return nil, err
	}

	ch := common.NewChannel(c.config.AcceptChanLen)

	wg := group.NewGroup("UdpConn Listen"+" "+dst, nil, func() {
		listenerconn.Close()
		ch.Close()
	})

	listener := &udpConnListener{
		listenerconn: listenerconn,
		wg:           wg,
		accept:       ch,
	}

	u := &UdpConn{config: c.config, listener: listener}
	wg.Go("UdpConn Listen loopRecv"+" "+dst, func() error {
		return u.loopRecv()
	})

	return u, nil
}

func (c *UdpConn) Accept() (Conn, error) {
	c.checkConfig()

	if c.listener.wg == nil {
		return nil, errors.New("not listen")
	}
	for !c.listener.wg.IsExit() {
		s := <-c.listener.accept.Ch()
		if s == nil {
			break
		}
		sonny := s.(*UdpConn)
		_, ok := c.listener.sonny.Load(sonny.listenersonny.dstaddr.String())
		if !ok {
			continue
		}
		if sonny.listenersonny.isclose {
			continue
		}
		return sonny, nil
	}
	return nil, errors.New("listener close")
}

func (c *UdpConn) loopRecv() error {
	c.checkConfig()

	buf := make([]byte, c.config.MaxPacketSize)
	for !c.listener.wg.IsExit() {
		n, srcaddr, err := c.listener.listenerconn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		data := make([]byte, n)
		copy(data, buf[0:n])
		srcaddrstr := srcaddr.String()

		v, ok := c.listener.sonny.Load(srcaddrstr)
		if !ok {
			sonny := &udpConnListenerSonny{
				dstaddr:    srcaddr,
				fatherconn: c.listener.listenerconn,
				recvch:     common.NewChannel(c.config.RecvChanLen),
			}

			u := &UdpConn{config: c.config, listenersonny: sonny}
			if !u.listenersonny.recvch.WriteTimeout(data, c.config.RecvChanPushTimeout) {
				loggo.Debug("udp conn %s push %d data to %s recv channel timeout", c.Info(), len(data), u.Info())
			}
			c.listener.sonny.Store(srcaddrstr, u)

			c.listener.accept.Write(u)
		} else {
			u := v.(*UdpConn)
			if !u.listenersonny.recvch.WriteTimeout(data, c.config.RecvChanPushTimeout) {
				loggo.Debug("udp conn %s push %d data to %s recv channel timeout", c.Info(), len(data), u.Info())
			}
		}

		c.listener.sonny.Range(func(key, value interface{}) bool {
			u := value.(*UdpConn)
			if u.listenersonny.isclose {
				c.listener.sonny.Delete(key)
			}
			return true
		})
	}
	return nil
}

func (c *UdpConn) checkConfig() {
	if c.config == nil {
		c.config = DefaultUdpConfig()
	}
}

func (c *UdpConn) SetConfig(config *UdpConfig) {
	c.config = config
}

func (c *UdpConn) GetConfig() *UdpConfig {
	c.checkConfig()
	return c.config
}
