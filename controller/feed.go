package controller

import (
	"douyin/models"
	"net/http"
	"time"
	"douyin/service"
	"github.com/gofiber/fiber/v2"
)

type FeedResponse struct {
	Response
	VideoList []models.VideoInfo `json:"video_list,omitempty"`
	NextTime  int64          `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *fiber.Ctx) error {
	// token := c.Query("token")
	// claims, _ := service.ParseToken(token)
	// fromId := uint(claims.ID)
	
    var DemoVideoList []models.VideoInfo
	DemoVideoInfo := models.NewVideoInfo(&DemoVideo)
	DemoVideoInfo.Author = service.GenerateUserInfo(&DemoUser)
	// DemoVideoInfo.Author.IsFollow = service.HasRelation(fromId,DemoUser.ID)
	DemoVideoList = append(DemoVideoList,DemoVideoInfo)
	return c.Status(http.StatusOK).JSON(FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: DemoVideoList,
		NextTime:  time.Now().Unix(),
	})

}
