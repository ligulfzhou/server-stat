package commands

import (
	"github.com/spf13/cobra"
	"log"
	"term-server-stat/internal/app"
	conf "term-server-stat/internal/config"
)

var (
	config   string
	mode     string
	StartCmd = &cobra.Command{
		Use:   "show",
		Short: "initialize the database",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/config.yaml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func run() {
	cfg, err := conf.ReadConfig(config)
	if err != nil {
		log.Fatal("parse config file failed : " + err.Error())
	}

	app := app.NewApp(cfg)
	app.Start()
}
