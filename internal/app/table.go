package app

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"server-stat/pkg/humanize"
	"server-stat/pkg/pad"
	"strings"
)

const dots = "..."

type RowCell struct {
	LeftMargin  int
	RightMargin int
	Width       int
	LeftAlign   bool
	Color       func(a ...interface{}) string
	Text        string
}

// SupportedCoinTableHeaders are all the supported coin table header columns
var tableHeaders = []string{
	"Alias",
	"Uptime",
	"Cpu",
	"Load(1/5/15)",
	"Mem(free/all)",
	"Net(total[r/w])",
	"Net Speed",
	"Disk(total sectors[r/w])",
	"Disk IOPS[r,w]",
	//"Disk Speed",
}

// GetStatTable returns the table for diplaying the coins
func (app *App) GetStatTable() [][]*RowCell {
	rows := [][]*RowCell{}
	app.ClearSyncMap(app.State.tableColumnWidths)
	app.ClearSyncMap(app.State.tableColumnAlignLeft)

	stats := app.GetStatList()
	for _, stat := range stats {
		rowCells := []*RowCell{}
		for _, header := range tableHeaders {
			leftMargin := 1
			rightMargin := 1
			switch header {
			case "Alias":
				txt := TruncateString(stat.Alias, 12)
				app.SetTableColumnWidthFromString(header, txt)
				app.SetTableColumnAlignLeft(header, true)
				rowCells = append(rowCells, &RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       nil,
					Text:        txt,
				})
			case "Uptime":
				txt := dots
				if stat.Connected && stat.HostStat != nil {
					txt = FormatUptime(stat.HostStat.Uptime)
				}
				app.SetTableColumnWidthFromString(header, txt)
				app.SetTableColumnAlignLeft(header, true)
				rowCells = append(rowCells, &RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       nil,
					Text:        txt,
				})
			case "Cpu":
				txt := dots
				if stat.Connected && stat.CpuCnt != 0 {
					txt = fmt.Sprintf("%d", stat.CpuCnt)
				}
				app.SetTableColumnWidthFromString(header, txt)
				app.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       nil,
					Text:        txt,
				})
			case "Load(1/5/15)":
				txt := dots
				if stat.Connected && stat.Load != nil {
					txt = fmt.Sprintf("%s/%s/%s",
						humanize.Monetaryf(stat.Load.Load1, 2),
						humanize.Monetaryf(stat.Load.Load15, 2),
						humanize.Monetaryf(stat.Load.Load15, 2))
				}
				app.SetTableColumnWidthFromString(header, txt)
				app.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       nil,
					Text:        txt,
				})
			case "Mem(free/all)":
				txt := dots
				if stat.Connected && stat.MemStat != nil {
					txt = fmt.Sprintf("%s/%s",
						FormatUsage(stat.MemStat.Free),
						FormatUsage(stat.MemStat.Total))
				}
				app.SetTableColumnWidthFromString(header, txt)
				app.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       nil,
					Text:        txt,
				})
			case "Net(total[r/w])":
				txt := dots
				if stat.Connected && stat.NetStat != nil {
					txt = fmt.Sprintf("%s,%s",
						FormatUsage(stat.NetStat.RxTotal),
						FormatUsage(stat.NetStat.TxTotal))
				}
				app.SetTableColumnWidthFromString(header, txt)
				app.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       nil,
					Text:        txt,
				})
			case "Net Speed":
				txt := dots
				if stat.Connected && stat.NetStat != nil {
					txt = fmt.Sprintf("%s,%s",
						FormatSpeed(stat.NetStat.RxSpeed),
						FormatSpeed(stat.NetStat.TxSpeed))
				}
				app.SetTableColumnWidthFromString(header, txt)
				app.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       nil,
					Text:        txt,
				})
			case "Disk(total sectors[r/w])":
				txt := dots
				if stat.Connected && stat.ProcDiskStats != nil {
					txt = fmt.Sprintf("%d,%d",
						stat.ProcDiskStats.SectorsRead,
						stat.ProcDiskStats.SectorsWritten)
				}
				app.SetTableColumnWidthFromString(header, txt)
				app.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       nil,
					Text:        txt,
				})
			case "Disk IOPS[r,w]":
				txt := dots
				if stat.Connected && stat.ProcDiskStats != nil {
					txt = fmt.Sprintf("%d,%d",
						stat.ProcDiskStats.ReadIOPS,
						stat.ProcDiskStats.WriteIOPS)
				}
				app.SetTableColumnWidthFromString(header, txt)
				app.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       nil,
					Text:        txt,
				})
			case "Disk Speed":
				txt := dots
				if stat.Connected && stat.ProcDiskStats != nil {
					txt = fmt.Sprintf("%s,%s",
						FormatSpeed(stat.ProcDiskStats.ReadSpeed),
						FormatSpeed(stat.ProcDiskStats.WriteSpeed))
				}
				app.SetTableColumnWidthFromString(header, txt)
				app.SetTableColumnAlignLeft(header, false)
				rowCells = append(rowCells, &RowCell{
					LeftMargin:  leftMargin,
					RightMargin: rightMargin,
					LeftAlign:   false,
					Color:       nil,
					Text:        txt,
				})
			}
		}
		rows = append(rows, rowCells)
	}
	for _, row := range rows {
		for idx, header := range tableHeaders {
			row[idx].Width = app.GetTableColumnWidth(header)
		}
	}
	return rows
}

func (app *App) GetStatList() []*Stat {
	//stats := make([]*Stat, len(app.Cfg.Servers.List))
	stats := []*Stat{}
	for _, alias := range app.Cfg.Servers.List {
		var stat *Stat
		istat, ok := app.State.Stats.Load(alias)
		if ok {
			if tstat, ok := istat.(Stat); ok {
				stat = &tstat
			} else {
				stat = NewEmptyStat(alias)
			}
		} else {
			stat = NewEmptyStat(alias)
		}
		//fmt.Printf("stat %v \n", stat)
		stats = append(stats, stat)
	}
	return stats
}

// RefreshTable read stats and put to table
func (app *App) RefreshTable() error {

	app.Gui.Update(func(gui *gocui.Gui) error {
		v, err := app.Gui.View(TableViewName)
		if err != nil {
			return err
		}
		v.Clear()

		rows := app.GetStatTable()
		//fmt.Println(rows)
		for _, row := range rows {
			//fmt.Println(row)
			line := []string{}
			for i, col := range tableHeaders {
				width := app.GetTableColumnWidth(col)

				leftAlign := app.GetTableColumnAlignLeft(col)
				padfn := pad.Left
				padLeft := 1
				if i == 0 {
					padLeft = 0
				}
				if leftAlign {
					padfn = pad.Right
				}

				colStr := fmt.Sprintf(
					"%s%s%s",
					strings.Repeat(" ", padLeft),
					padfn(row[i].Text, width+(1-padLeft), " "),
					strings.Repeat(" ", 1),
				)
				line = append(line, colStr)
			}
			fmt.Fprintln(v, strings.Join(line, ""))
		}
		return nil
	})
	return nil
}
