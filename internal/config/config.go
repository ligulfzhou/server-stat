package config

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	Log struct {
		File  string `mapstructure:"file"`
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
	SSh struct {
		Frequency           int    `mapstructure:"frequency"`
		DefaultPassword     string `mapstructure:"default_password"`
		DefaultIdentityfile string `mapstructure:"default_identityfile"`
	} `mapstructure:"ssh"`
	Servers struct {
		List []string `mapstructure:"list"`
	} `mapstructure:"servers"`
}

func ReadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(content))))
	if err != nil {
		return nil, err
	}
	var config Config
	err = viper.Unmarshal(&config)
	return &config, err
}
