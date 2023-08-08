package public

import (
	"douyin/config"
	"douyin/utils/jwt"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Config = config.Config{}
	DBConn *gorm.DB
	Jwt    *jwt.JWT
)

func InitDatabase() {
	mysqlConfig := Config.Database.Mysql
	var err error
	DBConn, err = gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
			mysqlConfig.UserName,
			mysqlConfig.Password,
			mysqlConfig.Host,
			mysqlConfig.Database,
			mysqlConfig.Charset,
		),
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// DBConn.AutoMigrate(&models.User{}, &models.Message{}, &models.Video{}, &models.Comment{})
}

func InitJWT() {
	// TODO: ADD jwt signing key to configuration file
	Jwt = jwt.NewJWT([]byte("test_key"))
}
