package rpc

import (
	"errors"
	"Orca_Puppet/pkg/go-engine/common"
	"sync"
	"time"
)

var grpccallMap sync.Map

func NewCall(timeoutms int) *RpcCall {
	c := &RpcCall{
		timeoutms: timeoutms,
		id:        common.Guid(),
		retc:      make(chan int, 1),
	}
	grpccallMap.Store(c.id, c)
	return c
}

func PutRet(id string, ret ...interface{}) {
	v, ok := grpccallMap.Load(id)
	if !ok {
		return
	}
	rc := v.(*RpcCall)
	rc.result = ret
	rc.retc <- 1
}

type RpcCall struct {
	timeoutms int
	id        string
	result    []interface{}
	retc      chan int
}

func (r *RpcCall) Id() string {
	return r.id
}

func (r *RpcCall) Call(f func()) ([]interface{}, error) {
	f()

	select {
	case _ = <-r.retc:
		grpccallMap.Delete(r.id)
		return r.result, nil
	case <-time.After(time.Duration(r.timeoutms) * time.Millisecond):
		grpccallMap.Delete(r.id)
		return nil, errors.New("time out")
	}
}
