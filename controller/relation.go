package controller

import (
	"douyin/models"
	"douyin/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"fmt"
)

type UserListResponse struct {
	Response
	UserList []models.UserInfo `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *fiber.Ctx) error {
	token := c.Query("token")
	toIdInt, _ := strconv.Atoi(c.Query("to_user_id"))
	toId := uint(toIdInt)
	actionTypeInt, _ := strconv.Atoi(c.Query("action_type"))
	actionType := uint(actionTypeInt)
	if claims, err := service.ParseToken(token); err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}else{
		fromId := uint(claims.ID)
		if fromId == toId {
			return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "不能关注自己"})
		}

	// 关注/取关
		if actionType == 1{
			err = service.FollowAction(fromId, toId)
		}else{
			err = service.CancleAction(fromId, toId)
		}
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 1,  StatusMsg:err.Error()})
		} else {
			return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 0,StatusMsg:"关注/取消关注成功！"})
		}
	
	}

}

func FollowList(c *fiber.Ctx) error {
	token := c.Query("token")
	userIdInt, _ := strconv.Atoi(c.Query("user_id"))
	userId := uint(userIdInt)
	if _, err := service.ParseToken(token); err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	tmpFollowList, err := service.FollowList(userId)
	fmt.Println(tmpFollowList)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:"查询关注列表失败",
			},
			UserList: nil,
		})
	}else{
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:"查询关注列表成功",
			},
			UserList: tmpFollowList,
		})
	}
}


func FollowerList(c *fiber.Ctx) error {
	token := c.Query("token")
	userIdInt, _ := strconv.Atoi(c.Query("user_id"))
	userId := uint(userIdInt)
	if _, err := service.ParseToken(token); err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	tmpFollowerList, err := service.FollowerList(userId)
	fmt.Println(tmpFollowerList)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:"查询粉丝列表失败",
			},
			UserList: nil,
		})
	}else{
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:"查询粉丝列表成功",
			},
			UserList: tmpFollowerList,
		})
	}
}


func FriendList(c *fiber.Ctx) error {
	token := c.Query("token")
	userIdInt, _ := strconv.Atoi(c.Query("user_id"))
	userId := uint(userIdInt)
	if _, err := service.ParseToken(token); err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	tmpFriendList, err := service.FriendList(userId)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:"查询好友列表失败",
			},
			UserList: nil,
		})
	}else{
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:"查询好友列表成功",
			},
			UserList: tmpFriendList,
		})
	}
}