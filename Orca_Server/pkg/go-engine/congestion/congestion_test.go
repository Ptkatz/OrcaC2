package congestion

import (
	"fmt"
	"Orca_Server/pkg/go-engine/rbuffergo"
	"testing"
	"time"
)

type CongestionTest struct {
	buf         *rbuffergo.Rlistgo
	rtt_ms      int
	bw_mbps     int
	packet_size int
	last        time.Time
	send        int
}

func (ct *CongestionTest) SendMsg() {
	ct.buf.PushBack(time.Now())
}

func (ct *CongestionTest) RecvMsg() bool {
	err, v := ct.buf.Front()
	if err != nil {
		return false
	}
	t := v.(time.Time)
	if time.Now().Sub(t) < time.Millisecond*time.Duration(ct.rtt_ms) {
		return false
	}
	if time.Now().Sub(ct.last) >= time.Millisecond*100 {
		ct.last = time.Now()
		ct.send = 0
	}
	if ct.send > ct.bw_mbps*1024*1024/10 {
		return false
	}
	ct.send += ct.packet_size
	ct.buf.PopFront()
	return true
}

func (ct *CongestionTest) Start() {
	ct.last = time.Now()
}

func TestBB1(t *testing.T) {
	bb := BBCongestion{}
	bb.Init()

	ct := CongestionTest{}
	ct.rtt_ms = 200
	ct.bw_mbps = 10
	ct.packet_size = 500
	ct.buf = rbuffergo.NewRList(ct.bw_mbps * 1024 * 1024)
	ct.Start()

	send_mbps := 0
	recv_mbps := 0
	last := time.Now()
	lastupdate := time.Now()
	begin := time.Now()
	start := time.Now()
	for {
		for {
			if !bb.CanSend(0, ct.packet_size) {
				break
			}
			ct.SendMsg()
			send_mbps += ct.packet_size
		}

		if time.Now().Sub(lastupdate) > time.Second {
			lastupdate = time.Now()
			//fmt.Printf("begin maxfly %d flyeddata %d flyingdata %d status %d\n", bb.maxfly, bb.flyeddata, bb.flyingdata, bb.status)
			bb.Update()
			//fmt.Printf("end maxfly %d flyeddata %d flyingdata %d status %d\n", bb.maxfly, bb.flyeddata, bb.flyingdata, bb.status)
		}

		for {
			b := ct.RecvMsg()
			if !b {
				break
			}
			bb.RecvAck(0, ct.packet_size)
			recv_mbps += ct.packet_size
		}

		if time.Now().Sub(last) >= time.Second {
			last = time.Now()
			fmt.Printf("send %d MB/s recv %d MB/s %d/%d \n", send_mbps/1024/1024, recv_mbps/1024/1024, ct.buf.Size(), ct.buf.Capacity())
			recv_mbps = 0
			send_mbps = 0
		}

		time.Sleep(10)

		if time.Now().Sub(begin) >= 30*time.Second {
			begin = time.Now()
			fmt.Printf("change from %d to %d \n", ct.bw_mbps, ct.bw_mbps*2)
			ct.bw_mbps = ct.bw_mbps * 2
		}

		if time.Now().Sub(start) >= 60*time.Second {
			break
		}
	}
}

func TestBB2(t *testing.T) {
	bb := BBCongestion{}
	bb.Init()

	ct := CongestionTest{}
	ct.rtt_ms = 200
	ct.bw_mbps = 10
	ct.packet_size = 500
	ct.buf = rbuffergo.NewRList(ct.bw_mbps * 1024 * 1024)
	ct.Start()

	send_mbps := 0
	recv_mbps := 0
	last := time.Now()
	lastupdate := time.Now()
	begin := time.Now()
	start := time.Now()
	for {
		if time.Now().Sub(begin) >= 10*time.Second {
			if time.Now().Sub(begin) >= 30*time.Second {
				begin = time.Now()
				fmt.Printf("pause sending \n")
			}
		} else {
			for {
				if !bb.CanSend(0, ct.packet_size) {
					break
				}
				ct.SendMsg()
				send_mbps += ct.packet_size
			}
		}

		if time.Now().Sub(lastupdate) > time.Second {
			lastupdate = time.Now()
			//fmt.Printf("begin maxfly %d flyeddata %d flyingdata %d status %d\n", bb.maxfly, bb.flyeddata, bb.flyingdata, bb.status)
			bb.Update()
			//fmt.Printf("end maxfly %d flyeddata %d flyingdata %d status %d\n", bb.maxfly, bb.flyeddata, bb.flyingdata, bb.status)
		}

		for {
			b := ct.RecvMsg()
			if !b {
				break
			}
			bb.RecvAck(0, ct.packet_size)
			recv_mbps += ct.packet_size
		}

		if time.Now().Sub(last) >= time.Second {
			last = time.Now()
			fmt.Printf("send %d MB/s recv %d MB/s %d/%d \n", send_mbps/1024/1024, recv_mbps/1024/1024, ct.buf.Size(), ct.buf.Capacity())
			recv_mbps = 0
			send_mbps = 0
		}

		time.Sleep(10)

		if time.Now().Sub(start) >= 60*time.Second {
			break
		}
	}
}
