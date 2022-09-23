package frame

import (
	"container/list"
	"Orca_Puppet/pkg/go-engine/common"
	"Orca_Puppet/pkg/go-engine/congestion"
	"Orca_Puppet/pkg/go-engine/loggo"
	"Orca_Puppet/pkg/go-engine/rbuffergo"
	"github.com/golang/protobuf/proto"
	"strconv"
	"sync"
	"time"
)

type FrameStat struct {
	sendDataNum     int
	recvDataNum     int
	sendReqNum      int
	recvReqNum      int
	sendAckNum      int
	recvAckNum      int
	sendDataNumsMap map[int32]int
	recvDataNumsMap map[int32]int
	sendReqNumsMap  map[int32]int
	recvReqNumsMap  map[int32]int
	sendAckNumsMap  map[int32]int
	recvAckNumsMap  map[int32]int
	sendping        int
	sendpong        int
	recvping        int
	recvpong        int
	recvOldNum      int
	recvOutWinNum   int
}

const (
	hbTimeoutSecond = 10
)

type FrameMgr struct {
	frame_max_size int
	frame_max_id   int32
	sendb          *rbuffergo.RBuffergo
	recvb          *rbuffergo.RBuffergo

	sendblock sync.Locker
	recvblock sync.Locker

	recvlock      sync.Locker
	sendlock      sync.Locker
	windowsize    int32
	resend_timems int
	compress      int

	sendwin  *rbuffergo.ROBuffergo
	sendlist *list.List
	sendid   int32

	recvwin  *rbuffergo.ROBuffergo
	recvlist *list.List
	recvid   int32

	close        bool
	remoteclosed bool
	closesend    bool

	lastPingTime int64
	lastPongTime int64
	rttns        int64

	lastSendHBTime   int64
	lastRecvHBTime   int64
	lastRecvDataTime int64

	reqmap map[int32]int64

	connected bool

	fs            *FrameStat
	openstat      int
	lastPrintStat int64

	debugid string

	ct           congestion.Congestion
	ctLastSendId int32
}

func (fm *FrameMgr) SetDebugid(debugid string) {
	fm.debugid = debugid
}

func (fm *FrameMgr) SetCongestion(ct congestion.Congestion) {
	fm.ct = ct
	fm.ct.Init()
}

func NewFrameMgr(frame_max_size int, frame_max_id int, buffersize int, windowsize int, resend_timems int, compress int, openstat int) *FrameMgr {

	sendb := rbuffergo.New(buffersize, false)
	recvb := rbuffergo.New(buffersize, false)

	fm := &FrameMgr{
		frame_max_size: frame_max_size, frame_max_id: int32(frame_max_id),
		sendb: sendb, recvb: recvb,
		sendblock: &sync.Mutex{}, recvblock: &sync.Mutex{},
		recvlock: &sync.Mutex{}, sendlock: &sync.Mutex{},
		windowsize: int32(windowsize), resend_timems: resend_timems, compress: compress,
		sendwin:  rbuffergo.NewROBuffer(windowsize, 0, frame_max_id),
		sendlist: list.New(), sendid: 0,
		recvwin:  rbuffergo.NewROBuffer(windowsize, 0, frame_max_id),
		recvlist: list.New(), recvid: 0,
		close: false, remoteclosed: false, closesend: false,
		lastPingTime: time.Now().UnixNano(), lastPongTime: time.Now().UnixNano(),
		lastSendHBTime: time.Now().UnixNano(), lastRecvHBTime: time.Now().UnixNano(), lastRecvDataTime: time.Now().UnixNano(),
		rttns:     (int64)(resend_timems * 1000),
		reqmap:    make(map[int32]int64),
		connected: false, openstat: openstat, lastPrintStat: time.Now().UnixNano(),
	}

	if openstat > 0 {
		fm.resetStat()
	}
	return fm
}

func (fm *FrameMgr) GetSendBufferLeft() int {
	fm.sendblock.Lock()
	defer fm.sendblock.Unlock()
	left := fm.sendb.Capacity() - fm.sendb.Size()
	return left
}

