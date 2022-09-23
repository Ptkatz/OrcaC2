package rbuffergo

import "sync"

/*
type:		   [1]
iter:	 begin(2)	 end(8)
			|		   |
data:   _ _ * * * * * * _ _ _
buffer: _ _ _ _ _ _ _ _ _ _ _
index:  0 1 2 3 4 5 6 7 8 9 10
type:		   [2]
iter:	  end(2)   begin(7)
			|		 |
data:   * * _ _ _ _ _ * * * *
buffer: _ _ _ _ _ _ _ _ _ _ _
index:  0 1 2 3 4 5 6 7 8 9 10
type:		   [3]
iter:	  begin(4),end(4)
				|
data:   _ _ _ _ _ _ _ _ _ _ _
buffer: _ _ _ _ _ _ _ _ _ _ _
index:  0 1 2 3 4 5 6 7 8 9 10
type:		   [4]
iter:	  begin(4),end(4)
|				 |
data:   * * * * * * * * * * *
buffer: _ _ _ _ _ _ _ _ _ _ _
index:  0 1 2 3 4 5 6 7 8 9 10
*/

type RBuffergo struct {
	buffer        []byte
	datasize      int
	begin         int
	end           int
	storeDatasize int
	storeBegin    int
	storeEnd      int
	lock          sync.Locker
}

func New(len int, lock bool) *RBuffergo {
	buffer := &RBuffergo{}
	buffer.buffer = make([]byte, len)
	if lock {
		buffer.lock = &sync.Mutex{}
	}
	return buffer
}

func (b *RBuffergo) CanWrite(size int) bool {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}
	return b.datasize+size <= len(b.buffer)
}

func (b *RBuffergo) SkipWrite(size int) {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	if !(b.datasize+size <= len(b.buffer)) {
		return
	}

	b.datasize += size
	b.end += size
	if b.end >= len(b.buffer) {
		b.end -= len(b.buffer)
	}
}

func (b *RBuffergo) Write(data []byte) bool {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	if !(b.datasize+len(data) <= len(b.buffer)) {
		return false
	}
	// [1][3]
	if b.end >= b.begin {
		// 能装下
		if len(b.buffer)-b.end >= len(data) {
			copy(b.buffer[b.end:], data)
		} else {
			copy(b.buffer[b.end:], data[0:len(b.buffer)-b.end])
			copy(b.buffer, data[len(b.buffer)-b.end:])
		}
	} else /*[2]*/ {
		copy(b.buffer[b.end:], data)
	}

	b.datasize += len(data)
	b.end += len(data)
	if b.end >= len(b.buffer) {
		b.end -= len(b.buffer)
	}

	return true
}

func (b *RBuffergo) CanRead(size int) bool {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	return b.datasize >= size
}

func (b *RBuffergo) SkipRead(size int) {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	if !(b.datasize >= size) {
		return
	}

	b.datasize -= size
	b.begin += size
	if b.begin >= len(b.buffer) {
		b.begin -= len(b.buffer)
	}

	if b.lock == nil {
		if b.datasize == 0 {
			b.begin = 0
			b.end = 0
		}
	}
}

func (b *RBuffergo) Read(data []byte) bool {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	if !(b.datasize >= len(data)) {
		return false
	}

	// [2][4]
	if b.begin >= b.end {
		// 能读完
		if len(b.buffer)-b.begin >= len(data) {
			copy(data, b.buffer[b.begin:])
		} else {
			copy(data[0:len(b.buffer)-b.begin], b.buffer[b.begin:])
			copy(data[len(b.buffer)-b.begin:], b.buffer)
		}
	} else /* [1]*/ {
		copy(data, b.buffer[b.begin:])
	}

	b.datasize -= len(data)
	b.begin += len(data)
	if b.begin >= len(b.buffer) {
		b.begin -= len(b.buffer)
	}

	if b.lock == nil {
		if b.datasize == 0 {
			b.begin = 0
			b.end = 0
		}
	}

	return true
}

func (b *RBuffergo) Store() {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	b.storeDatasize = b.datasize
	b.storeBegin = b.begin
	b.storeEnd = b.end
}

func (b *RBuffergo) Restore() {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	b.datasize = b.storeDatasize
	b.begin = b.storeBegin
	b.end = b.storeEnd
}

func (b *RBuffergo) Clear() {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	b.datasize = 0
	b.begin = 0
	b.end = 0
}

func (b *RBuffergo) Size() int {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	return b.datasize
}

func (b *RBuffergo) Capacity() int {
	return len(b.buffer)
}

func (b *RBuffergo) Empty() bool {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	return b.datasize == 0
}

func (b *RBuffergo) Full() bool {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	return b.datasize == len(b.buffer)
}

func (b *RBuffergo) GetReadLineBuffer() []byte {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	if b.datasize < len(b.buffer)-b.begin {
		return b.buffer[b.begin : b.begin+b.datasize]
	} else {
		return b.buffer[b.begin : b.begin+len(b.buffer)-b.begin]
	}
}

func (b *RBuffergo) GetWriteLineBuffer() []byte {
	if b.lock != nil {
		b.lock.Lock()
		defer b.lock.Unlock()
	}

	if len(b.buffer)-b.datasize < len(b.buffer)-b.end {
		return b.buffer[b.end : b.end+len(b.buffer)-b.datasize]
	} else {
		return b.buffer[b.end : b.end+len(b.buffer)-b.end]
	}
}

func (b *RBuffergo) GetBuffer() []byte {
	return b.buffer
}
