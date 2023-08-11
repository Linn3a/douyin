package service

import (
	"douyin/models"
	"douyin/public"
	"fmt"
)

func CreateComment(newComment *models.Comment) error {
	err := public.DBConn.Create(newComment).Error
	return err
}

func DeleteComment(id uint) error {
	err := public.DBConn.Delete(&models.Comment{}, id).Error
	return err
}

func GetCommentsByVideoId(vid uint) ([]models.Comment, error) {
	comments := []models.Comment{}
	err := public.DBConn.Where("video_id=?", vid).Find(&comments).Error
	return comments, err
}

func CountCommentsByVideoId(vid uint) (int64, error) {
	comments := []models.Comment{}
	err := public.DBConn.Where("video_id=?", vid).Find(&comments).Error
	counts := int64(len(comments))
	return counts, err
}

func CountCommentsByVideoIds(vids []uint) (map[uint]int64, error) {
	var queryResults []map[string]interface{}
    err := public.DBConn.Table("comments").
        Select("video_id as vid, COUNT(id) as cid_count").
		Where("video_id IN ?", vids).
        Group("video_id").
        Find(&queryResults).Error
	counts := make(map[uint]int64, len(vids))
	for _, result := range(queryResults) {
		counts[result["vid"].(uint)] = int64(result["cid_count"].(int))
	}
	return counts, err
}

func GenerateCommentInfo(c *models.Comment) (models.CommentInfo, error) {
	uid := (*c).UserId
	userInfo, err := GetUserInfoById(uid)
	if err != nil {
		return models.CommentInfo{}, fmt.Errorf("get userinfo failed: %v", err)
	}
	commentInfo := models.NewCommentInfo(c, userInfo)
	return commentInfo, nil
}

func GenerateCommentInfos(comments *[]models.Comment) ([]models.CommentInfo, error) {
	uids := make([]uint, len(*comments))
	for ind, c := range *comments {
		uids[ind] = c.UserId
	}
	userInfoIdMap, err := GetUserInfosByIds(uids)
	if err != nil {
		return []models.CommentInfo{}, fmt.Errorf("get userinfos failed: %v", err)
	}
	commentInfos := make([]models.CommentInfo, len(*comments))
	for ind, c := range *comments {
		commentInfos[ind] = models.NewCommentInfo(
			&c,
			userInfoIdMap[c.UserId],
		)
	}
	return commentInfos, nil
}

func GetCommentInfosByVideoId(vid uint) ([]models.CommentInfo, error) {
	comments, err := GetCommentsByVideoId(vid)
	if err != nil {
		return []models.CommentInfo{}, fmt.Errorf("get comments failed: %v", err)
	}
	commentInfos, err := GenerateCommentInfos(&comments)
	if err != nil {
		return []models.CommentInfo{}, fmt.Errorf("generate commentInfos failed: %v", err)
	}
	return commentInfos, nil
}
