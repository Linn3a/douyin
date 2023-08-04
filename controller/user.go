package controller

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"sync/atomic"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *fiber.Ctx) error {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	if _, exist := usersLoginInfo[token]; exist {
		return c.Status(http.StatusOK).JSON(UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := User{
			Id:   userIdSequence,
			Name: username,
		}
		usersLoginInfo[token] = newUser
		return c.Status(http.StatusOK).JSON(UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    username + password,
		})
	}
}

func Login(c *fiber.Ctx) error {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	if user, exist := usersLoginInfo[token]; exist {
		return c.Status(http.StatusOK).JSON(UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	} else {
		return c.Status(http.StatusOK).JSON(UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

func UserInfo(c *fiber.Ctx) error {
	token := c.Query("token")

	if user, exist := usersLoginInfo[token]; exist {
		return c.Status(http.StatusOK).JSON(UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		return c.Status(http.StatusOK).JSON(UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}
}

//func Register(c *gin.Context) {
//	username := c.Query("username")
//	password := c.Query("password")
//
//	token := username + password
//
//	if _, exist := usersLoginInfo[token]; exist {
//		c.JSON(http.StatusOK, UserLoginResponse{
//			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
//		})
//	} else {
//		atomic.AddInt64(&userIdSequence, 1)
//		newUser := User{
//			Id:   userIdSequence,
//			Name: username,
//		}
//		usersLoginInfo[token] = newUser
//		c.JSON(http.StatusOK, UserLoginResponse{
//			Response: Response{StatusCode: 0},
//			UserId:   userIdSequence,
//			Token:    username + password,
//		})
//	}
//}
//
//func Login(c *gin.Context) {
//	username := c.Query("username")
//	password := c.Query("password")
//
//	token := username + password
//
//	if user, exist := usersLoginInfo[token]; exist {
//		c.JSON(http.StatusOK, UserLoginResponse{
//			Response: Response{StatusCode: 0},
//			UserId:   user.Id,
//			Token:    token,
//		})
//	} else {
//		c.JSON(http.StatusOK, UserLoginResponse{
//			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
//		})
//	}
//}
//
//func UserInfo(c *gin.Context) {
//	token := c.Query("token")
//
//	if user, exist := usersLoginInfo[token]; exist {
//		c.JSON(http.StatusOK, UserResponse{
//			Response: Response{StatusCode: 0},
//			User:     user,
//		})
//	} else {
//		c.JSON(http.StatusOK, UserResponse{
//			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
//		})
//	}
//}