func (fm *FrameMgr) WriteSendBuffer(data []byte) {
	fm.sendblock.Lock()
	defer fm.sendblock.Unlock()
	fm.sendb.Write(data)
	//loggo.Debug("debugid %v WriteSendBuffer %v %v", fm.debugid, fm.sendb.Size(), len(data))
}

func (fm *FrameMgr) Update() bool {
	cur := time.Now().UnixNano()

	fm.cutSendBufferToWindow(cur)

	tmpreq, tmpack, tmpackto := fm.preProcessRecvList()
	avtive := len(tmpreq) + len(tmpack) + len(tmpackto)
	fm.processRecvList(tmpreq, tmpack, tmpackto)

	fm.combineWindowToRecvBuffer(cur)

	fm.calSendList(cur)

	fm.ping()
	fm.hb()

	fm.second(cur)

	return avtive > 0
}

func (fm *FrameMgr) cutSendBufferToWindow(cur int64) {

	fm.sendblock.Lock()
	defer fm.sendblock.Unlock()

	sendall := false

	if fm.sendb.Size() < fm.frame_max_size {
		sendall = true
	}

	for fm.sendb.Size() >= fm.frame_max_size && fm.sendwin.Size() < int(fm.windowsize) {
		fd := &FrameData{Type: (int32)(FrameData_USER_DATA),
			Data: make([]byte, fm.frame_max_size)}
		fm.sendb.Read(fd.Data)

		if fm.compress > 0 && len(fd.Data) > fm.compress {
			newb := common.CompressData(fd.Data)
			if len(newb) < len(fd.Data) {
				fd.Data = newb
				fd.Compress = true
			}
		}

		f := &Frame{Type: (int32)(Frame_DATA),
			Id:   fm.sendid,
			Data: fd}

		fm.sendid++
		if fm.sendid >= fm.frame_max_id {
			fm.sendid = 0
		}

		err := fm.sendwin.Set(int(f.Id), f)
		if err != nil {
			loggo.Error("sendwin Set fail %v", err)
		}
		//loggo.Debug("debugid %v cut frame push to send win %v %v %v", fm.debugid, f.Id, fm.frame_max_size, fm.sendwin.Size())
	}

	if sendall && fm.sendb.Size() > 0 && fm.sendwin.Size() < int(fm.windowsize) {
		fd := &FrameData{Type: (int32)(FrameData_USER_DATA),
			Data: make([]byte, fm.sendb.Size())}
		fm.sendb.Read(fd.Data)

		if fm.compress > 0 && len(fd.Data) > fm.compress {
			newb := common.CompressData(fd.Data)
			if len(newb) < len(fd.Data) {
				fd.Data = newb
				fd.Compress = true
			}
		}

		f := &Frame{Type: (int32)(Frame_DATA),
			Id:   fm.sendid,
			Data: fd}

		fm.sendid++
		if fm.sendid >= fm.frame_max_id {
			fm.sendid = 0
		}

		err := fm.sendwin.Set(int(f.Id), f)
		if err != nil {
			loggo.Error("sendwin Set fail %v", err)
		}
		//loggo.Debug("debugid %v cut small frame push to send win %v %v %v", fm.debugid, f.Id, len(f.Data.Data), fm.sendwin.Size())
	}

	if fm.sendb.Empty() && fm.close && !fm.closesend && fm.sendwin.Size() < int(fm.windowsize) {
		fd := &FrameData{Type: (int32)(FrameData_CLOSE)}

		f := &Frame{Type: (int32)(Frame_DATA),
			Id:   fm.sendid,
			Data: fd}

		fm.sendid++
		if fm.sendid >= fm.frame_max_id {
			fm.sendid = 0
		}

		err := fm.sendwin.Set(int(f.Id), f)
		if err != nil {
			loggo.Error("sendwin Set fail %v", err)
		}
		fm.closesend = true
		//loggo.Debug("debugid %v close frame push to send win %v %v", fm.debugid, f.Id, fm.sendwin.Size())
	}
}

