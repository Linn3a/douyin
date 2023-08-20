package config

type Database struct {
	Mysql Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	OSS   OSS   `mapstructure:"oss" json:"oss" yaml:"oss"`
}

type Mysql struct {
	DriverName          string `mapstructure:"driver_name" json:"driver_name " yaml:"driver_name"`
	Host                string `mapstructure:"host" json:"host" yaml:"host"`
	Port                string `mapstructure:"port" json:"port" yaml:"port"`
	Database            string `mapstructure:"database" json:"database" yaml:"database"`
	UserName            string `mapstructure:"username" json:"username" yaml:"username"`
	Password            string `mapstructure:"password" json:"password" yaml:"password"`
	Charset             string `mapstructure:"charset" json:"charset" yaml:"charset"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns        int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	LogMode             string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer" json:"enable_file_log_writer" yaml:"enable_file_log_writer"`
	LogFilename         string `mapstructure:"log_filename" json:"log_filename" yaml:"log_filename"`
}

// OSS object database
type OSS struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id" json:"access_key_id" yaml:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret" yaml:"access_key_secret"`
	BucketName      string `mapstructure:"bucket_name" json:"bucket_name" yaml:"bucket_name"`
}
