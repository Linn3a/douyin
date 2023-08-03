package models

import "time"

type User struct {
	// base info
	ID              int64     `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name"`
	Password        string    `json:"password"`
	Avatar          string    `json:"avatar"`
	BackgroundImage string    `json:"background_image"`
	Signature       string    `json:"signature"`
	CreatedAt       time.Time `json:"created_at"`

	//	自引用follow
	Follow []*User `gorm:"many2many:user_follows"`
	//  发出消息
	SendMessage []Message `gorm:"foreignKey:from_user_id"`
	//  收到消息
	ReceiveMessage []Message `gorm:"foreignKey:to_user_id"`
}
