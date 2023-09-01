package controller

import (
	"douyin/models"
	"douyin/service"
	"douyin/utils/jwt"
	"douyin/utils/validator"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RelationActionRequest struct {
	Token      string `query:"token" validate:"required"`
	ToUserID   string `query:"to_user_id" validate:"required"`
	ActionType string `query:"action_type" validate:"required,oneof=1 2"`
}

type UserListRequest struct {
	UserID string `query:"user_id" validate:"required"`
	Token  string `query:"token" validate:"required"`
}

type UserListResponse struct {
	Response
	UserList []models.UserInfo `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *fiber.Ctx) error {
	request := RelationActionRequest{}
	emptyResponse := Response{}
	if err := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return err
	}
	var fromId uint
	if err := jwt.JwtClient.AuthTokenValid(c, &emptyResponse, &fromId, request.Token); err != nil {
		return err
	}
	actionType, _ := strconv.Atoi(request.ActionType)
	toIdInt, _ := strconv.Atoi(request.ToUserID)
	toId := uint(toIdInt)
	if fromId == toId {
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "不能关注自己"})
	}
	// 关注/取关
	var err error
	if actionType == 1 {
		err = service.FollowAction(fromId, toId)
	} else {
		err = service.CancleAction(fromId, toId)
	}
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: err.Error()})
	} else {
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 0, StatusMsg: "关注/取消关注成功！"})
	}

}

func FollowList(c *fiber.Ctx) error {
	request := UserListRequest{}
	emptyResponse := UserListResponse{}
	if err := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return err
	}
	uidInt, _ := strconv.Atoi(request.UserID)
	uid := uint(uidInt)
	if err := jwt.JwtClient.AuthCurUser(c, &emptyResponse, request.Token, uid); err != nil {
		return err
	}

	tmpFollowList, err := service.FollowList(uid)
	fmt.Println(tmpFollowList)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "查询关注列表失败",
			},
			UserList: nil,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "查询关注列表成功",
			},
			UserList: tmpFollowList,
		})
	}
}

func FollowerList(c *fiber.Ctx) error {
	request := UserListRequest{}
	emptyResponse := UserListResponse{}
	if err := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return err
	}
	uidInt, _ := strconv.Atoi(request.UserID)
	uid := uint(uidInt)
	if err := jwt.JwtClient.AuthCurUser(c, &emptyResponse, request.Token, uid); err != nil {
		return err
	}

	tmpFollowerList, err := service.FollowerList(uid)
	fmt.Println(tmpFollowerList)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "查询粉丝列表失败",
			},
			UserList: nil,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "查询粉丝列表成功",
			},
			UserList: tmpFollowerList,
		})
	}
}

func FriendList(c *fiber.Ctx) error {
	request := UserListRequest{}
	emptyResponse := UserListResponse{}
	if err := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return err
	}
	uidInt, _ := strconv.Atoi(request.UserID)
	uid := uint(uidInt)
	if err := jwt.JwtClient.AuthCurUser(c, &emptyResponse, request.Token, uid); err != nil {
		return err
	}

	tmpFriendList, err := service.FriendList(uid)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "查询好友列表失败",
			},
			UserList: nil,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "查询好友列表成功",
			},
			UserList: tmpFriendList,
		})
	}
}
