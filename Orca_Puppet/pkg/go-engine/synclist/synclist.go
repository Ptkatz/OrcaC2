package synclist

import (
	"container/list"
	"sync"
)

type List struct {
	data *list.List
	lock sync.Mutex
}

func NewList() *List {
	q := new(List)
	q.data = list.New()
	return q
}

func (q *List) Push(v interface{}) {
	defer q.lock.Unlock()
	q.lock.Lock()
	q.data.PushFront(v)
}

func (q *List) Pop() interface{} {
	defer q.lock.Unlock()
	q.lock.Lock()
	iter := q.data.Back()
	if iter == nil {
		return nil
	}
	v := iter.Value
	q.data.Remove(iter)
	return v
}

func (q *List) Len() int {
	defer q.lock.Unlock()
	q.lock.Lock()
	return q.data.Len()
}

func (q *List) Range(f func(value interface{})) {
	defer q.lock.Unlock()
	q.lock.Lock()
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		f(iter.Value)
	}
}

func (q *List) Contain(v interface{}) bool {
	defer q.lock.Unlock()
	q.lock.Lock()
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		if v == iter.Value {
			return true
		}
	}
	return false
}

func (q *List) ContainBy(v interface{}, f func(left interface{}, right interface{}) bool) bool {
	defer q.lock.Unlock()
	q.lock.Lock()
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		if f(v, iter.Value) {
			return true
		}
	}
	return false
}
