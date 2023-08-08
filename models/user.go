package models

import (
	"gorm.io/gorm"
)

type User struct {
	// base info
	gorm.Model
	Name            string `json:"name"`
	Password        string `json:"password"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`

	//	自引用follow
	Follow []*User `gorm:"many2many:user_follows;joinForeignKey:followed_id;joinReferences:follower_id"`
	// followed: 被关注人
	// follower: 粉丝

	//  发出消息
	SendMessage []Message `gorm:"foreignKey:from_user_id"`
	//  收到消息
	ReceiveMessage []Message `gorm:"foreignKey:to_user_id"`
	//  发表的评论
	Comment []Comment
	//  发布的视频
	CreatedVideo []Video `gorm:"foreignKey:AuthorID"`
	//  喜欢的视频
	LikeVideo []Video `gorm:"many2many:video_likes"` // many2many关系 连接表名video_likes
}

// 用于响应http请求的结构
// TODO: add info from other service
type UserInfo struct {
	id   int64
	name string
	// follow_count     int64
	// follower_count   int64
	// is_follow        bool
	avatar           string
	background_image string
	signature        string
	// total_favorited  int64
	// work_count       int64
	// favorite_count   int64
}

func NewUserInfo(u *User) UserInfo {
	return UserInfo{
		id:               int64((*u).ID),
		name:             (*u).Name,
		avatar:           (*u).Avatar,
		background_image: (*u).BackgroundImage,
		signature:        (*u).Signature,
	}
}