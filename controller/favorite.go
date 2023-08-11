package controller

import (
	"douyin/models"
	"douyin/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FavoriteAction(c *fiber.Ctx) error {
	token := c.Query("token")
	claimPtr, err := service.ParseToken(token)
	if err != nil {
		fmt.Printf("User Unauthorized: %v\n", err)
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "Unauthorized"})
	}
	uid := uint((*claimPtr).ID)
	videoId := c.Query("video_id")
	vid, _ := strconv.Atoi(videoId)
	actionType := c.Query("action_type")
	if _, err := service.GetUserById(uid); err != nil {
		fmt.Printf("user don't exist: %v\n", err)
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 2, StatusMsg: "User doesn't exist"})
	}
	if actionType == "1" {
		if err := service.AddFavoriteVideo(uid, uint(vid)); err != nil {
			return c.Status(http.StatusOK).JSON(Response{StatusCode: 3, StatusMsg: "add favorite failed"})
		}
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 0})
	} else {
		if err := service.DeleteFavoriteVideo(uid, uint(vid)); err != nil {
			return c.Status(http.StatusOK).JSON(Response{StatusCode: 4, StatusMsg: "delete favorite failed"})
		}
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 0})
	}
}

func FavoriteList(c *fiber.Ctx) error {
	token := c.Query("token")
	claimPtr, err := service.ParseToken(token)
	if err != nil {
		fmt.Printf("User Unauthorized: %v\n", err)
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "Unauthorized"})
	}
	uid := uint((*claimPtr).ID)

	videos, err := service.GetFavoriteVideos(uid)
	if err != nil {
		fmt.Printf("videos get error: %v\n", err)
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 2, StatusMsg: "videos get error"})
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
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 3, StatusMsg: "videos get error"})
	}
	favoriteCounts, err := service.CountFavoritedUsersByIds(videoIds)
	if err != nil {
		fmt.Printf("favorite counts get error: %v\n", err)
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 4, StatusMsg: "favorite counts get error"})
	}
	commentCounts, err := service.CountCommentsByVideoIds(videoIds)
	if err != nil {
		fmt.Printf("comment counts get error: %v\n", err)
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 5, StatusMsg: "comment counts get error"})
	}

	videoInfos := make([]models.VideoInfo, len(videos))
	for ind, video := range videos {
		userInfo := userInfos[video.AuthorID]
		favoriteCount := favoriteCounts[video.ID]
		commentCount := commentCounts[video.ID]
		videoInfos[ind] = models.NewVideoInfo(&video, &userInfo, favoriteCount, commentCount)
	}
	return c.Status(http.StatusOK).JSON(VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
