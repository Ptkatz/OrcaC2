package conn

import (
	"bytes"
	"context"
	"errors"
	"Orca_Server/pkg/go-engine/common"
	"Orca_Server/pkg/go-engine/group"
	"Orca_Server/pkg/go-engine/rbuffergo"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type HttpConfig struct {
	MaxPacketSize       int
	RecvChanLen         int
	AcceptChanLen       int
	RecvChanPushTimeout int
	BufferSize          int
	MaxRetryNum         int
	CloseWaitTimeoutMs  int
	HBTimeoutMs         int
	MaxMsgIndex         int
}

func DefaultHttpConfig() *HttpConfig {
	return &HttpConfig{
		MaxPacketSize:       1024 * 100,
		RecvChanLen:         128,
		AcceptChanLen:       128,
		RecvChanPushTimeout: 100,
		BufferSize:          1024 * 1024,
		MaxRetryNum:         10,
		CloseWaitTimeoutMs:  5000,
		HBTimeoutMs:         10000,
		MaxMsgIndex:         100,
	}
}

const (
	ProtoConnnect = "connect"
	ProtoData     = "data"
	ProtoClose    = "close"

	ProtoCodeOK   = 200
	ProtoCodeFull = 403
	ProtoCodeFail = 404
)

type RhttpConn struct {
	id            string
	isclose       bool
	info          string
	config        *HttpConfig
	dialer        *httpConnDialer
	listenersonny *httpConnListenerSonny
	listener      *httpConnListener
	cancel        context.CancelFunc
	sendb         *rbuffergo.RBuffergo
	recvb         *rbuffergo.RBuffergo
	closelock     sync.Mutex
}

type httpConnDialer struct {
	wg    *group.Group
	addr  string
	url   string
	index int
	retry int
}

type httpConnListenerSonny struct {
	fwg          *group.Group
	addr         string
	expectIndex  int
	lastRecvTime time.Time
	lastSend     []byte
}

type httpConnListener struct {
	wg           *group.Group
	addr         string
	listenerconn *net.TCPListener
	sonny        sync.Map
	accept       *common.Channel
}

func (c *RhttpConn) Name() string {
	return "http"
}

func (c *RhttpConn) Read(p []byte) (n int, err error) {
	c.checkConfig()

	if c.isclose {
		return 0, errors.New("read closed conn")
	}

	if len(p) <= 0 {
		return 0, errors.New("read empty buffer")
	}

	var wg *group.Group
	if c.dialer != nil {
		wg = c.dialer.wg
	} else if c.listener != nil {
		return 0, errors.New("listener can not be read")
	} else if c.listenersonny != nil {
		wg = c.listenersonny.fwg
	} else {
		return 0, errors.New("empty conn")
	}

	for !c.isclose {
		if c.recvb.Size() <= 0 {
			if wg != nil && wg.IsExit() {
				return 0, errors.New("closed conn")
			}
			time.Sleep(time.Millisecond * 100)
			continue
		}

		size := copy(p, c.recvb.GetReadLineBuffer())
		c.recvb.SkipRead(size)
		return size, nil
	}

	return 0, errors.New("read closed conn")
}

func (c *RhttpConn) Write(p []byte) (n int, err error) {
	c.checkConfig()

	if c.isclose {
		return 0, errors.New("write closed conn")
	}

	if len(p) <= 0 {
		return 0, errors.New("write empty data")
	}

	var wg *group.Group
	if c.dialer != nil {
		wg = c.dialer.wg
	} else if c.listener != nil {
		return 0, errors.New("listener can not be write")
	} else if c.listenersonny != nil {
		wg = c.listenersonny.fwg
	} else {
		return 0, errors.New("empty conn")
	}

	totalsize := len(p)
	cur := 0

	for !c.isclose {
		size := totalsize - cur
		svleft := c.sendb.Capacity() - c.sendb.Size()
		if size > svleft {
			size = svleft
		}

		if size <= 0 {
			if wg != nil && wg.IsExit() {
				return 0, errors.New("closed conn")
			}
			time.Sleep(time.Millisecond * 100)
			continue
		}

		c.sendb.Write(p[cur : cur+size])
		cur += size

		if cur >= totalsize {
			return totalsize, nil
		}

		time.Sleep(time.Millisecond * 100)
	}

	return 0, errors.New("write closed conn")
}

func (c *RhttpConn) Close() error {
	c.checkConfig()

	if c.isclose {
		return nil
	}

	c.closelock.Lock()
	defer c.closelock.Unlock()

	//loggo.Debug("start Close %s", c.Info())

	if c.cancel != nil {
		c.cancel()
	}
	if c.dialer != nil {
		if c.dialer.wg != nil {
			//loggo.Debug("start Close dialer %s", c.Info())
			c.dialer.wg.Stop()
			c.dialer.wg.Wait()
		}
	} else if c.listener != nil {
		if c.listener.listenerconn != nil {
			c.listener.listenerconn.Close()
		}
		if c.listener.wg != nil {
			//loggo.Debug("start Close listener %s", c.Info())
			c.listener.wg.Stop()
			c.listener.sonny.Range(func(key, value interface{}) bool {
				u := value.(*RhttpConn)
				u.Close()
				return true
			})
			c.listener.wg.Wait()
		}
	} else if c.listenersonny != nil {
		//loggo.Debug("start Close listenersonny %s", c.Info())
	}
	c.isclose = true

	//loggo.Debug("Close ok %s", c.Info())

	return nil
}

func (c *RhttpConn) Info() string {
	c.checkConfig()

	if c.info != "" {
		return c.info
	}
	if c.dialer != nil {
		c.info = c.id + "<--rhttp dialer-->" + c.dialer.addr
	} else if c.listener != nil {
		c.info = "rhttp listener--" + c.listener.addr
	} else if c.listenersonny != nil {
		c.info = c.id + "<--rhttp listenersonny-->" + c.listenersonny.addr
	} else {
		c.info = "empty http conn"
	}
	return c.info
}

func (c *RhttpConn) postData(url string, d []byte) (int, []byte, error) {

	data := bytes.NewReader(d)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Close = true

	tp := http.Transport{}
	tp.Dial = func(network, addr string) (net.Conn, error) {
		var d net.Dialer
		if gControlOnConnSetup != nil {
			d = net.Dialer{Control: gControlOnConnSetup}
		}
		return d.Dial(network, addr)
	}

	client := &http.Client{}
	client.Transport = &tp
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}

func (c *RhttpConn) Dial(dst string) (Conn, error) {
	c.checkConfig()

	id := common.UniqueId()

	url := dst + "/" + id

	if !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}

	code, ret, err := c.postData(url+"?type="+ProtoConnnect, []byte{})
	if err != nil {
		return nil, err
	}

	if code != ProtoCodeOK {
		return nil, errors.New("dial fail " + string(ret))
	}

	wg := group.NewGroup("RhttpConn Dialer"+" "+id, nil, nil)

	sendb := rbuffergo.New(c.config.BufferSize, true)
	recvb := rbuffergo.New(c.config.BufferSize, true)

	dialer := &httpConnDialer{wg: wg, url: url, index: 0, retry: 0, addr: dst}

	u := &RhttpConn{id: id, config: c.config, dialer: dialer, sendb: sendb, recvb: recvb}

	wg.Go("RhttpConn updateDialerSonny"+" "+u.Info(), func() error {
		return u.updateDialerSonny()
	})

	return u, nil
}

