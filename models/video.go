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
	ID            int64
	Author        UserInfo
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
	Title         string
}

func NewVideoInfo(v *Video, ui *UserInfo,
	favoriteCount int64, commentCount int64) VideoInfo {
	return VideoInfo{
		ID:            int64((*v).ID),
		Author:        (*ui),
		PlayUrl:       (*v).PlayUrl,
		CoverUrl:      (*v).CoverUrl,
		FavoriteCount: favoriteCount,
		CommentCount:  commentCount,
		IsFavorite:    (favoriteCount > 0),
		Title:         (*v).Title,
	}
}
