package bootstrap

import (
	"douyin/config"
	"fmt"
	"github.com/spf13/viper"
)

func Init() {
	v := viper.New()
	v.SetConfigFile("config/config.json")
	v.SetConfigType("json")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config failed: %s \n", err))
	}
	//println(v)
	//fmt.Printf("%+v", v)
	testConfig := config.Config{}
	if err := v.Unmarshal(&testConfig); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v", testConfig)
}
