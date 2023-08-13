package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `json:"content"`
	VideoId uint   `json:"video_id"` // has many 关系下游
	UserId  uint   `json:"user_id"`  // has many 关系下游
}

// 用于响应http请求
type CommentInfo struct {
	ID         int64    `json:"id"`
	User       UserInfo `json:"user"`
	Content    string   `json:"content"`
	CreateDate string   `json:"create_date"`
}

func NewCommentInfo(c *Comment, ui UserInfo) CommentInfo {
	return CommentInfo{
		ID:         int64((*c).ID),
		User:       ui,
		Content:    (*c).Content,
		CreateDate: time.Now().String(),
	}
}
