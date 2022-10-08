package setchannel

// 添加鼠标动作通道
func AddMouseActionChan(id string, m chan string) {
	mutex.Lock()
	MouseActionChan[id] = m
	mutex.Unlock()
}

// 获取指定鼠标动作通道
func GetMouseActionChan(id string) (m chan string, exist bool) {
	mutex.Lock()
	m, exist = MouseActionChan[id]
	mutex.Unlock()
	return
}

// 删除指定鼠标动作通道
func DeleteMouseActionChan(id string) {
	mutex.Lock()
	if m, ok := MouseActionChan[id]; ok {
		close(m)
		delete(MouseActionChan, id)
	}
	mutex.Unlock()
}

// 添加键盘动作通道
func AddKeyboardActionChan(id string, m chan string) {
	mutex.Lock()
	KeyboardActionChan[id] = m
	mutex.Unlock()
}

// 获取指定键盘动作通道
func GetKeyboardActionChan(id string) (m chan string, exist bool) {
	mutex.Lock()
	m, exist = KeyboardActionChan[id]
	mutex.Unlock()
	return
}

// 删除指定键盘动作通道
func DeleteKeyboardActionChan(id string) {
	mutex.Lock()
	if m, ok := KeyboardActionChan[id]; ok {
		close(m)
		delete(KeyboardActionChan, id)
	}
	mutex.Unlock()
}
