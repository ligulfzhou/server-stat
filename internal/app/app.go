package app

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"sync"
	"term-server-stat/internal/config"
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
}

type State struct {
	Stats sync.Map //map[string]Stat
}

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

	for idx := 0; idx < len(app.Ctls); idx++ {
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
		case stat := <-app.StatCollector:
			app.getStat(stat)
		}
	}
}

func (app *App) getStat(stat Stat) {
	fmt.Println(stat)

}
