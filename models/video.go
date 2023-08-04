package models

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title    string    `json:"title"`
	PlayUrl  string    `json:"play_url"`  // 播放外链
	CoverUrl string    `json:"cover_url"` // 封面外链
	AuthorID uint      `json:"author_id"` // has many关系下游 在User中重设外键为AuthorID
	Comment  []Comment // has many关系上游 对应外键字段VideoID，引用ID
}