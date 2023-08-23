package models

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title         string    `json:"title"`
	PlayUrl       string    `json:"play_url"`  // 播放外链
	CoverUrl      string    `json:"cover_url"` // 封面外链
	AuthorID      uint      `json:"author_id"` // has many关系下游 在User中重设外键为AuthorID
	Comment       []Comment // has many关系上游 对应外键字段VideoID，引用ID
	FavoritedUser []*User   `gorm:"many2many:video_likes"`
}

type VideoInfo struct {
	ID            int64    `json:"id"`
	Author        *UserInfo `json:"author"`
	PlayUrl       string   `json:"play_url"`
	CoverUrl      string   `json:"cover_url"`
	FavoriteCount int64    `json:"favorite_count"`
	CommentCount  int64    `json:"comment_count"`
	IsFavorite    bool     `json:"is_favorite"`
	Title         string   `json:"title"`
}

func NewVideoInfo(v *Video) VideoInfo {
	return VideoInfo{
		ID:            int64((*v).ID),
		Author:        nil,
		PlayUrl:       (*v).PlayUrl,
		CoverUrl:      (*v).CoverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Title:         (*v).Title,
	}
}
