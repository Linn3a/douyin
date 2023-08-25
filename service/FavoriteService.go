package service

import (
	"douyin/models"
	"strings"
	"strconv"
	"douyin/middleware/rabbitmq"
	"fmt"
)

// audience2video
func AddFavoriteVideo(uid uint, vid uint) error {
	// like消息加入消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(uid)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(vid)))
	rabbitmq.RmqLikeAdd.Publish(sb.String())
	fmt.Println("like消息入队成功")
	return nil
	// user, err := GetUserById(uid)
	// if err != nil {
	// 	return fmt.Errorf("user not found: %v", err)
	// }
	// video, err := GetVideoById(vid)
	// if err != nil {
	// 	return fmt.Errorf("video not found: %v", err)
	// }
	// user := models.User{}
	// user.ID = uid
	// video := models.Video{}
	// video.ID = vid
	// err := models.DB.Model(&user).Association("LikeVideo").Append(&video)
	// return err
}

// audience2video
func DeleteFavoriteVideo(uid uint, vid uint) error {
	// like取消消息加入消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(uid)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(vid)))
	rabbitmq.RmqLikeDel.Publish(sb.String())
	fmt.Println("like取消消息入队成功")
	return nil
	// user, err := GetUserById(uid)
	// if err != nil {
	// 	return fmt.Errorf("user not found: %v", err)
	// }
	// video, err := GetVideoById(vid)
	// if err != nil {
	// 	return fmt.Errorf("video not found: %v", err)
	// }
	// user := models.User{}
	// user.ID = uid
	// video := models.Video{}
	// video.ID = vid
	// err := models.DB.Model(&user).Association("LikeVideo").Delete(video)
	// return err
}

// audience2video
func GetFavoriteVideos(uid uint) ([]models.Video, error) {
	// user, err := GetUserById(uid)
	// if err != nil {
	// 	return []models.Video{}, fmt.Errorf("user not found: %v", err)
	// }
	user := models.User{}
	user.ID = uid
	videos := make([]models.Video, 10)
	err := models.DB.Model(&user).Association("LikeVideo").Find(&videos)
	return videos, err
}

// audience2video
func CountFavoriteVideos(uid uint) (int64, error) {
	// user, err := GetUserById(uid)
	// if err != nil {
	// 	return 0, fmt.Errorf("user not found: %v", err)
	// }
	user := models.User{}
	user.ID = uid
	count := models.DB.Model(&user).Association("LikeVideo").Count()
	return count, nil
}

// video2audience
func CountFavoritedUsers(vid uint) (int64, error) {
	// video, err := GetVideoById(vid)
	// if err != nil {
	// 	return 0, fmt.Errorf("video not found: %v", err)
	// }
	video := models.Video{}
	video.ID = vid
	count := models.DB.Model(&video).Association("FavoritedUser").Count()
	return count, nil
}

// multiple video2audience
func CountFavoritedUsersByIds(vids []uint) (map[uint]int64, error) {
	var queryResults []map[string]interface{}
	err := models.DB.Table("video_likes").
		Select("video_id as vid, COUNT(user_id) as uid_count").
		Where("video_id IN ?", vids).
		Group("video_id").
		Find(&queryResults).Error
	counts := make(map[uint]int64, len(vids))
	for _, result := range queryResults {
		counts[uint(result["vid"].(uint64))] = result["uid_count"].(int64)
	}
	return counts, err
}

// author2audience
func CountUserFavorited(uid uint) (int64, error) {
	// join video and video_likes
	// videos, err := GetVideosByUserId(uid)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to find all created videos: %v", err)
	// }
	// count := int64(0)
	// for _, video := range videos {
	// 	count += models.DB.Model(&video).Association("FavoritedUser").Count()
	// }
	// return count, nil
	var count int64
	models.DB.Table("videos").
	Joins("JOIN video_likes ON videos.id = video_likes.video_id").
	Where("videos.author_id = ?", uid).
	Count(&count)
	return count, nil
}
