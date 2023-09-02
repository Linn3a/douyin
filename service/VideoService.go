package service

import (
	"douyin/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// MaxVideoNum 每次最多返回的视频流数量
const (
	MaxVideoNum = 30
)
// redis 缓存查询

func GetUserWorkCount(u *models.UserInfo) error {
	count, err := models.RedisClient.SCard(RedisCtx, BASIC_PUBLISH_KEY+strconv.Itoa(int(u.ID))).Result()
	if err != nil {
		return err
	}
	u.WorkCount = count
	return nil
}

func GetVideoIdsByUserId(uid uint) ([]uint, error) {
	strVids, err := models.RedisClient.SMembers(RedisCtx, BASIC_PUBLISH_KEY+strconv.Itoa(int(uid))).Result()
	if err != nil {
		return []uint{}, err
	}
	vids := make([]uint, len(strVids))
	for i, svid := range strVids {
		tmp, _ := strconv.Atoi(svid)
		vids[i] = uint(tmp)
	}
	return vids, nil
}

// ----------------------------------

func GenerateVideoInfo(v *models.Video) models.VideoInfo {
	return models.VideoInfo{
		ID:            int64(v.ID),
		Author:        nil,
		PlayUrl:       v.PlayUrl,
		CoverUrl:      v.CoverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Title:         v.Title,
	}
}

func GetVideoInfoById(vid uint) (models.VideoInfo, error) {
	video, err := GetVideoById(vid)
	videoInfo := GenerateVideoInfo(&video)
	authorInfo, err := GetUserInfoById(video.AuthorID)
	err = GetVideoFavoriteCount(&videoInfo)
	err = GetVideoCommentCount(&videoInfo)
	videoInfo.Author = &authorInfo
	return videoInfo, err
}

func GetVideoInfosByIds(vids []uint) ([]models.VideoInfo, error) {
	videos, err := GetVideosByIds(vids)
	videoInfos := make([]models.VideoInfo, len(vids))
	authorIds := make([]uint, len(vids))

	for ind, video := range videos {
		authorIds[ind] = video.AuthorID
		videoInfos[ind] = GenerateVideoInfo(&video)
	}
	authorInfos, err := GetUserInfosByIds(authorIds)

	for i, videoInfo := range videoInfos {
		err = GetVideoFavoriteCount(&videoInfos[i])
		err = GetVideoCommentCount(&videoInfos[i])
		authorInfo := authorInfos[uint(videoInfo.ID)]
		videoInfos[i].Author = &authorInfo
	}
	return videoInfos, err
}

func GetFeedVideoIds(latest_time int64) ([]uint, error) {
	var maxTime int64
	if latest_time == 0 {
		maxTime = time.Now().Unix()
	} else {
		maxTime = latest_time
	}
	
	strVids, err := models.RedisClient.ZRevRangeByScore(RedisCtx, BASIC_RECENT_PUBLISH_KEY, &redis.ZRangeBy{
		Min: "0",
		Max: strconv.FormatInt(maxTime, 10),
		Count: MaxVideoNum,
	}).Result()
	if err != nil {
		return []uint{}, err
	}
	vids := make([]uint, len(strVids))
	for i, strVid := range strVids {
		tmp, _ := strconv.Atoi(strVid)
		vids[i] = uint(tmp)
	}
	return vids, nil
} 

// =============================================================================================================

func GetVideoById(vid uint) (models.Video, error) {
	video := models.Video{}
	err := models.DB.First(&video, vid).Error
	return video, err
}

func GetVideosByIds(vids []uint) ([]models.Video, error) {
	videos := make([]models.Video, len(vids))
	err := models.DB.Where("vid in ?", vids).Find(&videos).Error
	return videos, err
}
func GetVideosByUpdateAt() ([]models.Video, error) {
	videos := make([]models.Video, 10)
	err := models.DB.Order("updated_at desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, err
}

func GetVideosLenght() (int64, error) {
	var video models.Video
	var lenght int64
	err := models.DB.Model(&video).Count(&lenght).Error
	if err != nil {
		return lenght, err
	}
	return lenght, err
}

// func GetVideosByIds(vids []uint) ([]models.Video, error) {
// 	videos := make([]models.Video, len(vids))
// 	err := models.DB.Where("vid in ?", vids).Find(&videos).Error
// 	return videos, err
// }

func CreateVideo(title string, playUrl string, coverUrl string, uid uint) error {
	video := models.Video{
		Title: title,
		PlayUrl: playUrl,
		CoverUrl: coverUrl,
		AuthorID: uid,
	}
	err := models.DB.Create(&video).Error
	models.RedisClient.SAdd(RedisCtx, BASIC_PUBLISH_KEY+strconv.Itoa(int(uid)), video.ID)
	models.RedisClient.ZAdd(RedisCtx, BASIC_RECENT_PUBLISH_KEY, &redis.Z{Score: float64(video.CreatedAt.Unix()), Member: video.ID})
	return err
}
