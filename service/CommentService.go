package service

import (
	"douyin/models"
	"fmt"
	"strconv"
	"time"
)

func GetVideoCommentCount(v *models.VideoInfo) error {
	commentCount, err := models.RedisClient.SCard(RedisCtx, INTERACT_COMMENT_KEY+strconv.Itoa(int(v.ID))).Result()
	if err != nil {
		return fmt.Errorf("comment list count error: %v", err)
	}
	v.CommentCount = commentCount
	return nil
}

func GetCommentIdsByVideoId(vid uint) ([]uint, error) {
	tmp, err := models.RedisClient.SMembers(RedisCtx, INTERACT_COMMENT_KEY+strconv.Itoa(int(vid))).Result()
	cids := make([]uint, len(tmp))
	for i, tmp_cid := range tmp {
		int_cid, _ := strconv.Atoi(tmp_cid)
		cids[i] = uint(int_cid)
	}
	return cids, err
}

func CreateComment(uid uint, vid uint, text string) (*models.Comment, error) {
	comment := models.Comment{
		UserId:  uid,
		VideoId: vid,
		Content: text,
	}
	err := models.DB.Create(&comment).Error
	models.RedisClient.SAdd(RedisCtx, INTERACT_COMMENT_KEY+strconv.Itoa(int(vid)), comment.ID)
	return &comment, err
}

func DeleteComment(id uint) error {
	err := models.DB.Delete(&models.Comment{}, id).Error
	return err
}

func GetCommentsByIds(cids []uint) ([]models.Comment, error) {
	comments := make([]models.Comment, len(cids))
	err := models.DB.Where("id in (?)", cids).Find(&comments).Error
	return comments, err
}

func GenerateCommentInfo(c *models.Comment) models.CommentInfo {
	return models.CommentInfo{
		ID:         int64(c.ID),
		User:       nil,
		Content:    c.Content,
		CreateDate: time.Now().String(),
	}
}

// func GetCommentsByVideoId(vid uint) ([]models.Comment, error) {
// 	comments := []models.Comment{}
// 	err := models.DB.Where("video_id=?", vid).Find(&comments).Error
// 	return comments, err
// }