package service

import (
	"douyin/models"
	"douyin/public"
	"fmt"
)

func GetVideoById(vid uint) (models.Video, error) {
	video := models.Video{}
	err := public.DBConn.First(&video, vid).Error
	return video, err
}

// func GetVideosByIds(vids []uint) ([]models.Video, error) {
// 	videos := make([]models.Video, len(vids))
// 	err := public.DBConn.Where("vid in ?", vids).Find(&videos).Error
// 	return videos, err
// }

func GetVideosByUserId(uid uint) ([]models.Video, error) {
	if user, err := GetUserById(uid); err != nil {
		return []models.Video{}, fmt.Errorf("user not found: %v", err)
	} else {
		videos := make([]models.Video, 10)
		err := public.DBConn.Model(&user).Association("CreatedVideo").Find(&videos)
		return videos, err
	}
	
}