func (fm *FrameMgr) calSendList(cur int64) {

	for e := fm.sendwin.FrontInter(); e != nil; e = e.Next() {
		f := e.Value.(*Frame)
		if fm.ct != nil && f.Id < fm.ctLastSendId {
			continue
		}
		if !f.Acked && (f.Resend || cur-f.Sendtime > int64(fm.resend_timems*(int)(time.Millisecond))) &&
			cur-f.Sendtime > fm.rttns {
			if fm.ct != nil && f.Data != nil && len(f.Data.Data) > 0 && !fm.ct.CanSend(int(f.Id), len(f.Data.Data)) {
				fm.ctLastSendId = f.Id
				return
			}
			f.Sendtime = cur
			fm.sendFrame(f)
			f.Resend = false
			if fm.openstat > 0 {
				fm.fs.sendDataNum++
				fm.fs.sendDataNumsMap[f.Id]++
			}
			//loggo.Debug("debugid %v push frame to sendlist %v %v", fm.debugid, f.Id, len(f.Data.Data))
		}
	}
	fm.ctLastSendId = -1
}

func (fm *FrameMgr) GetSendList() *list.List {
	fm.sendlock.Lock()
	defer fm.sendlock.Unlock()
	ret := list.New()
	for e := fm.sendlist.Front(); e != nil; e = e.Next() {
		f := e.Value.(*Frame)
		ret.PushBack(f)
	}
	fm.sendlist.Init()
	return ret
}

func (fm *FrameMgr) sendFrame(f *Frame) {
	fm.sendlock.Lock()
	defer fm.sendlock.Unlock()
	fm.sendlist.PushBack(f)
}

func (fm *FrameMgr) OnRecvFrame(f *Frame) {
	fm.recvlock.Lock()
	defer fm.recvlock.Unlock()
	fm.recvlist.PushBack(f)
}

func (fm *FrameMgr) preProcessRecvList() (map[int32]int, map[int32]int, map[int32]*Frame) {
	fm.recvlock.Lock()
	defer fm.recvlock.Unlock()

	tmpreq := make(map[int32]int)
	tmpack := make(map[int32]int)
	tmpackto := make(map[int32]*Frame)
	for e := fm.recvlist.Front(); e != nil; e = e.Next() {
		f := e.Value.(*Frame)
		if f.Type == (int32)(Frame_REQ) {
			for _, id := range f.Dataid {
				tmpreq[id]++
				//loggo.Debug("debugid %v recv req %v %v", fm.debugid, f.Id, common.Int32ArrayToString(f.Dataid, ","))
			}
		} else if f.Type == (int32)(Frame_ACK) {
			for _, id := range f.Dataid {
				tmpack[id]++
				//loggo.Debug("debugid %v recv ack %v %v", fm.debugid, f.Id, common.Int32ArrayToString(f.Dataid, ","))
			}
		} else if f.Type == (int32)(Frame_DATA) {
			tmpackto[f.Id] = f
			if fm.openstat > 0 {
				fm.fs.recvDataNum++
				fm.fs.recvDataNumsMap[f.Id]++
			}
			//loggo.Debug("debugid %v recv data %v %v", fm.debugid, f.Id, len(f.Data.Data))
		} else if f.Type == (int32)(Frame_PING) {
			fm.processPing(f)
		} else if f.Type == (int32)(Frame_PONG) {
			fm.processPong(f)
		} else {
			loggo.Error("error frame type %v", f.Type)
		}
	}
	fm.recvlist.Init()
	return tmpreq, tmpack, tmpackto
}

