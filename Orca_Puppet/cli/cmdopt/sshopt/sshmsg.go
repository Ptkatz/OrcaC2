package sshopt

import "Orca_Puppet/cli/cmdopt/fileopt"

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

type SshUploadStruct struct {
	SshStruct    SshStruct
	FileMetaInfo fileopt.FileMetaInfo
}

type SshDownloadStruct struct {
	RequestFile fileopt.RequestFile
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
