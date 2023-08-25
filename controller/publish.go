package controller

import (
	"douyin/models"
	"douyin/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type VideoListResponse struct {
	Response
	VideoList []models.VideoInfo `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *fiber.Ctx) error {
	token := c.FormValue("token", "0")
	userID, err := service.GetUserID(token)
	if err != nil {
		log.Printf("Get user id error:%v", err)
		return c.Status(http.StatusOK).JSON(Response{
			StatusCode: 1,
			StatusMsg:  "Invalid token",
		})
	}

	data, err := c.FormFile("data")
	if err != nil {
		return c.Status(http.StatusOK).JSON(Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})

	}
	videoUrl, coverUrl, err := service.UploadVideoToOSS(data)
	if err != nil {
		return c.Status(http.StatusOK).JSON(Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	fmt.Printf("videoUrl:%v\n", videoUrl)
	fmt.Printf("coverUrl:%v\n", coverUrl)

	//println(data.Filename)
	//filename := filepath.Base(data.Filename)

	newVideo := models.Video{
		Title:    c.FormValue("title", "title"),
		PlayUrl:  videoUrl,
		CoverUrl: coverUrl,
		AuthorID: userID,
	}

	err = service.CreateVideo(newVideo)
	if err != nil {
		log.Printf("Mysql create video error:%v", err)
		return c.Status(http.StatusOK).JSON(Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(Response{
		StatusCode: 0,
		StatusMsg:  "upload successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *fiber.Ctx) error {
	token := c.FormValue("token", "0")
	userID, err := service.GetUserID(token)
	if err != nil {
		log.Printf("Get user id error:%v", err)
		return c.Status(http.StatusOK).JSON(Response{
			StatusCode: 1,
			StatusMsg:  "Invalid token",
		})
	}
	videos, err := service.GetVideosByUserId(userID)
	if err != nil {
		return err
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
	return c.Status(http.StatusOK).JSON(VideoListResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "upload successfully",
		},
		VideoList: videoInfos,
	})
}