func (fm *FrameMgr) processRecvList(tmpreq map[int32]int, tmpack map[int32]int, tmpackto map[int32]*Frame) {

	for id, num := range tmpreq {
		err, value := fm.sendwin.Get(int(id))
		if err != nil {
			loggo.Error("sendwin get id fail %v %v", id, err)
			continue
		}
		if value == nil {
			continue
		}
		f := value.(*Frame)
		if f.Id == id {
			f.Resend = true
			//loggo.Debug("debugid %v choose resend win %v %v", fm.debugid, f.Id, len(f.Data.Data))
		} else {
			loggo.Error("sendwin get id diff %v %v", id, f.Id)
			continue
		}
		if fm.openstat > 0 {
			fm.fs.recvReqNum += num
			fm.fs.recvReqNumsMap[id] += num
		}
	}

	for id, num := range tmpack {
		err, value := fm.sendwin.Get(int(id))
		if err != nil {
			continue
		}
		if value == nil {
			continue
		}
		f := value.(*Frame)
		if f.Id == id {
			f.Acked = true
			//loggo.Debug("debugid %v remove send win %v %v", fm.debugid, f.Id, len(f.Data.Data))
		} else {
			loggo.Error("sendwin get id diff %v %v", id, f.Id)
			continue
		}
		if fm.openstat > 0 {
			fm.fs.recvAckNum += num
			fm.fs.recvAckNumsMap[id] += num
		}
		if fm.ct != nil {
			fm.ct.RecvAck(int(id), len(f.Data.Data))
		}
	}

	for !fm.sendwin.Empty() {
		err, value := fm.sendwin.Front()
		if err != nil {
			loggo.Error("sendwin get Front fail %s", err)
			break
		}
		if value == nil {
			break
		}
		f := value.(*Frame)
		if f.Acked {
			err := fm.sendwin.PopFront()
			if err != nil {
				loggo.Error("sendwin PopFront fail ")
				break
			}
		} else {
			break
		}
	}

	if len(tmpackto) > 0 {
		tmpsize := common.MinOfInt(len(tmpackto), fm.frame_max_size/2/4)
		tmp := make([]int32, len(tmpackto))
		index := 0
		for id, rf := range tmpackto {
			if fm.addToRecvWin(rf) {
				tmp[index] = id
				index++
				if fm.openstat > 0 {
					fm.fs.sendAckNum++
					fm.fs.sendAckNumsMap[id]++
				}
				if index >= tmpsize {
					f := &Frame{Type: (int32)(Frame_ACK), Resend: false, Sendtime: 0,
						Id:     0,
						Dataid: tmp[0:index]}
					fm.sendFrame(f)
					index = 0
					tmp = make([]int32, len(tmpackto))
					//loggo.Debug("debugid %v send ack %v %v", fm.debugid, f.Id, common.Int32ArrayToString(f.Dataid, ","))
				}
				//loggo.Debug("debugid %v add data to win %v %v", fm.debugid, rf.Id, len(rf.Data.Data))
			}
		}
		if index > 0 {
			f := &Frame{Type: (int32)(Frame_ACK), Resend: false, Sendtime: 0,
				Id:     0,
				Dataid: tmp[0:index]}
			fm.sendFrame(f)
			//loggo.Debug("debugid %v send ack %v %v", fm.debugid, f.Id, common.Int32ArrayToString(f.Dataid, ","))
		}
	}
}

func (fm *FrameMgr) addToRecvWin(rf *Frame) bool {

	if !fm.isIdInRange(rf.Id, fm.frame_max_id) {
		//loggo.Debug("debugid %v recv frame not in range %v %v", fm.debugid, rf.Id, fm.recvid)
		if fm.isIdOld(rf.Id, fm.frame_max_id) {
			if fm.openstat > 0 {
				fm.fs.recvOldNum++
			}
			return true
		}
		if fm.openstat > 0 {
			fm.fs.recvOutWinNum++
		}
		return false
	}

	err := fm.recvwin.Set(int(rf.Id), rf)
	if err != nil {
		loggo.Error("recvwin Set fail %v", err)
		return false
	}
	return true
}

