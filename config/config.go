package config

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
