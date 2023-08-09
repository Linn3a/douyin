package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app" json:"app" yaml:"app"`
	Database Database `mapstructure:"database" json:"database" yaml:"database"`
}
type App struct {
	Env     string `mapstructure:"env" json:"env" yaml:"env"`
	Port    string `mapstructure:"port" json:"port" yaml:"port"`
	AppName string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	AppUrl  string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
}

var GlobalConfig Config

func InitConfig() error {
	v := viper.New()

	v.SetConfigFile("config/config.json")
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s \n", err))
	}

	if err := v.Unmarshal(&GlobalConfig); err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("%+v", GlobalConfig)
	fmt.Printf(GlobalConfig.App.Env)
	fmt.Printf(GlobalConfig.App.Port)
	fmt.Printf("%v", GlobalConfig.Database)
	return nil
}
