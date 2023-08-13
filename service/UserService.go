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
	if err := models.DB.Where("name = ?", userName).First(&tmp).Error; err != nil {
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

func GenerateToken(u *models.User) (string, error) {
	// TODO: Add expired time to token claims
	token, err := public.Jwt.CreateToken(jwt.CustomClaims{ID: int64((*u).ID)}); 
	if err != nil {
		return token, err
	}
	return token, nil
}

func ParseToken(token string) (*jwt.CustomClaims, error) {
	parsedToken, err := public.Jwt.ParseToken(token)
	if err != nil {
		return nil, err
	}
	return parsedToken, nil

}

func GenerateUserInfo(u *models.User) models.UserInfo {
	// TODO: add info for other field
	//	e.g. total_favorite, work_count,etc.
	userInfo := models.NewUserInfo(u)
	userInfo.FollowCount, _ = GetFollowCnt(u.ID)
	userInfo.FollowerCount, _ = GetFollowerCnt(u.ID)
	userInfo.FavoriteCount, _ = CountUserFavorited(u.ID)
	userInfo.TotalFavorited, _ = CountUserFavorited(u.ID)
	return userInfo
}

func GetUserInfoById(id uint) (models.UserInfo, error) {
	user, err := GetUserById(id)
	if err != nil {
		return models.UserInfo{}, err
	}
	userInfo := GenerateUserInfo(&user)
		return userInfo, nil
}

func GetUsersByIds(uids []uint) ([]models.User, error) {
	users := make([]models.User, len(uids))
	err := models.DB.Where("id in ?", uids).Find(&users).Error
	return users, err
}

func GetUserInfosByIds(uids []uint) (map[uint]models.UserInfo, error) {
	users, err := GetUsersByIds(uids)
	if err != nil {
		var userInfoIdMap map[uint]models.UserInfo
		return userInfoIdMap, err
	}
	// fmt.Printf("users: %v\n", users)
	userInfoIdMap := make(map[uint]models.UserInfo, len(users))
	for _, u := range users {
		userInfoIdMap[u.ID] = GenerateUserInfo(&u)
	}
	return userInfoIdMap, nil
}
