package synclist

import (
	"fmt"
	"testing"
)

func Test0001(t *testing.T) {
	tt := NewList()
	tt.Push(1)
	tt.Push(11)
	tt.Push(111)
	fmt.Println(tt.Len())
	fmt.Println(tt.Pop())
	fmt.Println(tt.Len())
}
