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
		ID: int64((*u).ID),
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

func GenerateUserInfo(u *models.User) models.UserInfo {
	return models.NewUserInfo(u)
}

func GetUserInfoById(id uint) (models.UserInfo, error) {
	if user, err := GetUserById(id); err == nil {
		userInfo := GenerateUserInfo(&user)
		return userInfo, nil
	} else {
		return models.UserInfo{}, err
	}
}

func GetUsersByIds(uids []uint) ([]models.User, error) {
	users := make([]models.User, len(uids))
	err := public.DBConn.Where("id in ?", uids).Find(&users).Error
	return users, err
}

func GetUserInfosByIds(uids []uint) (map[uint]models.UserInfo, error) {
	if users, err := GetUsersByIds(uids); err == nil {
		// userInfos := make([]models.UserInfo, len(users))
		userInfoIdMap := make(map[uint]models.UserInfo, len(users))
		for _, u := range(users) {
			userInfoIdMap[u.ID] = GenerateUserInfo(&u)
		}
		return userInfoIdMap, nil
	} else {
		var userInfoIdMap map[uint]models.UserInfo
		return userInfoIdMap, err
	}
}
