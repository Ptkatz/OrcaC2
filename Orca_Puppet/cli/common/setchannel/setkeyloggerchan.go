package setchannel

func AddKeyloggerQuitSignChan(id string, m chan interface{}) {
	mutex.Lock()
	KeyloggerQuitSignChan[id] = m
	mutex.Unlock()
}

func GetKeyloggerQuitSignChan(id string) (m chan interface{}, exist bool) {
	mutex.Lock()
	m, exist = KeyloggerQuitSignChan[id]
	mutex.Unlock()
	return
}

func DeleteKeyloggerQuitSignChan(id string) {
	mutex.Lock()
	if m, ok := KeyloggerQuitSignChan[id]; ok {
		close(m)
		delete(KeyloggerQuitSignChan, id)
	}
	mutex.Unlock()
}
