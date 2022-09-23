package pool

import (
	"container/list"
)

type PoolElement struct {
	Value interface{}
}

type Pool struct {
	use    map[*PoolElement]int
	free   *list.List
	allocf func() interface{}
}

func New(allocf func() interface{}) *Pool {
	p := &Pool{}
	p.allocf = allocf
	p.use = make(map[*PoolElement]int)
	p.free = list.New()
	return p
}

func (p *Pool) Alloc() *PoolElement {
	if p.free.Len() <= 0 {
		pe := PoolElement{Value: p.allocf()}
		p.free.PushBack(&pe)
	}
	e := p.free.Front()
	pe := e.Value.(*PoolElement)
	p.free.Remove(e)
	p.use[pe]++
	return pe
}

func (p *Pool) Free(pe *PoolElement) {
	if _, ok := p.use[pe]; ok {
		delete(p.use, pe)
		p.free.PushFront(pe)
	}
}

func (p *Pool) UsedSize() int {
	return len(p.use)
}

func (p *Pool) FreeSize() int {
	return p.free.Len()
}
