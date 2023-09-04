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
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return httpErr
	}
	var uid uint
	if err, httpErr := jwt.JwtClient.AuthTokenValid(c, &emptyResponse, &uid, request.Token); err != nil {
		return httpErr
	}
	vid, _ := strconv.Atoi(request.VideoID)
	actionType := request.ActionType
	if actionType == "1" {
		if err := service.AddFavoriteVideo(uid, uint(vid)); err != nil {
			fmt.Printf("add favorite failed: %v\n", err)
			return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 5, StatusMsg: "add favorite failed" + err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 0, StatusMsg: "add favorite success"})
	} else {
		if err := service.DeleteFavoriteVideo(uid, uint(vid)); err != nil {
			fmt.Printf("delete favorite failed: %v\n", err)
			return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 6, StatusMsg: "delete favorite failed" + err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 0, StatusMsg: "delete favorite success"})
	}
}

func FavoriteList(c *fiber.Ctx) error {
	request := FavoriteListRequest{}
	emptyResponse := VideoListResponse{}
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return httpErr
	}
	uidInt, _ := strconv.Atoi(request.UserID)
	uid := uint(uidInt)
	if err, httpErr := jwt.JwtClient.AuthTokenValid(c, &emptyResponse, &uid, request.Token); err != nil {
		return httpErr
	}

	vids, err := service.GetFavoriteVideoIds(uid)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 5, StatusMsg: "redis videos get error" + err.Error()}})
	}
	videoInfos, err := service.GetVideoInfosByIds(vids)
	if err != nil {
		fmt.Printf("video infos get error: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 6, StatusMsg: "mysql videos get error" + err.Error()}})
	}

	for i := 0; i < len(videoInfos); i++ {
		videoInfos[i].IsFavorite = true
	}
	
	return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: 
		Response{StatusCode: 0},
		VideoList: videoInfos,
	})
}
