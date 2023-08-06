package public

import(
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/Eacient/douyin/models"
)

var (
    DBConn *gorm.DB
)

func InitDatabase(){
	dsn := "root:1515438605@tcp(localhost:3306)/test-gorm?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DBConn,err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	  }
	DBConn.AutoMigrate(&models.User{}, &models.Message{}, &models.Video{}, &models.Comment{})

}