func (c *RhttpConn) updateDialerSonny() error {

	//loggo.Debug("start http conn %s", c.Info())

	buf := make([]byte, c.config.MaxPacketSize)
	var lastrecv []byte
	var lastsend []byte
	lastrecv = nil
	lastsend = nil
	for !c.dialer.wg.IsExit() {
		active := false

		if lastrecv != nil {
			if !c.recvb.Write(lastrecv) {
				time.Sleep(time.Microsecond * 100)
				continue
			}
			active = true
		}
		lastrecv = nil

		var send []byte
		if lastsend == nil {
			sendn := common.MinOfInt(c.sendb.Size(), len(buf))
			if sendn > 0 {
				if !c.sendb.Read(buf[0:sendn]) {
					//loggo.Error("sendb Read fail")
					return errors.New("sendb Read fail")
				}
				active = true
				send = buf[0:sendn]
			}
		} else {
			send = lastsend
			active = true
		}

		code, ret, err := c.postData(c.dialer.url+"?type="+ProtoData+"&index="+strconv.Itoa(c.dialer.index), send)
		if err != nil || code != ProtoCodeOK {
			if code != ProtoCodeFull {
				c.dialer.retry++
				if c.dialer.retry > c.config.MaxRetryNum {
					//loggo.Error("retry max %d", c.dialer.retry)
					break
				}
			}
			lastsend = send
			time.Sleep(time.Millisecond * 100)
			continue
		}
		lastsend = nil

		//loggo.Debug("dailer send ok %s %d %d %d", c.Info(), c.dialer.index, len(send), len(ret))

		c.dialer.index++
		if c.dialer.index >= c.config.MaxMsgIndex {
			c.dialer.index = 0
		}

		if len(ret) > 0 {
			if !c.recvb.Write(ret) {
				lastrecv = ret
				continue
			}
			active = true
		}

		if !active {
			time.Sleep(time.Microsecond * 100)
		}
	}

	//loggo.Debug("close http conn update %s", c.Info())

	startEndTime := time.Now()
	for !c.dialer.wg.IsExit() {
		now := time.Now()

		diffclose := now.Sub(startEndTime)
		if diffclose > time.Millisecond*time.Duration(c.config.CloseWaitTimeoutMs) {
			break
		}

		if c.recvb.Size() <= 0 {
			break
		}

		time.Sleep(time.Millisecond * 10)
	}

	//loggo.Debug("close http conn %s", c.Info())

	c.postData(c.dialer.url+"?type="+ProtoClose, []byte{})

	return errors.New("closed")
}

