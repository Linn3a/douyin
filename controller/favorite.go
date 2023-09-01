package controller

import (
	"douyin/service"
	"douyin/utils/jwt"
	"douyin/utils/validator"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type FavoriteListRequest struct {
	Token  string `query:"token" validate:"required"`   // 用户鉴权token
	UserID string `query:"user_id" validate:"required"` // 用户id
}

type FavoriteActionRequest struct {
	ActionType string `query:"action_type" validate:"required"` // 1-点赞，2-取消点赞
	Token      string `query:"token" validate:"required"`       // 用户鉴权token
	VideoID    string `query:"video_id" validate:"required"`    // 视频id
}

func FavoriteAction(c *fiber.Ctx) error {
	request := FavoriteActionRequest{}
	emptyResponse := Response{}
	if err := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return err
	}
	var uid uint
	if err := jwt.JwtClient.AuthTokenValid(c, &emptyResponse, &uid, request.Token); err != nil {
		return err
	}
	vid, _ := strconv.Atoi(request.VideoID)
	actionType := request.ActionType
	if actionType == "1" {
		if err := service.AddFavoriteVideo(uid, uint(vid)); err != nil {
			fmt.Printf("add favorite failed: %v\n", err)
			return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 5, StatusMsg: "add favorite failed" + err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 0})
	} else {
		if err := service.DeleteFavoriteVideo(uid, uint(vid)); err != nil {
			fmt.Printf("delete favorite failed: %v\n", err)
			return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 6, StatusMsg: "delete favorite failed" + err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 0})
	}
}

func FavoriteList(c *fiber.Ctx) error {
	request := c.Locals("request").(FavoriteListRequest)
	uid, _ := strconv.Atoi(request.UserID)
	if _, err := service.GetUserById(uint(uid)); err != nil {
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 4, StatusMsg: "user not exists" + err.Error()}})
	}

	vids, err := service.GetFavoriteVideoIds(uint(uid))
	videoInfos, err := service.GetVideoInfosByIds(vids)
	for i, _ := range videoInfos {
		videoInfos[i].IsFavorite = true
	}

	if err != nil {
		fmt.Printf("videos get error: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 5, StatusMsg: "videos get error" + err.Error()}})
	}

	return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: 
		Response{StatusCode: 0},
		VideoList: videoInfos,
	})
}
