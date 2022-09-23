package tmap

import (
	"fmt"
	"testing"
	"time"
)

func Test0001(t *testing.T) {
	tt := NewTMap()
	tt.Add(1, "1", 1000)
	tt.Add(2, "2", 2000)
	time.Sleep(1500 * time.Millisecond)
	fmt.Println(tt.Get(1))
	fmt.Println(tt.Get(2))
	time.Sleep(1500 * time.Millisecond)
	fmt.Println(tt.Get(2))
}
