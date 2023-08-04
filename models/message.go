package models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Content string `json:"content"`

	FromUserID int64 `json:"from_user_id"`
	ToUserID   int64 `json:"to_user_id"`
}
