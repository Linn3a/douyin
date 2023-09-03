package models

import (
	"time"

	"gorm.io/gorm"
)

// type Message struct {
// 	Id 		uint   `json:"id"`
// 	Content string `json:"content"`
// 	FromUserID int64 `json:"from_user_id"`
// 	ToUserID   int64 `json:"to_user_id"`
// 	CreateTime int64 `json:"create_time"`
// }

type Message struct {
	gorm.Model
	Content    string `json:"content"`
	FromUserID uint   `json:"from_user_id"`
	ToUserID   uint   `json:"to_user_id"`
	CreateTime time.Time `json:"create_time"`
}

type MessageInfo struct {
	ID         int64  `json:"id"`
	ToUserID   int64  `json:"to_user_id"`
	FromUserID int64  `json:"from_user_id"`
	Content    string `json:"content"`
	// CreateTime string `json:"create_time"`
	CreateTime int `json:"create_time"`
}
