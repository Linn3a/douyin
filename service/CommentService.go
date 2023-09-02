package service

import (
	"douyin/models"
	"fmt"
	"time"
	"strings"
	"strconv"
	"douyin/middleware/rabbitmq"
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
	// err := models.DB.Create(&comment).Error
	// 关注消息加入消息队列
	sb := strings.Builder{}
	sb.WriteString(comment.Content)
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(comment.UserId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(comment.VideoId)))
	rabbitmq.RmqCommentAdd.Publish(sb.String())
	fmt.Println("评论消息入队成功")
	models.RedisClient.SAdd(RedisCtx, INTERACT_COMMENT_KEY+strconv.Itoa(int(vid)), comment.ID)
	return &comment, nil
}


func DeleteComment(id uint) error {
	rabbitmq.RmqCommentDel.Publish(strconv.FormatInt(int64(id),10))
	fmt.Println("删除评论消息入队成功")
	return nil
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