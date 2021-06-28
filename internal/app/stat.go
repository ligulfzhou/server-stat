package app

import "term-server-stat/pkg/psutil"

type Stat struct {
	Connected     bool
	Error         string
	Alias         string
	CpuCnt        int
	ProcDiskStats *psutil.ProcDiskStats
	Load          *psutil.AvgStat
	NetStat       *psutil.NetStats
	MemStat       *psutil.VirtualMemoryStat
}

func NewEmptyStat(alias string) *Stat {
	return &Stat{
		Connected: false,
		Alias:     alias,
	}
}
