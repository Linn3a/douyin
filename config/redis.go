package config

type Redis struct {
	Address      string `mapstructure:"address" json:"address " yaml:"address"`
	DB           int    `mapstructure:"db" json:"db" yaml:"db"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	ReadTimeout  int    `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout"`
	PoolTimeout  int    `mapstructure:"pool_timeout" json:"pool_timeout" yaml:"pool_timeout"`
}
