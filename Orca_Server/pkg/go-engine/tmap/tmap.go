package tmap

import (
	"sync"
	"time"
)

type TMapValue struct {
	v         interface{}
	valid     time.Time
	timeoutms int
}

type TMap struct {
	m sync.Map
}

func NewTMap() *TMap {
	return &TMap{}
}

func (t *TMap) Add(k interface{}, v interface{}, timeoutms int) {
	vv := &TMapValue{v: v, valid: time.Now(), timeoutms: timeoutms}
	t.m.Store(k, vv)
}

func (t *TMap) Del(k interface{}) {
	t.m.Delete(k)
}

func (t *TMap) Get(k interface{}) interface{} {
	vv, ok := t.m.Load(k)
	if !ok {
		return nil
	}
	tv := vv.(*TMapValue)
	if time.Now().Sub(tv.valid) > time.Duration(tv.timeoutms)*time.Millisecond {
		return nil
	}
	return tv.v
}

func (t *TMap) Valid(k interface{}) bool {
	vv, ok := t.m.Load(k)
	if !ok {
		return false
	}
	tv := vv.(*TMapValue)
	tv.valid = time.Now()
	return true
}

func (t *TMap) Update() {

	tmp := make(map[interface{}]*TMapValue)
	t.m.Range(func(key, value interface{}) bool {
		tv := value.(*TMapValue)
		tmp[key] = tv
		return true
	})

	now := time.Now()
	for k, tv := range tmp {
		diff := now.Sub(tv.valid)
		if diff > time.Duration(tv.timeoutms)*time.Millisecond {
			t.m.Delete(k)
		}
	}
}
