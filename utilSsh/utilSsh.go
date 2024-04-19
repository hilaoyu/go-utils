package utilSsh

import (
	"fmt"
	"github.com/hilaoyu/go-utils/utilProxy"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"strings"
	"time"
)

type SshExecFunc func(sshClient *ssh.Client, sftpClient *sftp.Client) error

const PROXY_TYPE_SOCKS5 = "socks5"

type SshClient struct {
	addr       string
	user       string
	password   string
	sshClient  *ssh.Client
	sftpClient *sftp.Client

	useProxy            string
	proxySocks5         utilProxy.UtilProxy
	proxySocks5Addr     string
	proxySocks5user     string
	proxySocks5Password string
}

func NewSshClient() (client *SshClient) {
	client = &SshClient{}
	return
}

func (c *SshClient) UseProxySocks5(proxyAddr string, proxyUser string, proxyPassword string) (err error) {
	proxyAddr = strings.TrimSpace(proxyAddr)
	c.useProxy = PROXY_TYPE_SOCKS5
	if "" == proxyAddr {
		err = fmt.Errorf("proxy addr can't be empty")
		return
	}
	c.proxySocks5Addr = proxyAddr
	c.proxySocks5user = proxyUser
	c.proxySocks5Password = proxyPassword

	d, err := utilProxy.NewProxySocks5(c.proxySocks5Addr, c.proxySocks5user, c.proxySocks5Password)
	if nil != err {
		err = fmt.Errorf("proxy connect err: %+v", err)
		return
	}
	c.proxySocks5 = d

	return
}

func (c *SshClient) Connect(addr string, user string, password string, timeout ...time.Duration) (err error) {
	c.addr = addr
	c.user = user
	c.password = password
	c.addr = strings.TrimSpace(c.addr)
	if "" == c.addr {
		err = fmt.Errorf("ssh addr can't be empty")
		return
	}
	if "" == c.user {
		err = fmt.Errorf("ssh user can't be empty")
		return
	}
	if "" == c.password {
		err = fmt.Errorf("ssh password can't be empty")
		return
	}
	if !strings.Contains(c.addr, ":") {
		c.addr += ":22"
	}

	connTimeout := time.Duration(10) * time.Second
	if len(timeout) > 0 && timeout[0] > 0 {
		connTimeout = timeout[0]
	}

	config := &ssh.ClientConfig{}
	config.SetDefaults()
	config.User = c.user
	config.Auth = []ssh.AuthMethod{ssh.Password(c.password)}
	config.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }
	config.Timeout = connTimeout

	switch c.useProxy {
	case PROXY_TYPE_SOCKS5:
		if nil == c.proxySocks5 {
			err = fmt.Errorf("proxy not connected")
			return
		}

		conn, err1 := c.proxySocks5.NewConn("tcp", c.addr, config.Timeout)
		if nil != err1 {
			err = fmt.Errorf("dial tcp error: %+v", err1)
			return
		}

		sshConn, chans, reqs, err1 := ssh.NewClientConn(conn, c.addr, config)
		if nil != err1 {
			err = fmt.Errorf("ssh conn error: %+v", err1)
			return
		}

		c.sshClient = ssh.NewClient(sshConn, chans, reqs)
		break
	default:
		c.sshClient, err = ssh.Dial("tcp", c.addr, config)
		if nil != err {
			return err
		}
		break
	}

	c.sftpClient, err = sftp.NewClient(c.sshClient)

	return
}

func (c *SshClient) Exec(command string, wait ...bool) (output string, err error) {
	if nil == c.sshClient {
		err = fmt.Errorf("ssh not connected")
		return
	}
	session, err := c.sshClient.NewSession()
	if err != nil {
		err = fmt.Errorf("ssh NewSession error: %+v", err)
		return
	}
	defer session.Close()

	var buf []byte
	var waitOutput bool = true
	if len(wait) > 0 {
		waitOutput = wait[0]
	}
	if waitOutput {
		buf, err = session.CombinedOutput(command)
	} else {
		err = session.Start(command)
	}

	output = string(buf)
	//fmt.Println(output)
	return output, err

}

func (c *SshClient) SendFile(localFile string, remoteFile string) (err error) {
	if nil == c.sftpClient {
		err = fmt.Errorf("sftp not connected")
		return
	}
	srcFile, err := os.Open(localFile)
	if err != nil {
		err = fmt.Errorf("open local file failed,file: %s ,error: %+v", localFile, err)
		return
	}
	//defer srcFile.Close()

	dstFile, err := c.sftpClient.Create(remoteFile)
	if err != nil {
		err = fmt.Errorf("create remote file failed,file: %s ,error: %+v", remoteFile, err)
		return
	}
	defer dstFile.Close()

	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		buf = buf[0:n]
		dstFile.Write(buf)
	}
	return nil
}

func (c *SshClient) Close() {
	if nil != c.sftpClient {
		c.sftpClient.Close()
	}
	if nil != c.sshClient {
		c.sshClient.Close()
	}

}
