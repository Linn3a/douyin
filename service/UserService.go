package service

import (
	// "gorm.io/gorm"
	// "sync/atomic"
	"douyin/models"
	"douyin/public"
)

func GetUserByName(userName string) (models.User, error) {
	tmp := models.User{}
	if err := public.DBConn.Where("name = ?", userName).First(&tmp).Error; err != nil {
		return tmp, err
	}
	return tmp, nil
}

func CreateUser(newUser *models.User) error {
	err := public.DBConn.Create(newUser).Error
	return err
}

func GetUserById(ID uint) (models.User, error) {
	tmp := models.User{}
	if err := public.DBConn.First(&tmp, ID).Error; err != nil {
		return tmp, err
	}
	return tmp, nil
}
