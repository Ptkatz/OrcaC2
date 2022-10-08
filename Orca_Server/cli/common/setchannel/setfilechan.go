package setchannel

// 添加文件分片消息通道
func AddFileSliceDataChan(id string, m chan interface{}) {
	mutex.Lock()
	FileSliceDataChan[id] = m
	mutex.Unlock()
}

// 获取指定文件分片消息通道
func GetFileSliceDataChan(id string) (m chan interface{}, exist bool) {
	mutex.Lock()
	m, exist = FileSliceDataChan[id]
	mutex.Unlock()
	return
}

// 删除指定文件分片消息通道
func DeleteFileSliceDataChan(id string) {
	mutex.Lock()
	if m, ok := FileSliceDataChan[id]; ok {
		close(m)
		delete(FileSliceDataChan, id)
	}
	mutex.Unlock()
}
