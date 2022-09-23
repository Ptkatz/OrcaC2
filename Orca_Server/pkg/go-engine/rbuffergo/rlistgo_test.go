package rbuffergo

import (
	"fmt"
	"strconv"
	"testing"
)

func TestRList1(t *testing.T) {
	rob := NewRList(3)
	rob.PushBack(1)
	rob.PushBack(2)
	err := rob.PushBack(3)
	if err != nil {
		fmt.Println(err)
	}
	err = rob.PushBack(3)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("size=" + strconv.Itoa(rob.size))

	for e := rob.FrontInter(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	fmt.Println(rob.Front())
	fmt.Println("loop start")
	for e := rob.FrontInter(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	fmt.Println("loop end")
	rob.PopFront()
	fmt.Println("loop start")
	fmt.Println("size=" + strconv.Itoa(rob.size))
	for e := rob.FrontInter(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	fmt.Println("loop end")
	err = rob.PushBack(4)
	if err != nil {
		fmt.Println(err)
	}
	for e := rob.FrontInter(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	fmt.Println(rob.Full())
	fmt.Println(rob.Empty())
	rob.PopFront()
	rob.PopFront()
	rob.PopFront()
	fmt.Println(rob.Empty())
	for e := rob.FrontInter(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	err = rob.PushBack(5)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("loop start")
	fmt.Println("size=" + strconv.Itoa(rob.size))
	for e := rob.FrontInter(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	fmt.Println("loop end")
}
