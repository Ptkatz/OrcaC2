package rbuffergo

import (
	"testing"
)

func TestNew(t *testing.T) {
	rb := New(10, true)
	if rb.Capacity() != 10 {
		t.Error()
	}
}

func TestRBuffer_Write(t *testing.T) {
	rb := New(10, true)
	rb.Write([]byte{1, 2, 3})
	if rb.Size() != 3 {
		t.Error(rb.Size())
	}
	var tmp [3]byte
	rb.Read(tmp[0:])
	if tmp[0] != 1 || tmp[1] != 2 || tmp[2] != 3 {
		t.Error(tmp)
	}
	t.Log(rb.GetBuffer())
	rb.Write([]byte{1, 2, 3})
	rb.Write([]byte{1, 2, 3})
	rb.Write([]byte{1, 2, 3})
	t.Log(rb.GetBuffer())
	if rb.Size() != 9 {
		t.Error(rb.Size())
	}
	if rb.Write([]byte{1, 2, 3}) == true {
		t.Error()
	}
	t.Log(rb.GetBuffer())
	rb.Read(tmp[0:])
	rb.Read(tmp[0:])
	rb.Write([]byte{1, 2, 3})
	if rb.Size() != 6 {
		t.Error(rb.Size())
	}
	t.Log(rb.GetBuffer())

	var tmp1 [6]byte
	rb.Read(tmp1[0:])
	t.Log(tmp1)
}
