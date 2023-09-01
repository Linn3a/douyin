package controller

import (
	"douyin/models"
	"douyin/service"
	"douyin/utils/jwt"
	"douyin/utils/validator"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// type PublishActionRequest struct {
// 	Token string `query:"token" validate:"required"`
// 	Data  []byte `query:"data" validate:"required"`
// 	Title string `query:"title" validate:"required"`
// }

type PublishListRequest struct {
	UserID string `query:"user_id" validate:"required"`
	Token  string `query:"token" validate:"required"`
}

type VideoListResponse struct {
	Response
	VideoList []models.VideoInfo `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *fiber.Ctx) error {
	token := c.FormValue("token")
	var uid uint
	if err := jwt.JwtClient.AuthTokenValid(c, &Response{}, &uid, token); err != nil {
		return err
	}
	title := c.FormValue("title")
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

	if err := service.CreateVideo(title, videoUrl, coverUrl, uid); err != nil {
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
	request := PublishListRequest{}
	emptyResponse := VideoListResponse{}
	if err := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return err
	}
	uid, _ := strconv.Atoi(request.UserID)
	if err := jwt.JwtClient.AuthCurUser(c, &emptyResponse, request.Token, uint(uid)); err != nil {
		return err
	}

	vids, err := service.GetVideoIdsByUserId(uint(uid))
	if err != nil {
		return c.Status(http.StatusOK).JSON(VideoListResponse{
			Response: Response{
				StatusCode: 5,
				StatusMsg: err.Error(),
			},
		})
	}

	videoInfos, err := service.GetVideoInfosByIds(vids) 
	if err != nil {
		return c.Status(http.StatusOK).JSON(VideoListResponse{
			Response: Response{
				StatusCode: 6,
				StatusMsg: err.Error(),
			},
		})
	}

	return c.Status(http.StatusOK).JSON(VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoInfos,
	})
}
