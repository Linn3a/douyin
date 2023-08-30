package controller

import (
	"douyin/models"
	"douyin/service"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type FeedResponse struct {
	Response
	VideoList []models.VideoInfo `json:"video_list,omitempty"`
	NextTime  int64              `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *fiber.Ctx) error {
	// token := c.Query("token")
	// claims, _ := service.ParseToken(token)
	// fromId := uint(claims.ID)

	var DemoVideoList []models.Video
	var err error
	// DemoVideoInfo := models.NewVideoInfo(&DemoVideo)
	// DemoVideoInfo.Author = service.GenerateUserInfo(&DemoUser)
	// DemoVideoInfo.Author.IsFollow = service.HasRelation(fromId,DemoUser.ID)

	// DemoVideoList = append(DemoVideoList, DemoVideoInfo)
	DemoVideoList, err = service.GetVideosByUpdateAt()
	if err != nil {
		return c.Status(http.StatusOK).JSON(Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	var lenght int64
	lenght, err = service.GetVideosLenght()
	if err != nil {
		return c.Status(http.StatusOK).JSON(Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	videoList := make([]models.VideoInfo, lenght)
	for i, v := range DemoVideoList {
		videoList[i] = models.NewVideoInfo(&v)
	}
	return c.Status(http.StatusOK).JSON(FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})
}
