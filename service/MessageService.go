package service

import(
	"douyin/models"

)

func AddMessage(message models.Message) error {
	if err := models.DB.Table("messages").Create(&message).Error; err != nil {
		return err
	}
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