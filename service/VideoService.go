package service

import (
	"douyin/models"
	"douyin/utils/log"
	"fmt"
	"strconv"

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
	if err != nil {
		fmt.Printf("%v\n", err)
		return models.VideoInfo{}, err
	}
	videoInfo := GenerateVideoInfo(&video)
	authorInfo, err := GetUserInfoById(video.AuthorID)
	if err != nil {
		fmt.Printf("%v\n", err)
		return models.VideoInfo{}, err
	}
	if err := GetVideoFavoriteCount(&videoInfo); err != nil {
		fmt.Printf("%v\n", err)
		return models.VideoInfo{}, err
	}
	if err = GetVideoCommentCount(&videoInfo); err != nil {
		fmt.Printf("%v\n", err)
		return models.VideoInfo{}, err
	}
	videoInfo.Author = &authorInfo
	return videoInfo, nil
}

func GetVideoInfosByIds(vids []uint) ([]models.VideoInfo, error) {
	videoInfos := make([]models.VideoInfo, len(vids))
	videos, err := GetVideosByIds(vids)
	if err != nil {
		log.FieldLog("gorm", "error", "get video info by id failed")
		fmt.Printf("%v\n", err)
		return videoInfos, err
	}
	authorIds := make([]uint, len(vids))

	for ind, video := range videos {
		authorIds[ind] = video.AuthorID
		videoInfos[ind] = GenerateVideoInfo(&video)
	}
	authorInfos, err := GetUserInfoMapByIds(authorIds)
	if err != nil {
		return videoInfos, err
	}

	for i := 0; i < len(videoInfos); i++ {
		if err := GetVideoFavoriteCount(&videoInfos[i]); err != nil {
			fmt.Printf("%v\n", err)
			return videoInfos, err
		}
		if err := GetVideoCommentCount(&videoInfos[i]); err != nil {
			fmt.Printf("%v\n", err)
			return videoInfos, err
		}
		authorInfo := authorInfos[videos[i].AuthorID]
		videoInfos[i].Author = &authorInfo
	}
	return videoInfos, nil
}

func GetFeedVideoIds(latest_time *int64) ([]uint, error) {

	zVids, err := models.RedisClient.ZRevRangeByScoreWithScores(RedisCtx, BASIC_RECENT_PUBLISH_KEY, &redis.ZRangeBy{
		Min:   "0",
		Max:   strconv.FormatInt(*latest_time, 10),
		Count: MaxVideoNum,
	}).Result()
	if err != nil {
		log.FieldLog("redis", "error", "get video feed ids failed")
		return []uint{}, err
	}

	*latest_time = int64(zVids[len(zVids)-1].Score)
	vids := make([]uint, len(zVids))
	for i, zVid := range zVids {
		tmpStr := zVid.Member.(string)
		tmpInt, _ := strconv.Atoi(tmpStr)
		tmp := uint(tmpInt)
		vids[i] = tmp
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
	err := models.DB.Where("id in ?", vids).Find(&videos).Error
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
		Title:    title,
		PlayUrl:  playUrl,
		CoverUrl: coverUrl,
		AuthorID: uid,
	}
	err := models.DB.Create(&video).Error
	models.RedisClient.SAdd(RedisCtx, BASIC_PUBLISH_KEY+strconv.Itoa(int(uid)), video.ID)
	models.RedisClient.ZAdd(RedisCtx, BASIC_RECENT_PUBLISH_KEY, &redis.Z{Score: float64(video.CreatedAt.UnixMilli()), Member: video.ID})
	return err
}
