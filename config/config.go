package config

import (
	"douyin/utils/log"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	App      App      `mapstructure:"app" json:"app" yaml:"app"`
	Database Database `mapstructure:"database" json:"database" yaml:"database"`
	Redis    Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	Rabbitmq Rabbitmq `mapstructure:"rabbitmq" json:"rabbitmq" yaml:"rabbitmq"`
}
type App struct {
	Env     string `mapstructure:"env" json:"env" yaml:"env"`
	Port    string `mapstructure:"port" json:"port" yaml:"port"`
	AppName string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	AppUrl  string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
}
type Rabbitmq struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
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
		log.FieldLog("viper", "error", "unmarshal config failed")
		return err
	}
	log.FieldLog("viper", "info", fmt.Sprintf("init config success:%+v", GlobalConfig))
	return nil
}
