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
	curKey := INTERACT_USER_TOT_FAVORITE_KEY+strconv.Itoa(int(u.ID))
	if n, _ := models.RedisClient.Exists(RedisCtx, curKey).Result(); n == 0 {
		// 未命中则为0，因为当前场景下没有缓存过期
		u.TotalFavorited = 0
		return nil
	}
	totFavorited, err := models.RedisClient.Get(RedisCtx, curKey).Result()
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
	for i := 0; i < len(favoriteVids); i++ {
		tmp, _ := strconv.Atoi(favoriteVids[i])
		uintIds[i] = uint(tmp)
	}
	return uintIds, nil
}
//---------------------------------


// audience2video
func AddFavoriteVideo(uid uint, vid uint) error {
	if err := models.RedisClient.SAdd(RedisCtx, INTERACT_USER_FAVORITE_KEY+strconv.Itoa(int(uid)), vid).Err(); err != nil {
		return err
	}
	if err := models.RedisClient.SAdd(RedisCtx, INTERACT_VIDEO_FAVORITE_KEY+strconv.Itoa(int(vid)), uid).Err(); err != nil {
		return err
	}
	if err := models.RedisClient.Incr(RedisCtx, INTERACT_USER_TOT_FAVORITE_KEY+strconv.Itoa(int(uid))).Err(); err != nil {
		return err
	}
	// like消息加入消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(uid)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(vid)))
	rabbitmq.RmqLikeAdd.Publish(sb.String())
	fmt.Println("like消息入队成功")
	return nil
}

// audience2video
func DeleteFavoriteVideo(uid uint, vid uint) error {
	if err := models.RedisClient.SRem(RedisCtx, INTERACT_USER_FAVORITE_KEY+strconv.Itoa(int(uid)), vid).Err(); err != nil {
		return err
	}
	if err := models.RedisClient.SRem(RedisCtx, INTERACT_VIDEO_FAVORITE_KEY+strconv.Itoa(int(vid)), uid).Err(); err != nil {
		return err
	}
	if err := models.RedisClient.Decr(RedisCtx, INTERACT_USER_TOT_FAVORITE_KEY+strconv.Itoa(int(uid))).Err(); err != nil {
		return err
	}
	// like取消消息加入消息队列
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(uid)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(vid)))
	rabbitmq.RmqLikeDel.Publish(sb.String())
	fmt.Println("like取消消息入队成功")
	return nil
}



// audience2video
func GetFavoriteVideos(uid uint) ([]models.Video, error) {
	user := models.User{}
	user.ID = uid
	videos := make([]models.Video, 10)
	err := models.DB.Model(&user).Association("LikeVideo").Find(&videos)
	return videos, err
}

