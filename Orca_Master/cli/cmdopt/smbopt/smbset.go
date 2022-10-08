package smbopt

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

var MySmb SmbOption

func InitMsOption() {
	MySmb = SmbOption{
		Host:   "",
		User:   "",
		Pwd:    "",
		Hash:   "",
		Domain: "",
	}
}

func SmbSet(host, user, pwd, hash, domain string) {
	MySmb = SmbOption{
		Host:   host,
		User:   user,
		Pwd:    pwd,
		Hash:   hash,
		Domain: domain,
	}
}
