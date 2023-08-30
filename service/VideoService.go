package service

import (
	"douyin/models"
	"fmt"
	"strconv"
)

// MaxVideoNum 每次最多返回的视频流数量
const (
	MaxVideoNum = 30
)

// redis 缓存查询

func GetVideoCommentCount(v *models.VideoInfo) error {
	commentCount, err := models.RedisClient.LLen(RedisCtx, INTERACT_COMMENT_KEY+strconv.Itoa(int(v.ID))).Result()
	if err != nil {
		return fmt.Errorf("comment list count error: %v", err)
	}
	v.CommentCount = commentCount
	return nil
}

// ----------------------------------

func GenerateVideoInfo(curId uint, v *models.Video) {

}

func GetVideoInfoById(curId uint, vid uint) {

}

func GetVideoInfoByIds(curId uint, vids []uint) {

}

// =============================================================================================================

func GetVideoById(vid uint) (models.Video, error) {
	video := models.Video{}
	err := models.DB.First(&video, vid).Error
	return video, err
}

func GetVideosByUpdateAt() ([]models.Video, error) {
	videos := make([]models.Video, 10)
	err := models.DB.Order("updated_at desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, err
}

func GetVideosLenght() (int64, error) {
	var video models.Video
	var lenght int64
	err := models.DB.Model(&video).Count(&lenght).Error
	if err != nil {
		return lenght, err
	}
	return lenght, err
}

// func GetVideosByIds(vids []uint) ([]models.Video, error) {
// 	videos := make([]models.Video, len(vids))
// 	err := models.DB.Where("vid in ?", vids).Find(&videos).Error
// 	return videos, err
// }

func GetVideosByUserId(uid uint) ([]models.Video, error) {
	if user, err := GetUserById(uid); err != nil {
		return []models.Video{}, fmt.Errorf("user not found: %v", err)
	} else {
		videos := make([]models.Video, 10)
		err := models.DB.Model(&user).Association("CreatedVideo").Find(&videos)
		return videos, err
	}
}

func CreateVideo(video models.Video) error {
	return models.DB.Create(&video).Error
}
