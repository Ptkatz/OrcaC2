package frame

import (
	"fmt"
	"Orca_Server/pkg/go-engine/common"
	"Orca_Server/pkg/go-engine/rbuffergo"
	"testing"
)

func Test0001(t *testing.T) {

	fm := FrameMgr{}
	fm.frame_max_id = 100000
	fm.frame_max_size = 888
	fm.recvid = 4
	fm.windowsize = 100
	lr := &Frame{}
	rr := &Frame{}
	lr.Id = 1
	rr.Id = 4

	lr.Id = 99
	rr.Id = 8

	fm.recvid = 9000
	lr.Id = 9998
	rr.Id = 9999

	fm.recvid = 9000
	lr.Id = 9998
	rr.Id = 8

	fm.recvid = 0
	lr.Id = 9998
	rr.Id = 8

	fm.recvid = 0
	fm.windowsize = 5
	fmt.Println("fm.isIdInRange  = ", fm.isIdInRange(4, 10))

	fm.recvid = 0
	fm.windowsize = 5
	fmt.Println("fm.isIdInRange  = ", fm.isIdInRange(5, 10))

	fm.recvid = 4
	fm.windowsize = 5
	fmt.Println("fm.isIdInRange  = ", fm.isIdInRange(1, 10))

	fm.recvid = 7
	fm.windowsize = 5
	fmt.Println("fm.isIdInRange  = ", fm.isIdInRange(1, 10))

	fm.recvid = 7
	fm.windowsize = 5
	fmt.Println("fm.isIdInRange  = ", fm.isIdInRange(2, 10))

	fm.recvid = 7
	fm.windowsize = 5
	fmt.Println("fm.isIdInRange  = ", fm.isIdInRange(9, 10))

	fm.recvid = 10
	fm.windowsize = 10000
	fmt.Println("fm.isIdInRange  = ", fm.isIdInRange(0, fm.frame_max_id))

	fm.recvid = 7
	fm.windowsize = 5
	fmt.Println("fm.isIdOld  = ", fm.isIdOld(2, 10))

	fm.recvid = 7
	fm.windowsize = 5
	fmt.Println("fm.isIdOld  = ", fm.isIdOld(1, 10))

	fm.recvid = 3
	fm.windowsize = 5
	fmt.Println("fm.isIdOld  = ", fm.isIdOld(1, 10))

	fm.recvid = 13
	fm.windowsize = 10000
	fmt.Println("fm.isIdOld  = ", fm.isIdOld(9, fm.frame_max_id))

	dd := common.CompressData(([]byte)("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	fmt.Println("fm.compressData  = ", len(dd))

	_, ddd := common.DeCompressData(dd)
	fmt.Println("fm.deCompressData  = ", ddd)

	mm := make(map[int32]int)
	mm[1] = 1
	mm[2] = 1
	mm[3] = 1
	mm[4] = 2
	mm[6] = 7
	mms := fm.printStatMap(&mm)
	fmt.Println("fm.printStatMap  = ", mms)
	fm.openstat = 1
	fm.resetStat()
	fm.sendwin = rbuffergo.NewROBuffer(100, 0, 10000)
	fm.recvwin = rbuffergo.NewROBuffer(100, 0, 10000)
	//fm.printStat(time.Now().UnixNano())
}
