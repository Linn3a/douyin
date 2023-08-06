package service

import (
	"github.com/gofiber/fiber/v2"
	// "gorm.io/gorm"
	// "sync/atomic"
	"github.com/Eacient/douyin/models"
	"github.com/Eacient/douyin/public"
	"fmt"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = make(map[string]models.User)


type UserLoginResponse struct {
	models.Response
	UserId uint  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	models.Response
	User models.User `json:"user"`
}

func Register(c *fiber.Ctx) error {
	username := c.Query("username")
	password := c.Query("password")

	var tmp = models.User{}
	if err := public.DBConn.Where("name = ?", username).First(&tmp).Error;err == nil {
		fmt.Println("The user exist")
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	}
	newUser := models.User{
		Name: username,
		Password:password,
	}
	if err := public.DBConn.Create(&newUser).Error; err != nil {
		fmt.Println("插入失败", err)
	}else{
		fmt.Println("插入成功", err)
		public.DBConn.Where("name = ?", username).First(&tmp)
		fmt.Println("id:",tmp.ID)
	}
	return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
		Response: models.Response{StatusCode: 0},
		UserId:   tmp.ID,
		Token:    username + password,
	})
	
}

func Login(c *fiber.Ctx) error {
	username := c.Query("username")
	password := c.Query("password")

	var tmp = models.User{}
	if err := public.DBConn.Where("name = ?", username).First(&tmp).Error;err != nil {
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: models.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
	}else{
		TmpPassword := tmp.Password
		if TmpPassword == password{
			return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
				Response: models.Response{StatusCode: 0},
				UserId:   tmp.ID,
				Token:    username + password,
			})
		}else {
			return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
				Response: models.Response{StatusCode: 1, StatusMsg: "Password doesn't match"},
			})
		}
	}
}

func UserInfo(c *fiber.Ctx) error {
	token := c.Query("token")
	uid := c.Query("user_id")
	
	var tmp = models.User{}
	if err := public.DBConn.Where("id = ?", uid).First(&tmp).Error;err == nil {
		username := tmp.Name
		password := tmp.Password
		TrueToken := username + password
		if TrueToken == token{
			return c.Status(fiber.StatusOK).JSON(UserResponse{
				Response: models.Response{StatusCode: 0},
				User:     tmp,
			})
		}
	}else{
		fmt.Println("UserInfo error",err)
	}
	return c.Status(fiber.StatusOK).JSON(UserResponse{
		Response: models.Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
	})
	
	
}

