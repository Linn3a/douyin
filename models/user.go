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
	Follow []*User `gorm:"many2many:user_follows"`
	//  发出消息
	SendMessage []Message `gorm:"foreignKey:from_user_id"`
	//  收到消息
	ReceiveMessage []Message `gorm:"foreignKey:to_user_id"`
}
