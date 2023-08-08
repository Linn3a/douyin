package service

import (
	// "gorm.io/gorm"
	// "sync/atomic"
	"douyin/models"
	"douyin/public"
	"douyin/utils/jwt"
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

func GenerateToken(u *models.User) (string, error) {
	// TODO: Add expired time to token claims
	if token, err := public.Jwt.CreateToken(jwt.CustomClaims{
		Id: int64((*u).ID),
	}); err == nil {
		return token, nil
	} else {
		return token, err
	}
}

func ParseToken(token string) (*jwt.CustomClaims, error) {
	if token, err := public.Jwt.ParseToken(token); err == nil {
		return token, nil
	} else {
		return nil, err
	}
}
