package controller

import (
	"douyin/models"
	"douyin/service"
	"fmt"
	"strconv"

	"douyin/utils/jwt"
	"douyin/utils/validator"

	"github.com/gofiber/fiber/v2"
)

type UserRegisterRequest struct {
	Username string `query:"username" validate:"required"`
	Password string `query:"password" validate:"required"`
}

type UserLoginRequest struct {
	Username string `query:"username" validate:"required"`
	Password string `query:"password" validate:"required"`
}

type UserRequest struct {
	UserID string `query:"user_id" validate:"required"`
	Token  string `query:"token" validate:"required"` // 用户鉴权token
}

type UserLoginResponse struct {
	Response
	UserID int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User models.UserInfo `json:"user"`
}

func Register(c *fiber.Ctx) error {

	request := UserRegisterRequest{}
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &UserLoginResponse{}, &request); err != nil {
		return httpErr
	}
	username := request.Username
	password := request.Password

	if _, err := service.GetUserByName(username); err == nil {
		fmt.Println("The user exits")
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 3,
				StatusMsg:  "User already exist",
			},
			UserID: 1,
			Token:  "",
		})
	}

	newUser := models.User{
		Name:     username,
		Password: password,
	}

	if err := service.CreateUser(&newUser); err != nil {
		fmt.Println("插入失败", err)
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 4,
				StatusMsg:  "User insertion error",
			},
			UserID: int64(newUser.ID),
			Token:  "",
		})
	}
	fmt.Println("插入成功")

	token, err := service.GenerateToken(&newUser)
	if err != nil {
		fmt.Println("创建token失败")
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 5,
				StatusMsg:  "Unable to create token",
			},
			UserID: int64(newUser.ID),
		})
	}

	fmt.Printf("token is : %s", token)
	return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserID: int64(newUser.ID),
		Token:  token,
	})
}

func Login(c *fiber.Ctx) error {
	request := UserLoginRequest{}
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &UserLoginResponse{}, &request); err != nil {
		return httpErr
	}
	username := request.Username
	password := request.Password

	user, err := service.GetUserByName(username)

	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 3,
				StatusMsg:  "User doesn't exist",
			},
			UserID: 1,
			Token:  "",
		})
	}

	if user.Password != password {
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 4,
				StatusMsg:  "Password doesn't match",
			},
			UserID: 1,
			Token:  "",
		})
	}

	token, err := service.GenerateToken(&user)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
			Response: Response{
				StatusCode: 5,
				StatusMsg:  "Unable to create token",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(UserLoginResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "Login successfully",
		},
		UserID: int64(user.ID),
		Token:  token,
	})
}

func UserInfo(c *fiber.Ctx) error {
	request := UserRequest{}
	emptyResponse := UserResponse{}
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return httpErr
	}
	uid, _ := strconv.Atoi(request.UserID)
	if err, httpErr := jwt.JwtClient.AuthCurUser(c, &emptyResponse, request.Token, uint(uid)); err != nil {
		return httpErr
	}
	user, err := service.GetUserById(uint(uid))
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserResponse{
			Response: Response{StatusCode: 5, StatusMsg: "user not exits"}})
	}
	userInfo := service.GenerateUserInfo(&user)
	err = service.GetUserFollowCount(&userInfo)
	err = service.GetUserFollowerCount(&userInfo)
	err = service.GetUserTotalFavorited(&userInfo)
	err = service.GetUserWorkCount(&userInfo)
	err = service.GetUserFavoriteCount(&userInfo)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserResponse{
			Response:  Response{StatusCode: 6, StatusMsg: "user info construct error"}})
	}
	return c.Status(fiber.StatusOK).JSON(
		UserResponse{
			Response: Response{StatusCode: 0},
			User:     userInfo,
		},
	)
}
