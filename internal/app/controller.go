package app

import (
	"fmt"
	"sync"
	"term-server-stat/internal/config"
	"term-server-stat/pkg/psutil"
)

type Controller struct {
	Config          *config.Config
	Psutils         *psutil.PSUtils
	SshConfig       *SshConfig
	connected       bool
	connectionError string
	Reload          chan struct{}
	Done            chan struct{}
	connMux         sync.Mutex
	fetchMux        sync.Mutex
}

func NewController(cfg *config.Config, sshConfig *SshConfig) *Controller {
	ctl := &Controller{
		Config:    cfg,
		SshConfig: sshConfig,
		Reload:    make(chan struct{}),
		Done:      make(chan struct{}),
	}

	return ctl
}

func (c *Controller) Start(uploader chan<- Stat) {
	for {
		select {
		case <-c.Reload:
			go c.Fetch(uploader)
		case <-c.Done:
			return
		}
	}
}

func (c *Controller) Fetch(uploader chan<- Stat) {
	c.fetchMux.Lock()
	defer c.fetchMux.Unlock()

	if c.Psutils == nil || !c.connected {
		c.ConnectServer()
	}

	stat := Stat{
		Alias:     c.SshConfig.Alias,
		Connected: c.connected,
	}
	if c.connected {
		cnt := 0
		cpuCnt, err := c.Psutils.CPUCount(true)
		if err == nil {
			stat.CpuCnt = cpuCnt
		} else {
			fmt.Printf("get cpucnt %v", err)
			cnt += 1
		}

		netStat, err := c.Psutils.GetNetStats()
		if err == nil {
			stat.NetStat = netStat
		} else {
			fmt.Printf("get netcnt %v", err)
			cnt += 1
		}

		memStat, err := c.Psutils.VirtualMemory()
		if err == nil {
			stat.MemStat = memStat
		} else {
			fmt.Printf("get memcnt %v", err)
			cnt += 1
		}

		loadStat, err := c.Psutils.ArgLoad()
		if err == nil {
			stat.Load = loadStat
		} else {
			fmt.Printf("get loadcnt %v", err)
			cnt += 1
		}
		if cnt >= 3 {
			stat.Connected = false
			go c.ConnectServer()
		}
	} else if c.connectionError != "" {
		stat.Error = c.connectionError
	}
	uploader <- stat
}

func (c *Controller) ConnectServer() {
	c.connMux.Lock()
	defer c.connMux.Unlock()

	c.Psutils = psutil.NewPSUtils(c.SshConfig.User, c.SshConfig.Password, c.SshConfig.HostName, c.SshConfig.IdentityFile, "", c.SshConfig.Port)
	connected, err := c.Psutils.Connect()
	c.connected = connected
	if err != nil && !connected {
		c.connectionError = err.Error()
	}
}
