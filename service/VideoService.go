package service

import (
	"douyin/models"
	"fmt"
	"time"
)

// MaxVideoNum 每次最多返回的视频流数量
const (
	MaxVideoNum = 30
)

type FeedVideoList struct {
	Videos   []*models.Video `json:"video_list,omitempty"`
	NextTime int64           `json:"next_time,omitempty"`
}

type QueryFeedVideoListFlow struct {
	userId     uint
	latestTime time.Time
	videos     []*models.Video
	nextTime   int64
	feedVideo  *FeedVideoList
}

// =============================================================================================================

func GetVideoById(vid uint) (models.Video, error) {
	video := models.Video{}
	err := models.DB.First(&video, vid).Error
	return video, err
}

// func GetVideosByIds(vids []uint) ([]models.Video, error) {
// 	videos := make([]models.Video, len(vids))
// 	err := models.DB.Where("vid in ?", vids).Find(&videos).Error
// 	return videos, err
// }

func GetVideosByUserId(uid uint) ([]models.Video, error) {
	if user, err := GetUserById(uid); err != nil {
		return []models.Video{}, fmt.Errorf("user not found: %v", err)
	} else {
		videos := make([]models.Video, 10)
		err := models.DB.Model(&user).Association("CreatedVideo").Find(&videos)
		return videos, err
	}
}

func CreateVideo(video models.Video) error {
	return models.DB.Create(&video).Error
}

// =============================================================================================================// =============================================================================================================// =============================================================================================================

func QueryFeedVideoList(userId uint, latestTime time.Time) (*FeedVideoList, error) {
	return NewQueryFeedVideoListFlow(userId, latestTime).Do()
}

func NewQueryFeedVideoListFlow(userId uint, latestTime time.Time) *QueryFeedVideoListFlow {
	return &QueryFeedVideoListFlow{userId: userId, latestTime: latestTime}
}

func (q *QueryFeedVideoListFlow) Do() (*FeedVideoList, error) {
	//所有传入的参数不填也应该给他正常处理
	q.checkNum()
	if err := q.packData(); err != nil {
		return nil, err
	}
	return q.feedVideo, nil
}

func (q *QueryFeedVideoListFlow) checkNum() {
	//上层通过把userId置零，表示userId不存在或不需要
	if q.userId > 0 {
		// 用户Id是有效的
	}
	if q.latestTime.IsZero() {
		q.latestTime = time.Now()
	}
}

func (q *QueryFeedVideoListFlow) packData() error {
	q.feedVideo = &FeedVideoList{
		Videos:   q.videos,
		NextTime: q.nextTime,
	}
	return nil
}
