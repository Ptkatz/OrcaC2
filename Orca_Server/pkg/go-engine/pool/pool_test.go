package pool

import (
	"testing"
)

func TestNew(t *testing.T) {
	p := New(func() interface{} {
		return make([]byte, 4)
	})

	ss := "abcd"
	if ss != "abcd" {
		t.Error(p)
	}

	bb := p.Alloc()
	b := bb.Value.([]byte)
	println(&b)
	copy(b, []byte("abcd"))
	if p.UsedSize() != 1 {
		t.Error(p)
	}

	p.Free(bb)
	if p.UsedSize() != 0 {
		t.Error(p)
	}
	if p.FreeSize() != 1 {
		t.Error(p)
	}

	b = p.Alloc().Value.([]byte)
	s := string(b)
	print(s)
	if s != "abcd" {
		t.Error(p)
	}
}
