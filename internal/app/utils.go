package app

import (
	"fmt"
	"sync"
)

const (
	TB = 1024 * 1024 * 1024 * 1024
	GB = 1024 * 1024 * 1024
	MB = 1024 * 1024
	KB = 1024

	DAY = 24 * 60 * 60
)

// ClearSyncMap clears a sync.Map
func (app *App) ClearSyncMap(syncMap sync.Map) {
	syncMap.Range(func(key interface{}, value interface{}) bool {
		syncMap.Delete(key)
		return true
	})
}

// TruncateString returns a truncated string
func TruncateString(value string, maxLen int) string {
	dots := "..."
	if len(value) > maxLen {
		value = fmt.Sprintf("%s%s", value[0:maxLen-3], dots)
	}
	return value
}

func FormatUsage(value int64) string {
	if value > TB {
		tb := float64(value) / TB
		return fmt.Sprintf("%.2fTB", tb)
	}

	if value > GB {
		gb := float64(value) / GB
		return fmt.Sprintf("%.02fGB", gb)
	}

	if value > MB {
		mb := float64(value) / MB
		return fmt.Sprintf("%.02fMB", mb)
	}
	if value > KB {
		kb := float64(value) / KB
		return fmt.Sprintf("%.02fKB", kb)
	}
	return fmt.Sprintf("%dB", value)
}

func FormatUptime(tm int64) (res string) {
	if tm > DAY {
		d := tm / DAY
		res += fmt.Sprintf("%ddays,", d)
		tm %= DAY
	}
	hours := tm / 60 / 60
	minutes := tm % (60 * 60) / 60
	res += fmt.Sprintf("%d:%d", hours, minutes)
	return
}

func FormatSpeed(value int64) string {
	return FormatUsage(value) + "/s"
}