func (c *RhttpConn) Listen(dst string) (Conn, error) {
	c.checkConfig()

	addr, err := net.ResolveTCPAddr("tcp", dst)
	if err != nil {
		return nil, err
	}
	listenerconn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	ch := common.NewChannel(c.config.AcceptChanLen)

	wg := group.NewGroup("RhttpConn Listen"+" "+dst, nil, func() {
		listenerconn.Close()
		ch.Close()
	})

	listener := &httpConnListener{
		addr:         dst,
		listenerconn: listenerconn,
		wg:           wg,
		accept:       ch,
	}

	u := &RhttpConn{id: common.UniqueId(), config: c.config, listener: listener}
	wg.Go("RhttpConn Listen loopRecv"+" "+dst, func() error {
		return u.loopRecv()
	})
	wg.Go("RhttpConn Listen checkSonnyClose"+" "+dst, func() error {
		return u.checkSonnyClose()
	})

	return u, nil
}

func (c *RhttpConn) Accept() (Conn, error) {
	c.checkConfig()

	if c.listener.wg == nil {
		return nil, errors.New("not listen")
	}
	for !c.listener.wg.IsExit() {
		s := <-c.listener.accept.Ch()
		if s == nil {
			break
		}
		sonny := s.(*RhttpConn)
		_, ok := c.listener.sonny.Load(sonny.id)
		if !ok {
			continue
		}
		if sonny.isclose {
			continue
		}
		return sonny, nil
	}
	return nil, errors.New("listener close")
}

