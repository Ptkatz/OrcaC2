package setchannel

// 添加下一张截图消息通道
func AddNextScreenChan(id string, m chan string) {
	mutex.Lock()
	NextScreenChan[id] = m
	mutex.Unlock()
}

// 获取下一张截图消息通道
func GetNextScreenChan(id string) (m chan string, exist bool) {
	mutex.Lock()
	m, exist = NextScreenChan[id]
	mutex.Unlock()
	return
}

// 删除下一张截图消息通道
func DeleteNextScreenChan(id string) {
	mutex.Lock()
	if m, ok := NextScreenChan[id]; ok {
		close(m)
		delete(NextScreenChan, id)
	}
	mutex.Unlock()
}