func (fm *FrameMgr) processRecvFrame(f *Frame) bool {
	if f.Data.Type == (int32)(FrameData_USER_DATA) {
		fm.recvblock.Lock()
		defer fm.recvblock.Unlock()

		left := fm.recvb.Capacity() - fm.recvb.Size()
		if left >= len(f.Data.Data) {
			src := f.Data.Data
			if f.Data.Compress {
				old, err := common.DeCompressData(src)
				if err != nil {
					loggo.Error("recv frame deCompressData error %v", f.Id)
					return false
				}
				if left < len(old) {
					return false
				}
				//loggo.Debug("debugid %v deCompressData recv frame %v %v %v", fm.debugid, f.Id, len(src), len(old))
				src = old
			}

			fm.lastRecvDataTime = time.Now().UnixNano()

			fm.recvb.Write(src)
			//loggo.Debug("debugid %v combined recv frame to recv buffer %v %v", fm.debugid, f.Id, len(src))
			return true
		}
		return false
	} else if f.Data.Type == (int32)(FrameData_CLOSE) {
		fm.remoteclosed = true
		//loggo.Debug("debugid %v recv remote close frame %v", fm.debugid, f.Id)
		return true
	} else if f.Data.Type == (int32)(FrameData_CONN) {
		fm.sendConnectRsp()
		fm.connected = true
		//loggo.Debug("debugid %v recv remote conn frame %v", fm.debugid, f.Id)
		return true
	} else if f.Data.Type == (int32)(FrameData_CONNRSP) {
		fm.connected = true
		//loggo.Debug("debugid %v recv remote conn rsp frame %v", fm.debugid, f.Id)
		return true
	} else if f.Data.Type == (int32)(FrameData_HB) {
		fm.lastRecvHBTime = time.Now().UnixNano()
		//loggo.Debug("debugid %v recv remote hb frame %v", fm.debugid, f.Id)
		return true
	} else {
		loggo.Error("recv frame type error %v", f.Data.Type)
		return false
	}
}

func (fm *FrameMgr) combineWindowToRecvBuffer(cur int64) {

	for {
		done := false
		err, value := fm.recvwin.Front()
		if err == nil {
			f := value.(*Frame)
			if f.Id == fm.recvid {
				delete(fm.reqmap, f.Id)
				if fm.processRecvFrame(f) {
					err := fm.recvwin.PopFront()
					if err != nil {
						loggo.Error("recvwin PopFront fail %v ", err)
					}
					done = true
					//loggo.Debug("debugid %v process recv frame ok %v %v", fm.debugid, f.Id, len(f.Data.Data))
				}
			}
		}
		if !done {
			//loggo.Debug("debugid %v combined need recvid %v ", fm.debugid, fm.recvid)
			break
		} else {
			fm.recvid++
			if fm.recvid >= fm.frame_max_id {
				fm.recvid = 0
			}
			//loggo.Debug("debugid %v combined ok add recvid %v ", fm.debugid, fm.recvid)
		}
	}

	reqtmp := make(map[int32]int)
	e := fm.recvwin.FrontInter()
	id := fm.recvid
	for len(reqtmp) < int(fm.windowsize) && len(reqtmp)*4 < fm.frame_max_size/2 && e != nil {
		f := e.Value.(*Frame)
		//loggo.Debug("debugid %v start add req id %v %v %v", fm.debugid, fm.recvid, f.Id, id)
		if f.Id != id {
			oldReq := fm.reqmap[f.Id]
			if cur-oldReq > fm.rttns {
				reqtmp[id]++
				fm.reqmap[f.Id] = cur
				//loggo.Debug("debugid %v add req id %v ", fm.debugid, id)
			}
		} else {
			e = e.Next()
		}

		id++
		if id >= fm.frame_max_id {
			id = 0
		}
	}

	if len(reqtmp) > 0 {
		f := &Frame{Type: (int32)(Frame_REQ), Resend: false, Sendtime: 0,
			Id:     0,
			Dataid: make([]int32, len(reqtmp))}
		index := 0
		for id := range reqtmp {
			f.Dataid[index] = id
			index++
			if fm.openstat > 0 {
				fm.fs.sendReqNum++
				fm.fs.sendReqNumsMap[id]++
			}
		}
		fm.sendFrame(f)
		//loggo.Debug("debugid %v send req %v %v", fm.debugid, f.Id, common.Int32ArrayToString(f.Dataid, ","))
	}
}

func (fm *FrameMgr) GetRecvBufferSize() int {
	fm.recvblock.Lock()
	defer fm.recvblock.Unlock()

	return fm.recvb.Size()
}

func (fm *FrameMgr) GetRecvReadLineBuffer() []byte {
	fm.recvblock.Lock()
	defer fm.recvblock.Unlock()

	ret := fm.recvb.GetReadLineBuffer()
	//loggo.Debug("debugid %v GetRecvReadLineBuffer %v %v", fm.debugid, fm.recvb.Size(), len(ret))
	return ret
}

