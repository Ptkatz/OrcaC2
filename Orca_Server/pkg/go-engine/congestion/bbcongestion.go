package congestion

import (
	"fmt"
	"Orca_Server/pkg/go-engine/common"
	"Orca_Server/pkg/go-engine/rbuffergo"
	"math"
	"strconv"
	"time"
)

const (
	bbc_status_init = 0
	bbc_status_prop = 1

	bbc_win            = 5
	bbc_maxfly_grow    = 2.1
	bbc_maxfly_compare = float64(1.5)
)

var prop_seq = []float64{1, 1, 1.5, 1}

type BBCongestion struct {
	status        int
	maxfly        int
	flyeddata     int
	lastflyeddata int
	flyingdata    int
	rateflywin    *rbuffergo.Rlistgo
	flyedwin      *rbuffergo.Rlistgo
	propindex     int
	last          time.Time
	lastratewin   float64
	lastflyedwin  int
}

func (bb *BBCongestion) Init() {
	bb.status = bbc_status_init
	bb.maxfly = 1024 * 1024
	bb.rateflywin = rbuffergo.NewRList(bbc_win)
	bb.flyedwin = rbuffergo.NewRList(bbc_win)
	bb.last = time.Now()
}

func (bb *BBCongestion) RecvAck(id int, size int) {
	bb.flyeddata += size
}

func (bb *BBCongestion) CanSend(id int, size int) bool {
	if bb.flyingdata > bb.maxfly {
		return false
	}
	bb.flyingdata += size
	return true
}

func (bb *BBCongestion) Update() {

	if bb.flyeddata <= 0 {
		return
	}

	currate := float64(bb.flyingdata) / float64(bb.flyeddata)
	if currate < 1 {
		currate = 1
	}

	if bb.rateflywin.Full() {
		bb.rateflywin.PopFront()
	}
	bb.rateflywin.PushBack(currate)

	lastratewin := math.MaxFloat64
	for e := bb.rateflywin.FrontInter(); e != nil; e = e.Next() {
		rate := e.Value.(float64)
		if rate < lastratewin {
			lastratewin = rate
		}
	}

	if bb.flyedwin.Full() {
		bb.flyedwin.PopFront()
	}
	bb.flyedwin.PushBack(bb.flyeddata)

	lastflyedwin := 0
	for e := bb.flyedwin.FrontInter(); e != nil; e = e.Next() {
		flyed := e.Value.(int)
		if flyed > lastflyedwin {
			lastflyedwin = flyed
		}
	}

	if bb.status == bbc_status_init {
		if float64(bb.flyeddata) <= bbc_maxfly_compare*float64(bb.lastflyeddata) {
			oldmaxfly := bb.maxfly
			bb.maxfly = int(float64(oldmaxfly) / bbc_maxfly_grow)
			bb.status = bbc_status_prop
			//loggo.Debug("bbc_status_init flyeddata %d maxfly %d change", bb.flyeddata, bb.maxfly)
		} else {
			oldmaxfly := bb.maxfly
			bb.maxfly = int(float64(oldmaxfly) * bbc_maxfly_grow)
			//loggo.Debug("bbc_status_init grow flyeddata %d oldmaxfly %d maxfly %d", bb.flyeddata, oldmaxfly, bb.maxfly)
		}
		bb.lastflyeddata = bb.flyeddata
	} else if bb.status == bbc_status_prop {
		maxfly := float64(lastflyedwin) * lastratewin
		curmaxfly := int(maxfly)
		if curmaxfly > bb.maxfly {
			bb.maxfly = curmaxfly
		} else {
			if common.NearlyEqual(bb.flyingdata, bb.maxfly) {
				bb.maxfly = curmaxfly
			}
		}
		bb.maxfly = int(float64(bb.maxfly) * prop_seq[bb.propindex])
		//loggo.Debug("bbc_status_prop lastflyedwin %v lastrate %v maxfly %d prop %v", lastflyedwin, lastrate, bb.maxfly, prop_seq[bb.propindex])
		bb.propindex++
		bb.propindex = bb.propindex % len(prop_seq)
	} else {
		panic("error status " + strconv.Itoa(bb.status))
	}

	bb.flyeddata = 0
	bb.flyingdata = 0
	bb.lastratewin = lastratewin
	bb.lastflyedwin = lastflyedwin

	if bb.maxfly < 1024*1024 {
		bb.maxfly = 1024 * 1024
	}
}

func (bb *BBCongestion) Info() string {
	return fmt.Sprintf("status %v maxfly %v flyeddata %v lastratewin %v lastflyedwin %v", bb.status, bb.maxfly,
		bb.flyeddata, bb.lastratewin, bb.lastflyedwin)
}
