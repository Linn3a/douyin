package models

import "time"

type Message struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`

	FromUserID int64 `json:"from_user_id" gorm:"foreignKey:User"`
	ToUserID   int64 `json:"to_user_id" gorm:"foreignKey:User"`
}
