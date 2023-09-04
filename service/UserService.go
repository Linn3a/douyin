package service

import (
	// "gorm.io/gorm"
	// "sync/atomic"
	"douyin/models"
	"douyin/utils/jwt"
	"douyin/utils/log"
)

func GetUserByName(userName string) (models.User, error) {
	tmp := models.User{}
	if err := models.DB.Where("name = ?", userName).First(&tmp).Error; err != nil {
		log.FieldLog("gorm", "error", "get user by name error")
		return tmp, err
	}
	return tmp, nil
}

func CreateUser(newUser *models.User) error {
	err := models.DB.Create(newUser).Error
	return err
}

func GetUserById(ID uint) (models.User, error) {
	tmp := models.User{}
	if err := models.DB.First(&tmp, ID).Error; err != nil {
		return tmp, err
	}
	return tmp, nil
}

func GetUsersByIds(IDs []uint) ([]models.User, error) {
	tmp := []models.User{}
	if err := models.DB.Where("id in (?)", IDs).Find(&tmp).Error; err != nil {
		return tmp, err
	}
	return tmp, nil
}

func GenerateToken(u *models.User) (string, error) {
	token, err := jwt.JwtClient.NewToken(u.ID)
	if err != nil {
		log.FieldLog("jwt", "error", "create token error")
		return "token", err
	}
	return token, nil
}

func GenerateUserInfo(u *models.User) models.UserInfo {
	return models.UserInfo{
		ID:              int64(u.ID),
		Name:            u.Name,
		FollowCount:     0,
		FollowerCount:   0,
		IsFollow:        false,
		Avatar:          u.Avatar,
		BackgroundImage: u.BackgroundImage,
		Signature:       u.Signature,
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:   0,
	}
}

func GetUserInfoById(id uint) (models.UserInfo, error) {
	user, err := GetUserById(id)
	if err != nil {
		return models.UserInfo{}, err
	}
	userInfo := GenerateUserInfo(&user)
	err = GetUserFollowCount(&userInfo)
	if err != nil {
		return userInfo, err
	}
	err = GetUserFollowerCount(&userInfo)
	if err != nil {
		return userInfo, err
	}
	err = GetUserTotalFavorited(&userInfo)
	if err != nil {
		return userInfo, err
	}
	err = GetUserWorkCount(&userInfo)
	if err != nil {
		return userInfo, err
	}
	err = GetUserFavoriteCount(&userInfo)
	if err != nil {
		return userInfo, err
	}
	return userInfo, nil
}

func GetUserInfoMapByIds(ids []uint) (map[uint]models.UserInfo, error) {
	tmp := make(map[uint]models.UserInfo, len(ids))
	users, err := GetUsersByIds(ids)
	if err != nil {
		return tmp, err
	}
	for _, user := range users {
		tmpInfo := GenerateUserInfo(&user)
		err = GetUserFollowCount(&tmpInfo)
		if err != nil {
			return tmp, err
		}
		err = GetUserFollowerCount(&tmpInfo)
		if err != nil {
			return tmp, err
		}
		err = GetUserTotalFavorited(&tmpInfo)
		if err != nil {
			return tmp, err
		}
		err = GetUserWorkCount(&tmpInfo)
		if err != nil {
			return tmp, err
		}
		err = GetUserFavoriteCount(&tmpInfo)
		if err != nil {
			return tmp, err
		}
		tmp[user.ID] = tmpInfo
	}
	return tmp, nil
}

func GetUserInfosByIds(ids []uint) ([]models.UserInfo, error) {
	tmp := make([]models.UserInfo, len(ids))
	users, err := GetUsersByIds(ids)
	if err != nil {
		return tmp, err
	}
	for i, user := range users {
		tmpInfo := GenerateUserInfo(&user)
		err = GetUserFollowCount(&tmpInfo)
		if err != nil {
			return tmp, err
		}
		err = GetUserFollowerCount(&tmpInfo)
		if err != nil {
			return tmp, err
		}
		err = GetUserTotalFavorited(&tmpInfo)
		if err != nil {
			return tmp, err
		}
		err = GetUserWorkCount(&tmpInfo)
		if err != nil {
			return tmp, err
		}
		err = GetUserFavoriteCount(&tmpInfo)
		if err != nil {
			return tmp, err
		}
		tmp[i] = tmpInfo
	}
	return tmp, nil
}