func (c *RhttpConn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//loggo.Debug("ServeHTTP %v %v", r.Method, r.RequestURI)

	u, err := url.Parse(r.RequestURI)
	if err != nil {
		//loggo.Error("Parse fail %v", r.RequestURI)
		w.WriteHeader(ProtoCodeFail)
		w.Write([]byte("url Parse fail"))
		return
	}

	id := u.Path
	param := u.Query()
	types, ok := param["type"]
	if !ok || len(types) == 0 {
		//loggo.Error("no params type %v", r.RequestURI)
		w.WriteHeader(ProtoCodeFail)
		w.Write([]byte("no params type"))
		return
	}
	ty := types[0]

	v, ok := c.listener.sonny.Load(id)
	if !ok {
		if ty != ProtoConnnect {
			//loggo.Error("no sonny id %v", id)
			w.WriteHeader(ProtoCodeFail)
			w.Write([]byte("no sonny id"))
			return
		}

		sonny := &httpConnListenerSonny{fwg: c.listener.wg, expectIndex: 0, lastRecvTime: time.Now(), addr: c.listener.addr}

		sendb := rbuffergo.New(c.config.BufferSize, true)
		recvb := rbuffergo.New(c.config.BufferSize, true)

		u := &RhttpConn{id: id, config: c.config, listenersonny: sonny, sendb: sendb, recvb: recvb}

		c.listener.sonny.Store(id, u)

		c.listener.accept.Write(u)

		w.WriteHeader(ProtoCodeOK)

	} else {
		u := v.(*RhttpConn)
		u.listenersonny.lastRecvTime = time.Now()

		if ty != ProtoData && ty != ProtoClose {
			//loggo.Error("wrong type %v %v", id, ty)
			w.WriteHeader(ProtoCodeFail)
			w.Write([]byte("wrong type " + ty))
			return
		}

		if ty == ProtoClose {
			c.listener.sonny.Delete(u.id)
			w.WriteHeader(ProtoCodeOK)
			return
		}

		indexs, ok := param["index"]
		if !ok || len(indexs) == 0 {
			//loggo.Error("no index type %v", r.RequestURI)
			w.WriteHeader(ProtoCodeFail)
			w.Write([]byte("no params index"))
			return
		}
		index, err := strconv.Atoi(indexs[0])
		if err != nil {
			//loggo.Error("index fail %v", r.RequestURI)
			w.WriteHeader(ProtoCodeFail)
			w.Write([]byte("index fail"))
			return
		}

		newrecv := true
		if index != u.listenersonny.expectIndex {
			nextindex := index + 1
			if nextindex >= u.config.MaxMsgIndex {
				nextindex = 0
			}
			if nextindex == u.listenersonny.expectIndex {
				newrecv = false
			} else {
				//loggo.Error("index diff %v %v", r.RequestURI, u.listenersonny.expectIndex)
				w.WriteHeader(ProtoCodeFail)
				w.Write([]byte("index diff"))
				return
			}
		}

		if newrecv {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				//loggo.Error("read body fail %v", r.RequestURI)
				w.WriteHeader(ProtoCodeFail)
				w.Write([]byte("read body fail"))
				return
			}

			if !u.recvb.Write(body) {
				//loggo.Debug("body write fail %v %v", r.RequestURI, len(body))
				w.WriteHeader(ProtoCodeFull)
				w.Write([]byte("body write fail"))
				return
			}

			u.listenersonny.expectIndex++
			if u.listenersonny.expectIndex >= u.config.MaxMsgIndex {
				u.listenersonny.expectIndex = 0
			}

			sendn := common.MinOfInt(u.config.MaxPacketSize, u.sendb.Size())
			buff := make([]byte, sendn)
			u.sendb.Read(buff)

			w.WriteHeader(ProtoCodeOK)
			w.Write(buff)

			u.listenersonny.lastSend = buff
		} else {
			w.WriteHeader(ProtoCodeOK)
			w.Write(u.listenersonny.lastSend)
		}
	}
}

func (c *RhttpConn) loopRecv() error {
	c.checkConfig()
	http.Serve(c.listener.listenerconn, c)
	return nil
}

func (c *RhttpConn) checkSonnyClose() error {
	c.checkConfig()
	for !c.listener.wg.IsExit() {
		c.listener.sonny.Range(func(key, value interface{}) bool {
			u := value.(*RhttpConn)
			if u.isclose || time.Now().Sub(u.listenersonny.lastRecvTime) > time.Second*time.Duration(c.config.HBTimeoutMs) {
				c.listener.sonny.Delete(key)
			}
			return true
		})
		time.Sleep(time.Second)
	}
	return nil
}

func (c *RhttpConn) checkConfig() {
	if c.config == nil {
		c.config = DefaultHttpConfig()
	}
}

func (c *RhttpConn) SetConfig(config *HttpConfig) {
	c.config = config
}

func (c *RhttpConn) GetConfig() *HttpConfig {
	c.checkConfig()
	return c.config
}
