package models

type Relation struct {
	FollowedId uint `json:"followed_id"`
	FollowerId uint `json:"follower_id"`
}
