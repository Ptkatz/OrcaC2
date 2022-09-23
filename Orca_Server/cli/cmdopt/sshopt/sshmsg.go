package sshopt

type SshStruct struct {
	Node    string
	SSHHost string
	SSHUser string
	SSHPwd  string
}
type SshRunStruct struct {
	SshStruct SshStruct
	Command   string
}

type FileMetaInfo struct {
	Fid          string // 操作文件ID，随机生成的UUID
	SaveFileName string // 保存的文件名称
	SliceNum     int    // 基础分片数量
	SliceSize    int64  // 基础分片大小
	RemainSize   int64  // 剩余分片大小
	Md5sum       string // 文件md5值
}

type SshUploadStruct struct {
	SshStruct    SshStruct
	FileMetaInfo FileMetaInfo
}

// 请求文件
type RequestFile struct {
	SaveFileName string
	DestFileName string
}

type SshDownloadStruct struct {
	RequestFile RequestFile
	SshStruct   SshStruct
}

type SshTunnelStruct struct {
	SshStruct SshStruct
	Source    string
	Target    string
}

type SshTunnelRecord struct {
	SSHTunnel           *SSHTunnel
	SshTunnelBaseRecord SshTunnelBaseRecord
}

type SshTunnelBaseRecord struct {
	Uid      string
	ClientId string
	Source   string
	Target   string
}

var SshTunnelRecordLists = make([]SshTunnelRecord, 0)
