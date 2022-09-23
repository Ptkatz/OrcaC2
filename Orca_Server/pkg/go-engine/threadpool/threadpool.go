package threadpool

import (
	"Orca_Server/pkg/go-engine/common"
	"sync"
	"time"
)

type ThreadPool struct {
	workResultLock sync.WaitGroup
	max            int
	exef           func(interface{})
	ca             []chan interface{}
	control        chan int
	stat           ThreadPoolStat
}

type ThreadPoolStat struct {
	Datalen    []int
	Pushnum    []int
	Processnum []int
}

func NewThreadPool(max int, buffer int, exef func(interface{})) *ThreadPool {
	ca := make([]chan interface{}, max)
	control := make(chan int, max)
	for index, _ := range ca {
		ca[index] = make(chan interface{}, buffer)
	}

	stat := ThreadPoolStat{}
	stat.Datalen = make([]int, max)
	stat.Pushnum = make([]int, max)
	stat.Processnum = make([]int, max)

	tp := &ThreadPool{max: max, exef: exef, ca: ca, control: control, stat: stat}

	for index, _ := range ca {
		go tp.run(index)
	}

	return tp
}

func (tp *ThreadPool) AddJob(hash int, v interface{}) {
	tp.ca[common.AbsInt(hash)%len(tp.ca)] <- v
	tp.stat.Pushnum[common.AbsInt(hash)%len(tp.ca)]++
}

func (tp *ThreadPool) AddJobTimeout(hash int, v interface{}, timeoutms int) bool {
	select {
	case tp.ca[common.AbsInt(hash)%len(tp.ca)] <- v:
		tp.stat.Pushnum[common.AbsInt(hash)%len(tp.ca)]++
		return true
	case <-time.After(time.Duration(timeoutms) * time.Millisecond):
		return false
	}
}

func (tp *ThreadPool) Stop() {
	for i := 0; i < tp.max; i++ {
		tp.control <- i
	}
	tp.workResultLock.Wait()
}

func (tp *ThreadPool) run(index int) {
	defer common.CrashLog()

	tp.workResultLock.Add(1)
	defer tp.workResultLock.Done()

	for {
		select {
		case <-tp.control:
			return
		case v := <-tp.ca[index]:
			tp.exef(v)
			tp.stat.Processnum[index]++
		}
	}
}

func (tp *ThreadPool) GetStat() ThreadPoolStat {
	for index, _ := range tp.ca {
		tp.stat.Datalen[index] = len(tp.ca[index])
	}
	return tp.stat
}

func (tp *ThreadPool) ResetStat() {
	for index, _ := range tp.ca {
		tp.stat.Pushnum[index] = 0
		tp.stat.Processnum[index] = 0
	}
}