func (fm *FrameMgr) SkipRecvBuffer(size int) {
	fm.recvblock.Lock()
	defer fm.recvblock.Unlock()

	fm.recvb.SkipRead(size)
	//loggo.Debug("debugid %v SkipRead %v %v", fm.debugid, fm.recvb.Size(), size)
}

func (fm *FrameMgr) Close() {
	fm.close = true
}

func (fm *FrameMgr) IsRemoteClosed() bool {
	return fm.remoteclosed
}

func (fm *FrameMgr) ping() {
	cur := time.Now().UnixNano()
	if cur-fm.lastPingTime > (int64)(time.Second) {
		fm.lastPingTime = cur
		f := &Frame{Type: (int32)(Frame_PING), Resend: false, Sendtime: cur,
			Id: 0}
		fm.sendFrame(f)
		//loggo.Debug("debugid %v send ping %v", fm.debugid, cur)
		if fm.openstat > 0 {
			fm.fs.sendping++
		}
	}
}

func (fm *FrameMgr) hb() {
	cur := time.Now().UnixNano()
	if cur-fm.lastSendHBTime > (int64)(time.Second) && fm.sendwin.Size() < int(fm.windowsize) {
		fm.lastSendHBTime = cur

		fd := &FrameData{Type: (int32)(FrameData_HB)}

		f := &Frame{Type: (int32)(Frame_DATA),
			Id:   fm.sendid,
			Data: fd}

		fm.sendid++
		if fm.sendid >= fm.frame_max_id {
			fm.sendid = 0
		}

		err := fm.sendwin.Set(int(f.Id), f)
		if err != nil {
			loggo.Error("sendwin Set fail %v", err)
		}
		//loggo.Debug("debugid %v send hb %v", fm.debugid, f.Id)
	}
}

func (fm *FrameMgr) processPing(f *Frame) {
	rf := &Frame{Type: (int32)(Frame_PONG), Resend: false, Sendtime: f.Sendtime,
		Id: 0}
	fm.sendFrame(rf)
	if fm.openstat > 0 {
		fm.fs.recvping++
		fm.fs.sendpong++
	}
	//loggo.Debug("debugid %v recv ping %v", fm.debugid, f.Sendtime)
}

func (fm *FrameMgr) processPong(f *Frame) {
	cur := time.Now().UnixNano()
	if cur > f.Sendtime {
		rtt := cur - f.Sendtime
		fm.rttns = (fm.rttns + rtt) / 2
		if fm.openstat > 0 {
			fm.fs.recvpong++
		}
		fm.lastPongTime = cur
		//loggo.Debug("debugid %v recv pong %v %dms", fm.debugid, rtt, fm.rttns/1000/1000)
	}
}

func (fm *FrameMgr) isIdInRange(id int32, maxid int32) bool {
	begin := fm.recvid
	end := fm.recvid + fm.windowsize
	if end >= maxid {
		if id >= 0 && id < end-maxid {
			return true
		}
		end = maxid
	}
	if id >= begin && id < end {
		return true
	}
	return false
}

func (fm *FrameMgr) isIdOld(id int32, maxid int32) bool {
	begin := fm.recvid - fm.windowsize
	if begin < 0 {
		begin += maxid
	}
	end := fm.recvid

	if begin < end {
		return id >= begin && id < end
	} else {
		return (id >= begin && id < maxid) || (id >= 0 && id < end)
	}
}

func (fm *FrameMgr) IsConnected() bool {
	return fm.connected
}

func (fm *FrameMgr) Connect() {
	if fm.sendwin.Size() < int(fm.windowsize) {
		fd := &FrameData{Type: (int32)(FrameData_CONN)}

		f := &Frame{Type: (int32)(Frame_DATA),
			Id:   fm.sendid,
			Data: fd}

		fm.sendid++
		if fm.sendid >= fm.frame_max_id {
			fm.sendid = 0
		}

		err := fm.sendwin.Set(int(f.Id), f)
		if err != nil {
			loggo.Error("sendwin Set fail %v", err)
		}
		//loggo.Debug("debugid %v start connect", fm.debugid)
	} else {
		loggo.Error("Connect fail ")
	}
}

