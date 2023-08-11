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
	ID   int64 				`json:"id"`
	Name string				`json:"name"`
	FollowCount     int64	`json:"follow_count"`
	FollowerCount   int64	`json:"follower_count"`
	IsFollow        bool	`json:"is_follow"`
	Avatar           string	`json:"avatar"`
	BackgroundImage string	`json:"background_image"`
	Signature        string	`json:"signature"`
	TotalFavorited  int64	`json:"total_favorited"`
	WorkCount       int64	`json:"work_count"`
	FavoriteCount   int64	`json:"favorite_count"`
}

func NewUserInfo(u *User) UserInfo {
	return UserInfo{
		ID:               int64((*u).ID),
		Name:             (*u).Name,
		FollowCount:     0,
		FollowerCount:   0,
		IsFollow:        false,
		Avatar:           (*u).Avatar,
		BackgroundImage: (*u).BackgroundImage,
		Signature:        (*u).Signature,
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:	  0,
	}
}