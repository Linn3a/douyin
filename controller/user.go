package controller

import (
	"douyin/models"
	"douyin/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]models.User{
	"zhangleidouyin": {
		Model: gorm.Model{
			ID: 1,
		},
		Name:     "zhanglei",
		Password: "douyin",
		// Id:            1,
		// Name:          "zhanglei",
		// FollowCount:   10,
		// FollowerCount: 5,
		// IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User models.User `json:"user"`
}

func Register(c *fiber.Ctx) error {
	username := c.Query("username")
	password := c.Query("password")

	_, err := service.GetUserByName(username)
	if err == nil {
		fmt.Println("The suer exits")
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User already exist",
			},
		})
	}

	newUser := models.User{
		Name:     username,
		Password: password,
	}
	err = service.CreateUser(&newUser)
	if err != nil {
		fmt.Println("插入失败", err)
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 2,
				StatusMsg:  "User insertion error",
			},
		})
	}

	fmt.Println("插入成功")
	return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserId: int64(newUser.ID),
		Token:  username + password,
	})
}

func Login(c *fiber.Ctx) error {
	username := c.Query("username")
	password := c.Query("password")

	user, err := service.GetUserByName(username)

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User doesn't exist",
			},
		})
	}

	if user.Password != password {
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "Password doesn't match",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserId: int64(user.ID),
		Token:  username + password,
	})
}

func UserInfo(c *fiber.Ctx) error {
	token := c.Query("token")
	uid, _ := strconv.Atoi(c.Query("user_id"))

	if user, err := service.GetUserById(uint(uid)); err != nil && token == token {
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User not exits",
			},
		})
	} else {
		return c.Status(http.StatusOK).JSON(
			UserResponse{
				Response: Response{StatusCode: 0},
				User:     user,
			},
		)
	}
}
