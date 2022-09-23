package sshopt

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
	"time"
)

type Cli struct {
	user       string
	pwd        string
	ip         string
	port       string
	SshClient  *ssh.Client
	SftpClient *sftp.Client
}

func NewSSHClient(user, pwd, ip, port string) Cli {
	return Cli{
		user: user,
		pwd:  pwd,
		ip:   ip,
		port: port,
	}
}

// 不使用 HostKey， 使用密码
func (c *Cli) getConfig_nokey() *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: c.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.pwd),
		},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config
}

func (c *Cli) Connect() error {
	config := c.getConfig_nokey()
	client, err := ssh.Dial("tcp", c.ip+":"+c.port, config)
	if err != nil {
		return fmt.Errorf("connect server error: %w", err)
	}
	sftp, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("new sftp client error: %w", err)
	}

	c.SshClient = client
	c.SftpClient = sftp
	return nil
}

func (c Cli) Run(cmd string) (string, error) {
	if c.SshClient == nil {
		if err := c.Connect(); err != nil {
			return "", err
		}
	}
	session, err := c.SshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("create new session error: %w", err)
	}
	defer session.Close()

	buf, err := session.CombinedOutput(cmd)
	return string(buf), err
}

func (c Cli) LStateFile(remoteFile string) (os.FileInfo, error) {
	if c.SshClient == nil {
		if err := c.Connect(); err != nil {
			return nil, err
		}
	}
	stat, err := c.SftpClient.Stat(remoteFile)
	if err != nil {
		return nil, err
	}
	return stat, nil
}

func (c Cli) UploadFile(fileByte []byte, remoteFileName string) (int, error) {
	if c.SshClient == nil {
		if err := c.Connect(); err != nil {
			return -1, err
		}
	}

	ftpFile, err := c.SftpClient.Create(remoteFileName)
	if nil != err {
		return -1, fmt.Errorf("Create remote path failed: %w", err)
	}
	defer ftpFile.Close()

	if nil != err {
		return -1, fmt.Errorf("read local file failed: %w", err)
	}

	ftpFile.Write(fileByte)

	return 0, nil
}
