package service

import (
	"douyin/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	MaxMessageCount = 30
)

func AddMessage(toId uint, fromId uint, content string) error {
	message := models.Message{
		ToUserID:   toId,
		FromUserID: fromId,
		Content:    content,
		CreateTime: time.Now(),
		// CreateTime: time.Now(),
	}
	if err := models.DB.Table("messages").Create(&message).Error; err != nil {
		return err
	}
	key := GenerateMessageKey(message.FromUserID, message.ToUserID)
	models.RedisClient.ZAdd(RedisCtx, SOCIAL_MESSAGE_KEY+key, &redis.Z{Score: float64(message.CreateTime.UnixMilli()), Member: message.ID})
	return nil
}

// GetLatestMessageAfter 获取 from 和 to 之间 在 preMsgTime之后的最近消息
func GetLatestMessageAfter(fromId uint, toId uint, preMsgTime int64) ([]models.Message, error) {
	var msgList []models.Message
	err := models.DB.Table("messages").
		Where("(from_user_id = ? AND to_user_id = ? ) OR (from_user_id = ? AND to_user_id = ?)", fromId, toId, toId, fromId).
		Order("create_time asc").Where("create_time > ?", preMsgTime).Find(&msgList).Error

	return msgList, err
}

func GenerateMessageKey(fromId uint, toId uint) string {
	var key string
	if fromId < toId {
		key = strconv.Itoa(int(fromId)) + "-" + strconv.Itoa(int(toId))
	} else {
		key = strconv.Itoa(int(toId)) + "-" + strconv.Itoa(int(fromId))
	}
	return key
}

func GetMessagesIds(fromId uint, toId uint, preMsgTime *int64) ([]uint, error) {
	key := GenerateMessageKey(fromId, toId)
	// fmt.Println(key, strconv.FormatInt(*preMsgTime, 10), strconv.FormatInt(time.Now().UnixMilli(), 10))
	zMids, err := models.RedisClient.ZRangeByScoreWithScores(RedisCtx, SOCIAL_MESSAGE_KEY+key, &redis.ZRangeBy{
		Min:   strconv.FormatInt(*preMsgTime, 10),
		Max:   strconv.FormatInt(time.Now().UnixMilli(), 10),
		Count: MaxMessageCount,
	}).Result()
	if err != nil {
		return []uint{}, err
	}
	if len(zMids) == 0 {
		return []uint{}, nil
	}
	*preMsgTime = int64(zMids[len(zMids)-1].Score)
	mids := make([]uint, len(zMids))
	for i, zMid := range zMids {
		tmpStr := zMid.Member.(string)
		tmpInt, _ := strconv.Atoi(tmpStr)
		tmp := uint(tmpInt)
		// fmt.Println(tmp)
		mids[i] = tmp
	}
	return mids, nil
}

func GetMessagesByIds(mids []uint) ([]models.Message, error) {
	messages := make([]models.Message, len(mids))
	err := models.DB.Where("id in (?)", mids).Find(&messages).Error
	return messages, err
}

func GenerateMessageInfo(message *models.Message) models.MessageInfo {
	messageInfo := models.MessageInfo{
		ID:         int64(message.ID),
		ToUserID:   int64(message.ToUserID),
		FromUserID: int64(message.FromUserID),
		Content:    message.Content,
		// CreateTime: message.CreatedAt.Format("yyyy-MM-dd HH:MM:ss"),
		CreateTime: int(message.CreatedAt.UnixMilli()),
	}
	return messageInfo
}

func GetFriendLatestMessageId(fromId uint, toId uint) (uint, error) {
	key := GenerateMessageKey(fromId, toId)
	zMids, err := models.RedisClient.ZRevRangeByScoreWithScores(RedisCtx, SOCIAL_MESSAGE_KEY+key, &redis.ZRangeBy{
		Min:   "0",
		Max:   strconv.FormatInt(time.Now().UnixMilli(), 10),
		Count: 1,
	}).Result()
	mid := zMids[0].Member.(uint)
	return mid, err
}

// func GetMessageInfoById(mid uint) (models.MessageInfo, error) {
// 	message := models.Message{}
// 	err := models.DB.First(&message, mid).Error
// 	messageInfo := GenerateMessageInfo(&message)
// 	return messageInfo, err
// }

func GetFriendLatestMessageInfo(fromId uint, toId uint) (models.MessageInfo, error) {
	key := GenerateMessageKey(fromId, toId)
	zMids, err := models.RedisClient.ZRevRangeByScore(RedisCtx, SOCIAL_MESSAGE_KEY+key, &redis.ZRangeBy{
		Min:   "0",
		Max:   strconv.FormatInt(time.Now().UnixMilli(), 10),
		Count: 1,
	}).Result()
	if err != nil {
		return models.MessageInfo{}, err
	}
	if len(zMids) == 0 {
		return models.MessageInfo{}, nil
	}
	midInt, _ := strconv.Atoi(zMids[0])
	mid := uint(midInt)
	message := models.Message{}
	err = models.DB.First(&message, mid).Error
	messageInfo := GenerateMessageInfo(&message)
	return messageInfo, err
}
