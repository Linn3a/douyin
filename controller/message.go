package controller

import (
	"douyin/models"
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var tempChat = map[string][]models.Message{}

var messageIdSequence = int64(1)

type ChatResponse struct {
	Response
	MessageList []models.Message `json:"message_list"`
}

func MessageAction(c *fiber.Ctx) error {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")
	content := c.Query("content")

	if user, exist := usersLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(int64(user.ID), int64(userIdB))

		atomic.AddInt64(&messageIdSequence, 1)
		curMessage := models.Message{
			Model: gorm.Model{
				ID: uint(messageIdSequence),
				CreatedAt: time.Now(),
			},
			Content: content,
		}

		if messages, exist := tempChat[chatKey]; exist {
			tempChat[chatKey] = append(messages, curMessage)
		} else {
			tempChat[chatKey] = []models.Message{curMessage}
		}
		c.Status(http.StatusOK).JSON(Response{StatusCode: 0})
	} else {
		c.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	return c.JSON(fiber.Map{"MessageAction": "success"})
}

// MessageChat all users have same follow list
func MessageChat(c *fiber.Ctx) error {
	token := c.Query("token")
	toUserId := c.Query("to_user_id")

	if user, exist := usersLoginInfo[token]; exist {
		userIdB, _ := strconv.Atoi(toUserId)
		chatKey := genChatKey(int64(user.ID), int64(userIdB))

		c.Status(http.StatusOK).JSON(ChatResponse{Response: Response{StatusCode: 0}, MessageList: tempChat[chatKey]})
	} else {
		c.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	return c.JSON(fiber.Map{"MessageChat": "success"})
}

func genChatKey(userIdA int64, userIdB int64) string {
	if userIdA > userIdB {
		return fmt.Sprintf("%d_%d", userIdB, userIdA)
	}
	return fmt.Sprintf("%d_%d", userIdA, userIdB)
}