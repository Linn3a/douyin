package service

import (
	"context"
	"douyin/models"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// redis初始化和redis更新使用的函数

const (
	// 加速关系查询
	INTERACT_USER_FAVORITE_KEY     = "interact:favorite_videos:"
	INTERACT_VIDEO_FAVORITE_KEY    = "interact:favorited_by:"
	INTERACT_USER_TOT_FAVORITE_KEY = "interact:total_favorited:"
	SOCIAL_FOLLOWING_KEY           = "social:has_following:"
	SOCIAL_FOLLOWER_KEY            = "social:has_followers:"
	// 加速总数统计
	INTERACT_COMMENT_KEY     = "interact:has_comments:"
	INTERACT_MAX_COMMENT_KEY = "interact:max_comment_id:"
	// 加速总数统计
	BASIC_PUBLISH_KEY = "basic:publish_works:"
	// 加速排序
	SOCIAL_MESSAGE_KEY = "social:messages:"
	// 加速排序
	BASIC_RECENT_PUBLISH_KEY = "basic:recent_publish:"
)

var RedisCtx context.Context

func Init2Redis() error {
	RedisCtx = context.Background()
	InitFavorite2Redis()
	InitFollow2Redis()
	InitComment2Redis()
	InitPublish2Redis()
	InitMessages2Redis()
	return nil
}

// InitFavorite2Redis 可以迁移到favorite service
func InitFavorite2Redis() {
	// 遍历favorite数据库
	var favorites []models.Favorite
	models.DB.Table("video_likes").Joins("join videos on videos.id = video_likes.video_id").Scan(&favorites)
	for _, f := range favorites {
		// 取出uid，vid
		uid := f.UserId
		vid := f.VideoId
		aid := f.AuthorId
		// 更新三个key的内容
		models.RedisClient.SAdd(RedisCtx, INTERACT_USER_FAVORITE_KEY+strconv.Itoa(int(uid)), vid)
		models.RedisClient.SAdd(RedisCtx, INTERACT_VIDEO_FAVORITE_KEY+strconv.Itoa(int(vid)), uid)
		models.RedisClient.Incr(RedisCtx, INTERACT_USER_TOT_FAVORITE_KEY+strconv.Itoa(int(aid)))
	}
}

// InitFollow2Redis 可以迁移到relation service
func InitFollow2Redis() {
	// 遍历relation表
	var relations []models.Relation
	models.DB.Table("user_follows").Find(&relations)
	for _, r := range relations {
		// 取出fromid, toid
		fromId := r.FollowerId
		toId := r.FollowedId
		// 更新两个key的内容
		models.RedisClient.SAdd(RedisCtx, SOCIAL_FOLLOWING_KEY+strconv.Itoa(int(fromId)), toId)
		models.RedisClient.SAdd(RedisCtx, SOCIAL_FOLLOWER_KEY+strconv.Itoa(int(toId)), fromId)
	}

}

type CommentRelation struct {
	ID      uint `json:"id"`
	VideoID uint `json:"video_id"`
}

// 可以迁移到comment service
func InitComment2Redis() {
	// 遍历comment表
	var commentRelations []CommentRelation
	models.DB.Model(&models.Comment{}).Find(&commentRelations)
	maxCid := uint(0)
	for _, r := range commentRelations {
		// 取出vid
		vid := r.VideoID
		cid := r.ID
		// panic(cid)
		fmt.Println(cid)
		if cid > maxCid {
			maxCid = cid
		}
		// 更新一个key的内容
		models.RedisClient.SAdd(RedisCtx, INTERACT_COMMENT_KEY+strconv.Itoa(int(vid)), cid)
	}
	models.RedisClient.Set(RedisCtx, INTERACT_MAX_COMMENT_KEY, maxCid, 0)
}

type publishRelation struct {
	AuthorID  uint
	ID        uint
	CreatedAt time.Time
}

func InitPublish2Redis() {
	var publishRelations []publishRelation
	models.DB.Model(&models.Video{}).Find(&publishRelations)
	for _, r := range publishRelations {
		uid := r.AuthorID
		vid := r.ID
		ctime := r.CreatedAt
		models.RedisClient.SAdd(RedisCtx, BASIC_PUBLISH_KEY+strconv.Itoa(int(uid)), vid)
		models.RedisClient.ZAdd(RedisCtx, BASIC_PUBLISH_KEY+strconv.Itoa(int(uid)), &redis.Z{Score: float64(ctime.UnixMilli()), Member: vid})
	}
}

type messageRelation struct {
	ID         uint
	FromUserID uint
	ToUserID   uint
	CreatedAt  time.Time
}

func InitMessages2Redis() {
	var messageRelations []messageRelation
	models.DB.Model(&models.Message{}).Find(&messageRelations)
	for _, r := range messageRelations {
		mid := r.ID
		fromid := r.FromUserID
		toid := r.ToUserID
		ctime := r.CreatedAt
		key := GenerateMessageKey(fromid, toid)
		models.RedisClient.ZAdd(RedisCtx, SOCIAL_MESSAGE_KEY+key, &redis.Z{Score: float64(ctime.UnixMilli()), Member: mid})
	}
}

// 如果使用mq实现异步更新关系表，则不需要使用定时任务
// func SyncRedisTables() {
// 	// 使用定时任务实现
// 	SyncFavoriteTables()
// 	SyncFollowTables()
// }

// func SyncFavorite2DB() {
// 	// 遍历user set或者video set
// 		// 取uid
// 		// 查所有关系
// 		// 遍历set内容
// 			// 取vid
// 			// 检查关系是否存在
// 				// 不存在则插入
// 				// 多余则删除
// }

// func SyncFollow2DB() {
// 	// 遍历 from set或者to set
// 		// 取 from id
// 		// 查所有关系
// 		// 遍历set内容
// 			// 取 to id
// 			// 检查关系是否存在
// 				// 不存在则插入
// 				// 多余则删除
// }
