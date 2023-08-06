package controller

import (
	"douyin/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CommentListResponse struct {
	Response
	CommentList []models.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment models.Comment `json:"comment,omitempty"`
}

func CommentAction(c *fiber.Ctx) error {
	token := c.Query("token")
	actionType := c.Query("action_type")

	if user, exist := usersLoginInfo[token]; exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			return c.Status(http.StatusOK).JSON(CommentActionResponse{Response: Response{StatusCode: 0},
				Comment: models.Comment{
					Model: gorm.Model{
						ID: 1,
						CreatedAt: time.Date(2023, 05, 01, 0, 0, 0, 0, time.Local),
					},
					UserId:       user.ID,
					Content:    text,
				}})

		}
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 0})
	} else {
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: DemoComments,
	})

}