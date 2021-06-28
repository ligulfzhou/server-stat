package psutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	CipherList = []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"}
)

type PSUtils struct {
	user, password, host, keyPath, keyString string
	port                                     int
	cipherList                               []string

	platform string

	// for network
	NetworkInterface             string
	RxLastTmstamp, TxLastTmstamp int64
	RxLastTotal, TxLastTotal     int64

	// disk
	StorageDeviceNames  []string
	ProcDiskstatTmstamp int64
	LastDiskStat        ProcDiskStats

	// virtualization
	VirtualizationSystem, VirtualizationRole string

	//HostId
	HostId string

	// kernel version, useful for calc disk stat
	KernelVersion string

	client *ssh.Client
	mux    sync.Mutex
}

func NewPSUtils(user, password, host, keyPath, keyString string, port int) *PSUtils {

	path := keyPath
	if strings.Contains(keyPath, "~") {
		home, _ := os.UserHomeDir()
		path = home + strings.TrimLeft(keyPath, "~")
	}
	return &PSUtils{
		user:       user,
		password:   password,
		host:       host,
		keyPath:    path,
		keyString:  keyString,
		port:       port,
		cipherList: CipherList,
		client:     nil,
	}
}

func (ps *PSUtils) Connect() (bool, error) {
	ps.mux.Lock()
	defer ps.mux.Unlock()

	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		err          error
		errMsg       string

		config = ssh.Config{
			Ciphers: ps.cipherList,
		}
	)

	// get auth method
	auth = make([]ssh.AuthMethod, 0)

	// try password first
	if ps.password != "" {
		auth = append(auth, ssh.Password(ps.password))

		clientConfig = &ssh.ClientConfig{
			User: ps.user,
			Auth: auth,
			// 5 second shall be acceptable.
			Timeout:         5 * time.Second,
			Config:          config,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		addr = fmt.Sprintf("%s:%d", ps.host, ps.port)
		client, err = ssh.Dial("tcp", addr, clientConfig)
		if err == nil {
			ps.client = client
			return true, nil
		} else {
			fmt.Println(err.Error())
			errMsg += err.Error()
		}
	}

	// then check privkey
	if ps.keyPath != "" || ps.keyString != "" {
		var pemBytes []byte
		if ps.keyPath != "" {
			pemBytes, err = ioutil.ReadFile(ps.keyPath)
			if err != nil {
				return false, err
			}
		} else {
			pemBytes = []byte(ps.keyString)
		}

		var signer ssh.Signer
		signer, err = ssh.ParsePrivateKey(pemBytes)
		if err != nil {
			return false, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
		clientConfig = &ssh.ClientConfig{
			User: ps.user,
			Auth: auth,
			// 5 second shall be acceptable.
			Timeout:         5 * time.Second,
			Config:          config,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		addr = fmt.Sprintf("%s:%d", ps.host, ps.port)
		client, err = ssh.Dial("tcp", addr, clientConfig)
		if err == nil {
			ps.client = client
			return true, nil
		} else {
			fmt.Println(err.Error())
			errMsg += err.Error()
		}
	}

	if ps.keyPath != "" || ps.keyString != "" {
		var pemBytes []byte
		if ps.keyPath != "" {
			pemBytes, err = ioutil.ReadFile(ps.keyPath)
			if err != nil {
				return false, err
			}
		} else {
			pemBytes = []byte(ps.keyString)
		}

		var signer ssh.Signer
		signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(ps.password))
		if err != nil {
			return false, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
		clientConfig = &ssh.ClientConfig{
			User: ps.user,
			Auth: auth,
			// 5 second shall be acceptable.
			Timeout:         5 * time.Second,
			Config:          config,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		}

		addr = fmt.Sprintf("%s:%d", ps.host, ps.port)
		client, err = ssh.Dial("tcp", addr, clientConfig)
		if err == nil {
			ps.client = client
			return true, nil
		} else {
			fmt.Println(err.Error())
			errMsg += err.Error()
		}
	}

	fmt.Printf("connections failed with %s", errMsg)
	return false, fmt.Errorf("connections failed with %s", errMsg)
}

func (ps *PSUtils) SudoExec(command string) (string, error) {
	return ps.Exec("sudo " + command)
}

func (ps *PSUtils) Exec(command string) (string, error) {
	ps.mux.Lock()
	defer ps.mux.Unlock()

	session, err := ps.client.NewSession()
	if err != nil {
		return "", err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return "", err
	}

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		return "", err
	}

	// fmt.Println("\n\n\n---------------exec cmd----------------")
	// fmt.Printf("exec cmd: %s, res: %s \n", cmd, b.String())
	// fmt.Println("---------------exec cmd---------------- \n\n\n ")
	return StripString(b.String()), nil
}
