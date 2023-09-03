package models

import (
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
	ID         int64     `json:"id"`
	User       *UserInfo `json:"user"`
	Content    string    `json:"content"`
	CreateDate string    `json:"create_date"`
}
