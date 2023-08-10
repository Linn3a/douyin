package models


type Message struct {
	Id 		uint   `json:"id"`
	Content string `json:"content"`
	FromUserID int64 `json:"from_user_id"`
	ToUserID   int64 `json:"to_user_id"`
	CreateTime int64 `json:"create_time"`
}
