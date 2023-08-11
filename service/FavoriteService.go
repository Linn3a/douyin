package service

import (
	"douyin/models"
	"douyin/public"
	"fmt"
)

func AddFavoriteVideo(uid uint, vid uint) error {
	user, err := GetUserById(uid)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}
	video, err := GetVideoById(vid)
	if err != nil {
		return fmt.Errorf("video not found: %v", err)
	}
	err = public.DBConn.Model(&user).Association("LikeVideo").Append(&video)
	return err
}

func DeleteFavoriteVideo(uid uint, vid uint) error {
	user, err := GetUserById(uid)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}
	video, err := GetVideoById(vid)
	if err != nil {
		return fmt.Errorf("video not found: %v", err)
	}
	err = public.DBConn.Model(&user).Association("LikeVideo").Delete(video)
	return err
}

func GetFavoriteVideos(uid uint) ([]models.Video, error) {
	user, err := GetUserById(uid)
	if err != nil {
		return []models.Video{}, fmt.Errorf("user not found: %v", err)
	}
	videos := make([]models.Video, 10)
	err = public.DBConn.Model(&user).Association("LikeVideo").Find(&videos)
	return videos, err
}

func CountFavoriteVideos(uid uint) (int64, error) {
	user, err := GetUserById(uid)
	if err != nil {
		return 0, fmt.Errorf("user not found: %v", err)
	}
	count := public.DBConn.Model(&user).Association("LikeVideo").Count()
	return count, nil
}

func CountFavoritedUsers(vid uint) (int64, error) {
	video, err := GetVideoById(vid)
	if err != nil {
		return 0, fmt.Errorf("video not found: %v", err)
	}
	count := public.DBConn.Model(&video).Association("FavoritedUser").Count()
	return count, nil
}

func CountFavoritedUsersByIds(vids []uint) (map[uint]int64, error) {
	var queryResults []map[string]interface{}
    err := public.DBConn.Table("video_likes").
        Select("video_id as vid, COUNT(user_id) as uid_count").
		Where("video_id IN ?", vids).
        Group("video_id").
        Find(&queryResults).Error
	counts := make(map[uint]int64, len(vids))
	for _, result := range(queryResults) {
		counts[result["vid"].(uint)] = int64(result["uid_count"].(int))
	}
	return counts, err
}


func CountUserFavorited(uid uint) (int64, error) {
	videos, err := GetVideosByUserId(uid)
	if err != nil {
		return 0, fmt.Errorf("failed to find all created videos: %v", err)
	}
	count := int64(0)
	for _, video := range videos {
		count += public.DBConn.Model(&video).Association("FavoritedUser").Count()
	}
	return count, nil
}
