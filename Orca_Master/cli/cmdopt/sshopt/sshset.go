package sshopt

type SshOption struct {
	Node    string
	SSHHost string
	SSHUser string
	SSHPwd  string
}

var MySsh SshOption

func InitSshOption(node string) {
	MySsh = SshOption{
		Node:    node,
		SSHHost: "",
		SSHUser: "",
		SSHPwd:  "",
	}
}

func SshSet(node, sshHost, sshUser, sshPwd string) {
	MySsh = SshOption{
		Node:    node,
		SSHHost: sshHost,
		SSHUser: sshUser,
		SSHPwd:  sshPwd,
	}
}
