package smbopt

import "Orca_Puppet/cli/cmdopt/fileopt"

type SmbOption struct {
	Host   string
	User   string
	Pwd    string
	Hash   string
	Domain string
}

type SmbExecStruct struct {
	SmbStruct SmbOption
	Command   string
}

type SmbUploadStruct struct {
	SmbStruct    SmbOption
	FileMetaInfo fileopt.FileMetaInfo
}
