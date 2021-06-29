package commands

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"server-stat/internal/app"
	conf "server-stat/internal/config"
)

var (
	config string
	mode   string

	rootCmd = &cobra.Command{
		Use:   "server-stat",
		Short: "get server-stat in your ssh/config file",

		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := conf.ReadConfig(config)
			if err != nil {
				log.Fatal("parse config file failed : " + err.Error())
			}

			app := app.NewApp(cfg)
			app.Start()

			return nil
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.AddCommand(StartCmd)
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/config.yaml", "Start server with provided configuration file")
	rootCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}
