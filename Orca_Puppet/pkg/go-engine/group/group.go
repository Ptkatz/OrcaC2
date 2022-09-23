package group

import (
	"errors"
	"fmt"
	"Orca_Puppet/pkg/go-engine/common"
	"Orca_Puppet/pkg/go-engine/loggo"
	"sync"
	"sync/atomic"
	"time"
)

type Group struct {
	father   *Group
	son      map[*Group]int
	wg       int32
	errOnce  sync.Once
	err      error
	isexit   bool
	exitfunc func()
	donech   chan int
	sonname  map[string]int
	lock     sync.Mutex
	name     string
}

func NewGroup(name string, father *Group, exitfunc func()) *Group {
	g := &Group{
		father:   father,
		exitfunc: exitfunc,
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	g.donech = make(chan int)
	g.sonname = make(map[string]int)
	g.son = make(map[*Group]int)
	g.name = name

	if father != nil {
		father.addson(g)
	}

	return g
}

func (g *Group) addson(son *Group) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.son[son]++
}

func (g *Group) removeson(son *Group) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.son[son] == 0 {
		//loggo.Debug("removeson fail no son %s %s", g.name, son.name)
	}
	delete(g.son, son)
	if g.son[son] != 0 {
		loggo.Error("removeson fail has son %s %s", g.name, son.name)
	}
}

func (g *Group) add() {
	atomic.AddInt32(&g.wg, 1)
	if g.father != nil {
		g.father.add()
	}
}

func (g *Group) done() {
	atomic.AddInt32(&g.wg, -1)
	if g.father != nil {
		g.father.done()
	}
}

func (g *Group) IsExit() bool {
	return g.isexit
}

func (g *Group) Error() error {
	return g.err
}

func (g *Group) exit(err error) {
	g.errOnce.Do(func() {
		g.err = err
		g.isexit = true
		close(g.donech)
		if g.exitfunc != nil {
			g.exitfunc()
		}

		for son, _ := range g.son {
			son.exit(err)
		}
	})
}

func (g *Group) runningmap() string {
	g.lock.Lock()
	defer g.lock.Unlock()
	ret := ""
	tmp := make(map[string]int)
	for k, v := range g.sonname {
		if v > 0 {
			tmp[k] = v
		}
	}
	ret += fmt.Sprintf("%v", tmp) + "\n"
	for son, _ := range g.son {
		ret += son.runningmap()
	}
	return ret
}

func (g *Group) Done() <-chan int {
	return g.donech
}

func (g *Group) Go(name string, f func() error) {
	g.lock.Lock()
	defer g.lock.Unlock()
	if g.isexit {
		return
	}
	g.add()
	g.sonname[name]++

	go func() {
		defer common.CrashLog()
		defer g.done()

		if err := f(); err != nil {
			g.exit(err)
		}

		g.lock.Lock()
		defer g.lock.Unlock()
		g.sonname[name]--
	}()
}

func (g *Group) Stop() {
	g.exit(errors.New("stop"))
}

func (g *Group) Wait() error {
	last := int64(0)
	begin := int64(0)
	for g.wg != 0 {
		if g.isexit {
			cur := time.Now().Unix()
			if last == 0 {
				last = cur
				begin = cur
			} else {
				if cur-last > 30 {
					last = cur
					loggo.Error("Group Wait too long %s %d %s %v", g.name, g.wg,
						time.Duration((cur-begin)*int64(time.Second)).String(), g.runningmap())
				}
			}
		} else if g.father != nil {
			if g.father.IsExit() {
				g.exit(errors.New("father exit"))
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
	if g.father != nil {
		g.father.removeson(g)
	}
	return g.err
}
