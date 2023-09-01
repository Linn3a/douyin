package controller

import (
	"douyin/models"
	"douyin/service"
	"douyin/utils/jwt"
	"douyin/utils/validator"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

var tempChat = map[string][]models.Message{}

var messageIdSequence = int64(1)

type MessageActionRequest struct {
	Token      string `query:"token" validate:"required"`
	ToUserID   string `query:"to_user_id" validate:"required"`
	ActionType string `query:"action_type" validate:"required, oneof=1"`
	Content    string `query:"content" validate:"required"`
}

type MessageChatRequest struct {
	Token    string `query:"token" validate:"required"`
	ToUserID string `query:"to_user_id" validate:"required"`
}

type ChatResponse struct {
	Response
	MessageList []models.Message `json:"message_list"`
	PreMsgTime  int64            `json:"pre_msg_time"`
}

func MessageAction(c *fiber.Ctx) error {

	request := MessageActionRequest{}
	emptyResponse := Response{}
	if err := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return err
	}
	var fromId uint
	if err := jwt.JwtClient.AuthTokenValid(c, &emptyResponse, &fromId, request.Token); err != nil {
		return err
	}
	toUserId, _ := strconv.Atoi(request.ToUserID)
	content := request.Content

	// chatKey := genChatKey(int64(fromUserId), int64(toUserId))
	curMessage := models.Message{
		CreateTime: int64(time.Now().Unix()), //以秒为时间单位
		Content:    content,
		FromUserID: int64(fromId),
		ToUserID:   int64(toUserId),
	}
	service.AddMessage(curMessage)
	return c.Status(http.StatusOK).JSON(Response{StatusCode: 0, StatusMsg: "发送消息成功！"})

}

// MessageChat all users have same follow list
func MessageChat(c *fiber.Ctx) error {

	request := MessageActionRequest{}
	emptyResponse := Response{}
	if err := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return err
	}
	var fromId uint
	if err := jwt.JwtClient.AuthTokenValid(c, &emptyResponse, &fromId, request.Token); err != nil {
		return err
	}
	toUserId, _ := strconv.Atoi(request.ToUserID)


	// 上次消息时间
	var preMsgTime int64
	preMsgTimeStr := c.Query("pre_msg_time")
	if preMsgTimeStr == "" {
		preMsgTime = 1546926630
	} else {
		preMsgTime, _ = strconv.ParseInt(preMsgTimeStr, 10, 64)
	}
	msgList, err := service.GetLatestMessageAfter(fromId, uint(toUserId), preMsgTime)
	// 无消息
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(ChatResponse{
			Response:    Response{StatusCode: 1, StatusMsg: "no message"},
			MessageList: nil,
			PreMsgTime:  1546926630,
		})
	}
	var nextPreMsgTime int64
	if len(msgList) == 0 {
		nextPreMsgTime = 1546926630
	} else {
		nextPreMsgTime = msgList[len(msgList)-1].CreateTime
	}
	return c.Status(fiber.StatusOK).JSON(ChatResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "成功获取消息！"},
		MessageList: msgList,
		PreMsgTime:  nextPreMsgTime,
	})
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}
