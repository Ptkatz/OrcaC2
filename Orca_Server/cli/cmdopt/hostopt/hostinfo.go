package hostopt

// 命令消息类型
type CmdClass int32

// HostInfo 被控端信息
type HostInfo struct {
	SystemId  string //SystemId
	ClientId  string //主机标识
	Hostname  string //主机名
	Ip        string //上线ip
	ConnPort  string //上线端口
	Privilege string //执行权限
	Os        string //操作系统版本
	Version   string //连接客户端版本
}
