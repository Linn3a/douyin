package service

import (
	"douyin/models"
	"strings"
	"strconv"
	"douyin/middleware/rabbitmq"
	"fmt"
)

// redis 关系查询优化

func GetVideoIsFavorite(v *models.VideoInfo, uid uint) error {
	isFavorite, err := models.RedisClient.SIsMember(RedisCtx, INTERACT_VIDEO_FAVORITE_KEY+strconv.Itoa(int(v.ID)), uid).Result()
	if err != nil {
		return fmt.Errorf("user favorite set check error: %v", err)
	}
	v.IsFavorite = isFavorite
	return nil
}

func GetVideoFavoriteCount(v *models.VideoInfo) error {
	favoriteCount, err := models.RedisClient.SCard(RedisCtx, INTERACT_VIDEO_FAVORITE_KEY+strconv.Itoa(int(v.ID))).Result()
	if err != nil {
		return fmt.Errorf("video favorited set count error: %v", err)
	}
	v.FavoriteCount = favoriteCount
	return nil
}

func GetUserFavoriteCount(u *models.UserInfo) error {
	favoriteCount, err := models.RedisClient.SCard(RedisCtx, INTERACT_USER_FAVORITE_KEY+strconv.Itoa(int(u.ID))).Result()
	if err != nil {
		return fmt.Errorf("user favorited set count error: %v", err)
	}
	u.FavoriteCount = favoriteCount
	return nil
}

func GetUserTotalFavorited(u *models.UserInfo) error {
	totFavorited, err := models.RedisClient.Get(RedisCtx, INTERACT_USER_TOT_FAVORITE_KEY+strconv.Itoa(int(u.ID))).Result()
	if err != nil {
		return fmt.Errorf("user tot favorited get error: %v", err)
	}
	numTotFavorited, _ := strconv.Atoi(totFavorited)
	u.TotalFavorited = int64(numTotFavorited)
	return nil
}

func GetFavoriteVideoIds(uid uint) ([]uint, error) {
	favoriteVids, err := models.RedisClient.SMembers(RedisCtx, INTERACT_USER_FAVORITE_KEY+strconv.Itoa(int(uid))).Result()
	if err != nil {
		return []uint{}, err
	}
	uintIds := make([]uint, len(favoriteVids))
	for i, _ := range favoriteVids {
		tmp, _ := strconv.Atoi(favoriteVids[i])
		uintIds[i] = uint(tmp)
	}
	return uintIds, nil
}
//---------------------------------


// audience2video
func AddFavoriteVideo(uid uint, vid uint) error {
	err := models.RedisClient.SAdd(RedisCtx, INTERACT_USER_FAVORITE_KEY+strconv.Itoa(int(uid)), vid).Err()
	err = models.RedisClient.SAdd(RedisCtx, INTERACT_VIDEO_FAVORITE_KEY+strconv.Itoa(int(vid)), uid).Err()
	err = models.RedisClient.Incr(RedisCtx, INTERACT_USER_TOT_FAVORITE_KEY+strconv.Itoa(int(uid))).Err()
	// like消息加入消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(uid)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(vid)))
	rabbitmq.RmqLikeAdd.Publish(sb.String())
	fmt.Println("like消息入队成功")
	return err
	// user, err := GetUserById(uid)
	// if err != nil {
	// 	return fmt.Errorf("user not found: %v", err)
	// }
	// video, err := GetVideoById(vid)
	// if err != nil {
	// 	return fmt.Errorf("video not found: %v", err)
	// }
	// user := models.User{}
	// user.ID = uid
	// video := models.Video{}
	// video.ID = vid
	// err := models.DB.Model(&user).Association("LikeVideo").Append(&video)
	// return err
}

// audience2video
func DeleteFavoriteVideo(uid uint, vid uint) error {
	err := models.RedisClient.SRem(RedisCtx, INTERACT_USER_FAVORITE_KEY+strconv.Itoa(int(uid)), vid).Err()
	err = models.RedisClient.SRem(RedisCtx, INTERACT_VIDEO_FAVORITE_KEY+strconv.Itoa(int(vid)), uid).Err()
	err = models.RedisClient.Decr(RedisCtx, INTERACT_USER_TOT_FAVORITE_KEY+strconv.Itoa(int(uid))).Err()
	// like取消消息加入消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(uid)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(vid)))
	rabbitmq.RmqLikeDel.Publish(sb.String())
	fmt.Println("like取消消息入队成功")
	return err
	// user, err := GetUserById(uid)
	// if err != nil {
	// 	return fmt.Errorf("user not found: %v", err)
	// }
	// video, err := GetVideoById(vid)
	// if err != nil {
	// 	return fmt.Errorf("video not found: %v", err)
	// }
	// user := models.User{}
	// user.ID = uid
	// video := models.Video{}
	// video.ID = vid
	// err := models.DB.Model(&user).Association("LikeVideo").Delete(video)
	// return err
}



// audience2video
func GetFavoriteVideos(uid uint) ([]models.Video, error) {
	user := models.User{}
	user.ID = uid
	videos := make([]models.Video, 10)
	err := models.DB.Model(&user).Association("LikeVideo").Find(&videos)
	return videos, err
}