func (fm *FrameMgr) sendConnectRsp() {
	if fm.sendwin.Size() < int(fm.windowsize) {
		fd := &FrameData{Type: (int32)(FrameData_CONNRSP)}

		f := &Frame{Type: (int32)(Frame_DATA),
			Id:   fm.sendid,
			Data: fd}

		fm.sendid++
		if fm.sendid >= fm.frame_max_id {
			fm.sendid = 0
		}

		err := fm.sendwin.Set(int(f.Id), f)
		if err != nil {
			loggo.Error("sendwin Set fail %v", err)
		}
		//loggo.Debug("debugid %v send connect rsp", fm.debugid)
	} else {
		loggo.Error("sendConnectRsp fail ")
	}
}

func (fm *FrameMgr) resetStat() {
	fm.fs = &FrameStat{}
	fm.fs.sendDataNumsMap = make(map[int32]int)
	fm.fs.recvDataNumsMap = make(map[int32]int)
	fm.fs.sendReqNumsMap = make(map[int32]int)
	fm.fs.recvReqNumsMap = make(map[int32]int)
	fm.fs.sendAckNumsMap = make(map[int32]int)
	fm.fs.recvAckNumsMap = make(map[int32]int)
}

func (fm *FrameMgr) second(cur int64) {
	if cur-fm.lastPrintStat > (int64)(time.Second) {
		fm.lastPrintStat = cur

		ctinfo := ""
		if fm.ct != nil {
			ctinfo = fm.ct.Info()
		}

		if fm.openstat > 0 {
			fs := fm.fs
			loggo.Info("\nsendDataNum %v\nrecvDataNum %v\nsendReqNum %v\nrecvReqNum %v\nsendAckNum %v\nrecvAckNum %v\n"+
				"sendDataNumsMap %v\nrecvDataNumsMap %v\nsendReqNumsMap %v\nrecvReqNumsMap %v\nsendAckNumsMap %v\nrecvAckNumsMap %v\n"+
				"sendping %v\nrecvping %v\nsendpong %v\nrecvpong %v\n"+
				"sendwin %v\nrecvwin %v\n"+
				"recvOldNum %v\nrecvOutWinNum %v\n"+
				"rtt %v\n"+
				"ct %v\n",
				fs.sendDataNum, fs.recvDataNum,
				fs.sendReqNum, fs.recvReqNum,
				fs.sendAckNum, fs.recvAckNum,
				fm.printStatMap(&fs.sendDataNumsMap), fm.printStatMap(&fs.recvDataNumsMap),
				fm.printStatMap(&fs.sendReqNumsMap), fm.printStatMap(&fs.recvReqNumsMap),
				fm.printStatMap(&fs.sendAckNumsMap), fm.printStatMap(&fs.recvAckNumsMap),
				fs.sendping, fs.recvping,
				fs.sendpong, fs.recvpong,
				fm.sendwin.Size(), fm.recvwin.Size(),
				fs.recvOldNum, fs.recvOutWinNum,
				time.Duration(fm.rttns).String(),
				ctinfo)
			fm.resetStat()
		}

		if fm.ct != nil {
			fm.ct.Update()
		}

	}
}

func (fm *FrameMgr) printStatMap(m *map[int32]int) string {
	tmp := make(map[int]int)
	for _, v := range *m {
		tmp[v]++
	}
	max := 0
	for k := range tmp {
		if k > max {
			max = k
		}
	}
	var ret string
	for i := 1; i <= max; i++ {
		ret += strconv.Itoa(i) + "->" + strconv.Itoa(tmp[i]) + ","
	}
	if len(ret) <= 0 {
		ret = "none"
	}
	return ret
}

func (fm *FrameMgr) MarshalFrame(f *Frame) ([]byte, error) {
	resend := f.Resend
	sendtime := f.Sendtime
	mb, err := proto.Marshal(f)
	f.Resend = resend
	f.Sendtime = sendtime
	return mb, err
}

func (fm *FrameMgr) IsHBTimeout() bool {
	now := time.Now().UnixNano()
	if now-fm.lastRecvHBTime > int64(time.Second)*hbTimeoutSecond && now-fm.lastRecvDataTime > int64(time.Second)*hbTimeoutSecond {
		return true
	}
	return false
}
