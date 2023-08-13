package controller

import (
	"douyin/models"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

type VideoListResponse struct {
	Response
	VideoList []models.VideoInfo `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *fiber.Ctx) error {
	var token string
	if err := c.BodyParser(&token); err != nil {
		return err
	}
	if _, exist := usersLoginInfo[token]; !exist {
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

	data, err := c.FormFile("data")
	if err != nil {
		return c.Status(http.StatusOK).JSON(Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})

	}

	filename := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.ID, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveFile(data, saveFile); err != nil {
		return c.Status(http.StatusOK).JSON(Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})

	}

	return c.Status(http.StatusOK).JSON(Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
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

//// Publish check token then save upload file to public directory
//func Publish(c *gin.Context) {
//	token := c.PostForm("token")
//
//	if _, exist := usersLoginInfo[token]; !exist {
//		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
//		return
//	}
//
//	data, err := c.FormFile("data")
//	if err != nil {
//		c.JSON(http.StatusOK, Response{
//			StatusCode: 1,
//			StatusMsg:  err.Error(),
//		})
//		return
//	}
//
//	filename := filepath.Base(data.Filename)
//	user := usersLoginInfo[token]
//	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
//	saveFile := filepath.Join("./public/", finalName)
//	if err := c.SaveUploadedFile(data, saveFile); err != nil {
//		c.JSON(http.StatusOK, Response{
//			StatusCode: 1,
//			StatusMsg:  err.Error(),
//		})
//		return
//	}
//
//	c.JSON(http.StatusOK, Response{
//		StatusCode: 0,
//		StatusMsg:  finalName + " uploaded successfully",
//	})
//}

//// PublishList all users have same publish video list
//func PublishList(c *gin.Context) {
//	c.JSON(http.StatusOK, VideoListResponse{
//		Response: Response{
//			StatusCode: 0,
//		},
//		VideoList: DemoVideos,
//	})
//}
