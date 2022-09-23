package rbuffergo

import (
	"errors"
	"fmt"
)

/*
type:		   [1]
iter:	 begin(2)	 end(8)
			|		   |
data:   _ _ 1 2 3 * 5 * _ _ _
buffer: _ _ _ _ _ _ _ _ _ _ _
index:  0 1 2 3 4 5 6 7 8 9 10
type:		   [2]
iter:	  end(2)   begin(7)
			|		 |
data:   7 8 _ _ _ _ _ 4 * 5 6
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
data:   10 * * 13 3 4 5 * 7 8 *
buffer: _ _ _ _ _ _ _ _ _ _ _
index:  0 1 2 3 4 5 6 7 8 9 10
*/

type ROBuffergo struct {
	buffer []interface{}
	flag   []bool
	id     []int
	len    int
	maxid  int
	begin  int
	size   int
}

func NewROBuffer(len int, startid int, maxid int) *ROBuffergo {
	if startid >= maxid {
		return nil
	}
	if len >= maxid {
		return nil
	}
	buffer := &ROBuffergo{}
	buffer.buffer = make([]interface{}, len)
	buffer.flag = make([]bool, len)
	buffer.id = make([]int, len)
	buffer.len = len
	buffer.maxid = maxid
	for i, _ := range buffer.id {
		buffer.id[i] = startid + i
	}
	return buffer
}

func (b *ROBuffergo) Get(id int) (error, data interface{}) {
	if b.begin >= len(b.flag) {
		return fmt.Errorf("not init"), nil
	}
	cur := b.id[b.begin]
	if id >= b.maxid {
		return fmt.Errorf("id out of range %d %d %d", id, cur, b.maxid), nil
	}

	index := 0
	if id < cur {
		index = (b.begin + (id + b.maxid - cur)) % b.len
	} else {
		index = (b.begin + (id - cur)) % b.len
	}

	if b.id[index] != id {
		return fmt.Errorf("set id error %d %d %d", id, b.id[index], index), nil
	}

	if !b.flag[index] {
		return nil, nil
	}

	if b.buffer[index] == nil {
		return errors.New("data is nil"), nil
	}

	return nil, b.buffer[index]
}

func (b *ROBuffergo) Set(id int, data interface{}) error {
	if data == nil {
		return fmt.Errorf("data nil %d ", id)
	}
	if b.begin >= len(b.flag) {
		return fmt.Errorf("not init")
	}
	cur := b.id[b.begin]
	if id >= b.maxid {
		return fmt.Errorf("id out of range %d %d %d", id, cur, b.maxid)
	}

	index := 0
	if id < cur {
		index = (b.begin + (id + b.maxid - cur)) % b.len
	} else {
		index = (b.begin + (id - cur)) % b.len
	}

	if b.id[index] != id {
		return fmt.Errorf("set id error %d %d %d", id, b.id[index], index)
	}
	b.buffer[index] = data
	if !b.flag[index] {
		b.size++
	}
	b.flag[index] = true
	return nil
}

func (b *ROBuffergo) Front() (error, interface{}) {
	if b.begin >= len(b.flag) {
		return errors.New("no init"), nil
	}
	if !b.flag[b.begin] {
		return errors.New("no front data"), nil
	}
	if b.buffer[b.begin] == nil {
		return errors.New("front data is nil"), nil
	}
	return nil, b.buffer[b.begin]
}

func (b *ROBuffergo) PopFront() error {
	if !b.flag[b.begin] {
		return errors.New("no front data")
	}
	if b.buffer[b.begin] == nil {
		return errors.New("front data is nil")
	}
	old := b.begin
	b.begin++
	if b.begin >= b.len {
		b.begin = 0
	}
	cur := b.id[b.begin]
	b.buffer[old] = nil
	b.flag[old] = false
	b.id[old] = cur + b.len - 1
	if b.id[old] >= b.maxid {
		b.id[old] %= b.maxid
	}
	b.size--
	return nil
}

func (b *ROBuffergo) Size() int {
	return b.size
}

func (b *ROBuffergo) Full() bool {
	return b.size == b.len
}

func (b *ROBuffergo) Empty() bool {
	return b.size <= 0
}

type ROBuffergoInter struct {
	startindex int
	index      int
	Value      interface{}
	b          *ROBuffergo
}

func (b *ROBuffergo) FrontInter() *ROBuffergoInter {
	if b.begin >= len(b.flag) {
		return nil
	}
	if !b.flag[b.begin] {
		return nil
	}
	return &ROBuffergoInter{
		startindex: b.begin,
		index:      b.begin,
		Value:      b.buffer[b.begin],
		b:          b,
	}
}

func (bi *ROBuffergoInter) Next() *ROBuffergoInter {
	for {
		bi.index++
		if bi.index >= bi.b.len {
			bi.index %= bi.b.len
		}
		if bi.index == bi.startindex {
			return nil
		}
		if bi.b.flag[bi.index] {
			break
		}
	}
	bi.Value = bi.b.buffer[bi.index]
	return bi
}
