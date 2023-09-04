package models

import (
	"context"
	"douyin/config"
	"douyin/utils/log"
	"fmt"
	// "time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		config.GlobalConfig.Database.Mysql.UserName,
		config.GlobalConfig.Database.Mysql.Password,
		config.GlobalConfig.Database.Mysql.Host,
		config.GlobalConfig.Database.Mysql.Port,
		config.GlobalConfig.Database.Mysql.Database,
		config.GlobalConfig.Database.Mysql.Charset,
	)

	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		//panic("failed to connect database")
		return err
	}

	err = DB.AutoMigrate(&User{}, &Message{}, &Video{}, &Comment{})
	if err != nil {
		return err
	}
	return nil
}

var RedisClient *redis.Client

func InitRedis() error {
	log.FieldLog("redis", "info", "connect Redis")
	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.GlobalConfig.Redis.Address,
		//Password: config.GlobalConfig.Redis.Password,
		DB: config.GlobalConfig.Redis.DB,
		//超时
		// ReadTimeout:  time.Duration(config.GlobalConfig.Redis.ReadTimeout) * time.Second,
		// WriteTimeout: time.Duration(config.GlobalConfig.Redis.WriteTimeout) * time.Second,
		// PoolTimeout:  time.Duration(config.GlobalConfig.Redis.PoolTimeout) * time.Second,
	})
	if _, err := RedisClient.Ping(context.Background()).Result(); err != nil {
		log.FieldLog("redis", "error", fmt.Sprintf("connect Redis failed: %v", err))
		return err
	}
	log.FieldLog("redis", "info", "connect Redis success")
	return nil
}
