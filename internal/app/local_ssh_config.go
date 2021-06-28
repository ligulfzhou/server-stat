package app

import (
	"fmt"
	"github.com/kevinburke/ssh_config"
	"strconv"
	conf "term-server-stat/internal/config"
	"term-server-stat/pkg/psutil"
)

// DefaultIdentityFile github.com/kevinburke/ssh_config use following default value
const DefaultIdentityFile = "~/.ssh/identity"

type SshConfig struct {
	Alias        string
	HostName     string
	Port         int
	User         string
	IdentityFile string
	Password     string

	Psutil *psutil.PSUtils
}

type SshConfigFetcher struct {
	Config *conf.Config
}

func NewFetcher(config *conf.Config) *SshConfigFetcher {
	return &SshConfigFetcher{
		Config: config,
	}
}

func (f *SshConfigFetcher) Fetch() ([]SshConfig, error) {
	confs := make([]SshConfig, len(f.Config.Servers.List))
	for idx, alias := range f.Config.Servers.List {
		sshConfig := SshConfig{
			Alias:    alias,
			HostName: ssh_config.Get(alias, "HostName"),
			User:     ssh_config.Get(alias, "User"),
			Password: f.Config.SSh.DefaultPassword,
		}

		portStr := ssh_config.Get(alias, "Port")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("port %s not valid: %s", portStr, err.Error())
		} else {
			sshConfig.Port = port
		}

		identityFile := ssh_config.Get(alias, "IdentityFile")

		if identityFile != DefaultIdentityFile {
			sshConfig.IdentityFile = identityFile
		} else {
			sshConfig.IdentityFile = f.Config.SSh.DefaultIdentityfile
		}
		//if identityFile != DefaultIdentityFile {
		//	sshConfig.IdentityFile = identityFile
		//} else if f.Config.SSh.DefaultIdentityfile != "" {
		//	sshConfig.IdentityFile = f.Config.SSh.DefaultIdentityfile
		//}

		confs[idx] = sshConfig
	}
	return confs, nil
}
