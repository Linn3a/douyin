package controller

import (
	"douyin/models"
	"douyin/service"
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
	request := new(FavoriteActionRequest)
	if err := c.QueryParser(request); err != nil {
		fmt.Printf("request type wrong: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "request type wrong " + err.Error()})
	}
	if err := ValidateStruct(*request); err != nil {
		fmt.Printf("request invalid: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 2, StatusMsg: "request invalid " + err.Error()})
	}
	token := request.Token
	uid, err := service.GetUserID(token)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 3, StatusMsg: "get user id failed " + err.Error()})
	}
	vid, _ := strconv.Atoi(request.VideoID)
	// if _, err := service.GetUserById(uid); err != nil {
	// 	fmt.Printf("user don't exist: %v\n", err)
	// 	return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 2, StatusMsg: "User doesn't exist"})
	// }
	actionType := request.ActionType
	if actionType == "1" {
		// TODO: 已经点过赞的需要报错吗?
		if err := service.AddFavoriteVideo(uid, uint(vid)); err != nil {
			fmt.Printf("add favorite failed: %v\n", err)
			return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 4, StatusMsg: "add favorite failed" + err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 0})
	} else {
		// TODO: 先前没关注需要报错吗?
		if err := service.DeleteFavoriteVideo(uid, uint(vid)); err != nil {
			fmt.Printf("delete favorite failed: %v\n", err)
			return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 5, StatusMsg: "delete favorite failed" + err.Error()})
		}
		return c.Status(fiber.StatusOK).JSON(Response{StatusCode: 0})
	}
}

func FavoriteList(c *fiber.Ctx) error {
	request := new(FavoriteListRequest)
	if err := c.QueryParser(request); err != nil {
		fmt.Printf("request type wrong: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 1, StatusMsg: "request type wrong " + err.Error()}})
	}
	if err := ValidateStruct(*request); err != nil {
		fmt.Printf("request invalid: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 2, StatusMsg: "request invalid " + err.Error()}})
	}
	token := request.Token
	if _, err := service.ParseToken(token); err != nil {
		fmt.Printf("token invalid: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 3, StatusMsg: "token invalid" + err.Error()}})
	}

	uid, _ := strconv.Atoi(request.UserID)
	if _, err := service.GetUserById(uint(uid)); err != nil {
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 4, StatusMsg: "user not exists" + err.Error()}})
	}

	videos, err := service.GetFavoriteVideos(uint(uid))
	if err != nil {
		fmt.Printf("videos get error: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 5, StatusMsg: "videos get error" + err.Error()}})
	}

	authorIds := make([]uint, len(videos))
	videoIds := make([]uint, len(videos))
	for ind, video := range videos {
		authorIds[ind] = video.AuthorID
		videoIds[ind] = video.ID
	}

	userInfos, err := service.GetUserInfosByIds(authorIds)
	if err != nil {
		fmt.Printf("userInfos get error: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 6, StatusMsg: "videos get error" + err.Error()}})
	}
	favoriteCounts, err := service.CountFavoritedUsersByIds(videoIds)
	if err != nil {
		fmt.Printf("favorite counts get error: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 7, StatusMsg: "favorite counts get error" + err.Error()}})
	}
	commentCounts, err := service.CountCommentsByVideoIds(videoIds)
	if err != nil {
		fmt.Printf("comment counts get error: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(VideoListResponse{Response: Response{StatusCode: 8, StatusMsg: "comment counts get error" + err.Error()}})
	}

	videoInfos := make([]models.VideoInfo, len(videos))
	for ind, video := range videos {
		userInfo := userInfos[video.AuthorID]
		favoriteCount := favoriteCounts[video.ID]
		commentCount := commentCounts[video.ID]
		videoInfos[ind] = models.NewVideoInfo(&video, &userInfo, favoriteCount, commentCount)
	}
	return c.Status(fiber.StatusOK).JSON(VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoInfos,
	})
}
