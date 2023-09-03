package models

type Favorite struct{
	UserId uint `json:"user_id"`
	VideoId uint `json:"video_id"`
	AuthorId uint `json:"author_id`
}