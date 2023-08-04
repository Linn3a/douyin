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
