package service

import (
	"douyin/middleware/rabbitmq"
	"douyin/models"
	"douyin/utils/log"
	"fmt"
	"strconv"
	// "strings"
	"time"
	"encoding/json"
	"gorm.io/gorm"
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

	// err := models.DB.Create(&comment).Error
	// 关注消息加入消息队列
	cIdStr, _ := models.RedisClient.Get(RedisCtx, INTERACT_MAX_COMMENT_KEY).Result()
	models.RedisClient.Incr(RedisCtx, INTERACT_MAX_COMMENT_KEY)
	cIdInt, _ := strconv.Atoi(cIdStr)
	cId := uint(cIdInt)
	cId = cId + 1
	models.RedisClient.SAdd(RedisCtx, INTERACT_COMMENT_KEY+strconv.Itoa(int(vid)), cId)
	// sb := strings.Builder{}
	// sb.WriteString(strconv.Itoa(int(cId)))
	// sb.WriteString(" ")
	// sb.WriteString(text)
	// sb.WriteString(" ")
	// sb.WriteString(strconv.Itoa(int(uid)))
	// sb.WriteString(" ")
	// sb.WriteString(strconv.Itoa(int(vid)))

	comment := models.Comment{
		Model: gorm.Model{
			ID: cId,
		},
		Content: text,
		UserId:  uid,
		VideoId: vid,
	}
	jsonComment,err := json.Marshal(comment)
	if err != nil{
		fmt.Println("comment转换为json错误")
	}
	rabbitmq.RmqCommentAdd.Publish(string(jsonComment))
	log.FieldLog("commentMQ", "info", fmt.Sprintf("successfully add comment: %v", string(jsonComment)))

	return &comment, nil
}

func DeleteComment(id uint) error {
	rabbitmq.RmqCommentDel.Publish(strconv.FormatInt(int64(id), 10))
	log.FieldLog("commentMQ", "info", fmt.Sprintf("successfully delete comment id: %v", id))
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
