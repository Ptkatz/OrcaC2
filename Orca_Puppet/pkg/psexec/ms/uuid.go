package ms

// 此文件提供常用服务uuid

const (
	SRVSVC_UUID    = "4b324fc8-1670-01d3-1278-5a47bf6ee188"
	SRVSVC_VERSION = 2
	NTSVCS_UUID    = "367abb81-9844-35f1-ad32-98f038001003"
	NTSVCS_VERSION = 2
)

var UUIDMap = map[string]string{
	SRVSVC_UUID: "\\PIPE\\srvsvc",
	NTSVCS_UUID: "\\PIPE\\ntsvcs",
}
