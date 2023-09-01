package service

import (
	"context"
	"douyin/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// redis初始化和redis更新使用的函数

const (
	INTERACT_USER_FAVORITE_KEY     = "interact:favorite_videos:"
	INTERACT_VIDEO_FAVORITE_KEY    = "interact:favorited_by:"
	INTERACT_USER_TOT_FAVORITE_KEY = "interact:total_favorited:"
	SOCIAL_FOLLOWING_KEY           = "social:has_following:"
	SOCIAL_FOLLOWER_KEY            = "social:has_followers:"
	INTERACT_COMMENT_KEY           = "interact:has_comments:"
	BASIC_PUBLISH_KEY              = "basic:publish_works:"
	BASIC_RECENT_PUBLISH_KEY       = "basic:recent_publish:"
)

var RedisCtx = context.Background()

func Init2Redis() error {
	InitFavorite2Redis()
	InitFollow2Redis()
	InitComment2Redis()
	InitPublish2Redis()
	return nil
}

// 可以迁移到favorite service
func InitFavorite2Redis() {
	// 遍历favorite数据库
	var favorites []models.Favorite
	models.DB.Find(&favorites)
	for _, f := range favorites {
		// 取出uid，vid
		uid := f.UserId
		vid := f.VideoId
		// 更新三个key的内容
		models.RedisClient.SAdd(RedisCtx, INTERACT_USER_FAVORITE_KEY+strconv.Itoa(int(uid)), vid)
		models.RedisClient.SAdd(RedisCtx, INTERACT_VIDEO_FAVORITE_KEY+strconv.Itoa(int(vid)), uid)
		models.RedisClient.Incr(RedisCtx, INTERACT_USER_TOT_FAVORITE_KEY+strconv.Itoa(int(uid)))
	}
}

// 可以迁移到relation service
func InitFollow2Redis() {
	// 遍历relation表
	var relations []models.Relation
	models.DB.Find(&relations)
	for _, r := range relations {
		// 取出fromid, toid
		fromId := r.FollowerId
		toId := r.FollowedId
		// 更新两个key的内容
		models.RedisClient.SAdd(RedisCtx, SOCIAL_FOLLOWING_KEY+strconv.Itoa(int(fromId)), toId)
		models.RedisClient.SAdd(RedisCtx, SOCIAL_FOLLOWER_KEY+strconv.Itoa(int(toId)), fromId)
	}

}

type commentRelation struct {
	cid uint
	vid uint
}

// 可以迁移到comment service
func InitComment2Redis() {
	// 遍历comment表
	var commentRelations []commentRelation
	models.DB.Model(&models.Comment{}).Find(&commentRelations)
	for _, r := range commentRelations {
		// 取出vid
		vid := r.vid
		cid := r.cid
		// 更新一个key的内容
		models.RedisClient.SAdd(RedisCtx, INTERACT_COMMENT_KEY+strconv.Itoa(int(vid)), cid)
	}
}

type publishRelation struct {
	uid       uint
	vid       uint
	createdAt time.Time
}

func InitPublish2Redis() {
	var publishRelations []publishRelation
	models.DB.Model(&models.Video{}).Find(&publishRelations)
	for _, r := range publishRelations {
		uid := r.uid
		vid := r.vid
		ctime := r.createdAt
		models.RedisClient.SAdd(RedisCtx, BASIC_PUBLISH_KEY+strconv.Itoa(int(uid)), vid)
		models.RedisClient.ZAdd(RedisCtx, BASIC_PUBLISH_KEY+strconv.Itoa(int(uid)), &redis.Z{Score: float64(ctime.Unix()), Member: vid})
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
