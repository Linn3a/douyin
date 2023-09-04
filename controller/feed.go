package controller

import (
	"douyin/models"
	"douyin/service"
	"douyin/utils/jwt"
	"strconv"
	"time"

	// "net/http"
	// "time"

	"github.com/gofiber/fiber/v2"
)

type FeedResponse struct {
	Response
	VideoList []models.VideoInfo `json:"video_list"`
	NextTime  int64              `json:"next_time"`
}

// Feed same demo video list for every request
func Feed(c *fiber.Ctx) error {
	rawTimestamp := c.Query("latest_time")
	intTime, _ := strconv.ParseInt(rawTimestamp, 10, 64)
	var latestTime int64
	if intTime != 0 {
		latestTime = intTime
	} else {
		latestTime = time.Now().UnixMilli()
	}

	vids, err := service.GetFeedVideoIds(&latestTime)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(FeedResponse{Response: Response{StatusCode: 1, StatusMsg: "redis get videos error: " + err.Error()}})
	}
	if len(vids) == 0 {
		return c.Status(fiber.StatusOK).JSON(FeedResponse{Response: Response{StatusCode: 0, StatusMsg: "暂无发布视频"}, VideoList: []models.VideoInfo{}})
	}

	nextTime := latestTime
	videoInfos, err := service.GetVideoInfosByIds(vids)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(FeedResponse{Response: Response{StatusCode: 2, StatusMsg: "sql get videoinfos error: " + err.Error()}})
	}
	token := c.Query("token")
	var uid uint
	if err, _ = jwt.JwtClient.AuthTokenValid(c, &Response{}, &uid, token); err == nil {
		// 如果登陆 填充video favorite信息
		// 如果登陆 填充author follow信息
		for i := 0; i < len(videoInfos); i++ {
			err = service.GetVideoIsFavorite(&videoInfos[i], uid)
			if err != nil {
				return c.Status(fiber.StatusOK).JSON(FeedResponse{
					Response: Response{StatusCode: 2, StatusMsg: err.Error()},
				})
			}
			err = service.GetUserIsFollow(videoInfos[i].Author, uid)
			if err != nil {
				return c.Status(fiber.StatusOK).JSON(FeedResponse{
					Response: Response{StatusCode: 2, StatusMsg: err.Error()},
				})
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: err.Error()},
		VideoList: videoInfos,
		NextTime:  nextTime,
	})
}
