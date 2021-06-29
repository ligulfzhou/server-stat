package app

import (
	"server-stat/internal/config"
	"server-stat/pkg/psutil"
	"sync"
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
			cnt += 1
		}

		netStat, err := c.Psutils.GetNetStats()
		if err == nil {
			stat.NetStat = netStat
		} else {
			cnt += 1
		}

		memStat, err := c.Psutils.VirtualMemory()
		if err == nil {
			stat.MemStat = memStat
		} else {
			cnt += 1
		}

		loadStat, err := c.Psutils.ArgLoad()
		if err == nil {
			stat.Load = loadStat
		} else {
			cnt += 1
		}

		hostStat, err := c.Psutils.GetHostInfoStat()
		if err == nil {
			stat.HostStat = hostStat
		} else {
			cnt += 1
		}

		stat.ProcDiskStats = c.Psutils.GetDiskOverallStats()

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
