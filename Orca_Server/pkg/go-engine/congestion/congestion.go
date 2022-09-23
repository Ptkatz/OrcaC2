package congestion

type Congestion interface {
	Init()
	RecvAck(id int, size int)
	CanSend(id int, size int) bool
	Update()
	Info() string
}
