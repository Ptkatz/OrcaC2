package spider

import (
	"fmt"
	"Orca_Puppet/pkg/go-engine/common"
	"Orca_Puppet/pkg/go-engine/loggo"
	"Orca_Puppet/pkg/go-engine/threadpool"
	"sync/atomic"
	"time"
)

type LoopSpiderSlot interface {
	Name() string
	DefaultCur() string
	Crawl(cur string) bool
	NextCur(cur string) string
}

type LoopSpiderStatus struct {
	Cur      string
	CurInt   int
	UsedTime string
	Done     int64
	Speed    string
	Fail     int64
	OK       int64
}

type LoopSpider struct {
	Thread    int
	Buffer    int
	Cur       string
	done      int64
	startTime time.Time
	fail      int64
	ok        int64
	tp        *threadpool.ThreadPool
	lss       LoopSpiderSlot
}

func NewLoopSpider(lss LoopSpiderSlot) *LoopSpider {

	ls := LoopSpider{}
	ls.lss = lss

	err := common.LoadJson(".ls."+lss.Name()+".json", &ls)
	if err != nil {
		loggo.Error("NewLoopSpider LoadJson fail %s %s", lss.Name(), err)
		return nil
	}

	go saveLoopSpiderJson(&ls)

	ls.tp = threadpool.NewThreadPool(ls.Thread, ls.Buffer, func(i interface{}) {
		tmp := i.(string)
		if ls.lss.Crawl(tmp) {
			atomic.AddInt64(&ls.ok, 1)
		} else {
			atomic.AddInt64(&ls.fail, 1)
		}
		atomic.AddInt64(&ls.done, 1)
	})
	if ls.tp == nil {
		loggo.Error("NewLoopSpider NewThreadPool fail %s %s", lss.Name(), err)
		return nil
	}

	go crawlLoopSpider(&ls)

	return &ls
}

func crawlLoopSpider(ls *LoopSpider) {
	defer common.CrashLog()

	if len(ls.Cur) <= 0 {
		ls.Cur = ls.lss.DefaultCur()
	}
	ls.startTime = time.Now()

	for {
		tmp := ls.Cur
		for !ls.tp.AddJobTimeout(int(common.RandInt()), tmp, 1000) {
		}

		ls.Cur = ls.lss.NextCur(ls.Cur)
	}
}

func saveLoopSpiderJson(ls *LoopSpider) {
	defer common.CrashLog()

	for {
		common.SaveJson(".ls."+ls.lss.Name()+".json", &ls)
		time.Sleep(time.Second)
	}
}

func (ls *LoopSpider) GetLoopSpiderStatus() LoopSpiderStatus {

	res := LoopSpiderStatus{}

	res.Cur = ls.Cur
	res.CurInt = common.Hex2Num(ls.Cur, common.FULL_LETTERS)
	res.UsedTime = time.Now().Sub(ls.startTime).String()
	res.Done = ls.done
	res.Fail = ls.fail
	res.OK = ls.ok

	var speed int64
	if int64(time.Now().Sub(ls.startTime))/int64(time.Second) != 0 {
		speed = res.Done / (int64(time.Now().Sub(ls.startTime)) / int64(time.Second))
	}
	res.Speed = fmt.Sprintf("%d/s", speed)

	return res
}
