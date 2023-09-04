package controller

import (
	"douyin/models"
	"douyin/service"
	"douyin/utils/jwt"
	"douyin/utils/validator"
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

type FriendInfo struct {
	models.UserInfo
	Message string `json:"message"`
	MsgType int64  `json:"msgType"`
}

type FriendListResponse struct {
	Response
	UserList []FriendInfo `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *fiber.Ctx) error {
	request := RelationActionRequest{}
	emptyResponse := Response{}
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return httpErr
	}
	var fromId uint
	if err, httpErr := jwt.JwtClient.AuthTokenValid(c, &emptyResponse, &fromId, request.Token); err != nil {
		return httpErr
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
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return httpErr
	}
	uidInt, _ := strconv.Atoi(request.UserID)
	uid := uint(uidInt)
	if err, httpErr := jwt.JwtClient.AuthCurUser(c, &emptyResponse, request.Token, uid); err != nil {
		return httpErr
	}

	followingIds, err := service.GetFollowingIds(uid)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{Response: Response{StatusCode: 5, StatusMsg: "redis user get error: " + err.Error()}})
	}
	if len(followingIds) == 0 {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{Response: Response{StatusCode: 0,StatusMsg:  "暂无关注用户",},UserList: []models.UserInfo{}})
	}
	followingInfos, err := service.GetUserInfosByIds(followingIds)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{Response: Response{StatusCode: 6, StatusMsg: "userInfo get error: " + err.Error()}})
	}
	for i := 0; i < len(followingInfos); i++ {
		followingInfos[i].IsFollow = true
	}

	return c.Status(fiber.StatusOK).JSON(UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "查询关注列表成功",
		},
		UserList: followingInfos,
	})
}

func FollowerList(c *fiber.Ctx) error {
	request := UserListRequest{}
	emptyResponse := UserListResponse{}
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return httpErr
	}
	uidInt, _ := strconv.Atoi(request.UserID)
	uid := uint(uidInt)
	if err, httpErr := jwt.JwtClient.AuthCurUser(c, &emptyResponse, request.Token, uid); err != nil {
		return httpErr
	}

	followerIds, _ := service.GetFollowerIds(uid)
	if len(followerIds) == 0 {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{
			Response: Response{StatusCode: 0,StatusMsg:  "暂无粉丝"},UserList: []models.UserInfo{}})
	}
	followerInfos, err := service.GetUserInfosByIds(followerIds)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{Response: Response{StatusCode: 5,StatusMsg:  "查询粉丝列表失败"}})
	}
	for i := 0; i < len(followerIds); i++ {
		service.GetUserIsFollow(&followerInfos[i], uid)
	}
	return c.Status(fiber.StatusOK).JSON(UserListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "查询粉丝列表成功",
		},
		UserList: followerInfos,
	})
}

func FriendList(c *fiber.Ctx) error {
	request := UserListRequest{}
	emptyResponse := UserListResponse{}
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return httpErr
	}
	uidInt, _ := strconv.Atoi(request.UserID)
	uid := uint(uidInt)
	if err, httpErr := jwt.JwtClient.AuthCurUser(c, &emptyResponse, request.Token, uid); err != nil {
		return httpErr
	}

	friendIds, err := service.GetFriendIds(uid)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{Response: Response{StatusCode: 5, StatusMsg:  "redis 查询好友列表失败"}})
	}
	if len(friendIds) == 0 {
		c.Status(fiber.StatusOK).JSON(FriendListResponse{Response: Response{StatusCode: 0,StatusMsg:  "暂无好友"},UserList: []FriendInfo{}})
	}
	friendInfos, err := service.GetUserInfosByIds(friendIds)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(UserListResponse{Response: Response{StatusCode: 6, StatusMsg:  "查询好友列表失败",}})
	}
	friendList := make([]FriendInfo, len(friendIds))
	for i, friendInfo := range friendInfos {
		friendInfo.IsFollow = true
		friendList[i].UserInfo = friendInfo
		latestCommnetInfo, err := service.GetFriendLatestMessageInfo(uid, uint(friendInfo.ID))
		if err != nil {
			return c.Status(fiber.StatusOK).JSON(UserListResponse{Response: Response{StatusCode: 7, StatusMsg:  "查询好友最新消息失败",}})
		}
		friendList[i].Message = latestCommnetInfo.Content
		if latestCommnetInfo.FromUserID == int64(uid) {
			friendList[i].MsgType = 1
		} else {
			friendList[i].MsgType = 0
		}
	}

	return c.Status(fiber.StatusOK).JSON(FriendListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "查询好友列表成功",
		},
		UserList: friendList,
	})

}
