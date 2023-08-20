package controller

import (
	"douyin/models"
	"douyin/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
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
		StatusMsg:  " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: []models.VideoInfo{},
	})
}
