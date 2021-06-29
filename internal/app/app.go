package app

import (
	"github.com/jroimartin/gocui"
	"log"
	"server-stat/internal/config"
	"sync"
	"time"
)

type App struct {
	Cfg           *config.Config
	refreshMux    sync.Mutex
	refreshTicker *time.Ticker
	Lsc           *[]SshConfig
	Ctls          []*Controller
	Gui           *gocui.Gui
	StatCollector chan Stat

	State *State
}

type State struct {
	Stats sync.Map //map[string]Stat

	tableColumnWidths    sync.Map
	tableColumnAlignLeft sync.Map
}

const (
	TableViewName       = "table"
	TableHeaderViewName = "table_header"
)

func NewApp(cfg *config.Config) *App {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}

	// load sshConfigs from local
	fetcher := NewFetcher(cfg)
	sshConfigs, err := fetcher.Fetch()
	if err != nil {
		log.Fatal("fetch local sshConfig failed with " + err.Error())
	}

	if len(sshConfigs) <= 0 {
		log.Fatal("no local sshConfigs found...")
	}

	app := &App{
		Gui:           gui,
		Cfg:           cfg,
		Lsc:           &sshConfigs,
		StatCollector: make(chan Stat),

		State: &State{
			Stats:                sync.Map{},
			tableColumnWidths:    sync.Map{},
			tableColumnAlignLeft: sync.Map{},
		},
	}
	app.Gui.SetManagerFunc(app.Layout)

	ctrls := make([]*Controller, len(sshConfigs))
	for idx := 0; idx < len(sshConfigs); idx++ {
		ctrls[idx] = NewController(cfg, &sshConfigs[idx])
	}
	app.Ctls = ctrls

	if err := Keybindings(app.Gui); err != nil {
		log.Fatalf("binding key failed with %v", err)
	}

	if cfg.SSh.Frequency == 0 {
		app.refreshTicker = time.NewTicker(time.Duration(1) * time.Second)
	} else {
		app.refreshTicker = time.NewTicker(time.Duration(cfg.SSh.Frequency) * time.Second)
	}

	return app
}

func (app *App) Start() {
	defer app.Gui.Close()

	for idx, _ := range app.Ctls {
		go app.Ctls[idx].Start(app.StatCollector)
	}
	go app.Refresh()

	if err := app.Gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("main loop: %v", err)
	}
}

func (app *App) Refresh() {
	for {
		select {
		case <-app.refreshTicker.C:
			for _, ctrl := range app.Ctls {
				ctrl.Reload <- struct{}{}
			}
			go func() {
				app.RefreshAll()
			}()
		case stat := <-app.StatCollector:
			app.getStat(stat)
		}
	}
}

//store the received stat from controller
func (app *App) getStat(stat Stat) {
	app.State.Stats.Store(stat.Alias, stat)
}

func (app *App) RefreshAll() {
	app.RefreshTable()
	app.UpdateTableHeader()
}
