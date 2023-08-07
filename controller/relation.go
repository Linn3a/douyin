package controller

import (
	"douyin/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserListResponse struct {
	Response
	UserList []models.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *fiber.Ctx) error {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 0})
	} else {
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []models.User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []models.User{DemoUser},
	})
}

// FriendList all users have same friend list
func FriendList(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []models.User{DemoUser},
	})
}