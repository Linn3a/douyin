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

func GenerateCommentInfo(c *models.Comment) (models.CommentInfo, error) {
	uid := (*c).UserId
	if userInfo, err := GetUserInfoById(uid); err != nil {
		return models.CommentInfo{}, fmt.Errorf("get userinfo failed: %v", err)
	} else {
		commentInfo := models.NewCommentInfo(c, userInfo)
		return commentInfo, nil
	}
}

func GenerateCommentInfos(comments *[]models.Comment) ([]models.CommentInfo, error) {
	uids := make([]uint, len(*comments))
	for ind, c := range *comments {
		uids[ind] = c.UserId
	}
	if userInfoIdMap, err := GetUserInfosByIds(uids); err != nil {
		return []models.CommentInfo{}, fmt.Errorf("get userinfos failed: %v", err)
	} else {
		commentInfos := make([]models.CommentInfo, len(*comments))
		for ind, c := range *comments {
			commentInfos[ind] = models.NewCommentInfo(
				&c,
				userInfoIdMap[c.UserId],
			)
		}
		return commentInfos, nil
	}
}

func GetCommentInfosByVideoId(vid uint) ([]models.CommentInfo, error) {
	if comments, err := GetCommentsByVideoId(vid); err != nil {
		return []models.CommentInfo{}, fmt.Errorf("get comments failed: %v", err)
	} else {
		if commentInfos, err := GenerateCommentInfos(&comments); err != nil {
			return []models.CommentInfo{}, fmt.Errorf("generate commentInfos failed: %v", err)
		} else {
			return commentInfos, nil
		}
	}
}
