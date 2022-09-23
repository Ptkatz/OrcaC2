package common

import "time"

type Channel struct {
	ch     chan interface{}
	closed bool
}

func NewChannel(len int) *Channel {
	return &Channel{ch: make(chan interface{}, len)}
}

func (c *Channel) Close() {
	defer func() {
		if recover() != nil {
			c.closed = true
		}
	}()

	if !c.closed {
		c.closed = true
		close(c.ch)
	}
}

func (c *Channel) Write(v interface{}) {
	defer func() {
		if recover() != nil {
			c.closed = true
		}
	}()

	if !c.closed {
		c.ch <- v
	}
}

func (c *Channel) WriteTimeout(v interface{}, timeoutms int) bool {
	defer func() {
		if recover() != nil {
			c.closed = true
		}
	}()

	if !c.closed {

		select {
		case c.ch <- v:
			return true
		case <-time.After(time.Duration(timeoutms) * time.Millisecond):
			return false
		}
	}

	return true
}

func (c *Channel) Ch() <-chan interface{} {
	return c.ch
}
