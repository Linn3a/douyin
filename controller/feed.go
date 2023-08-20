package controller

import (
	"douyin/service"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// type FeedResponse struct {
// 	Response
// 	VideoList []models.VideoInfo `json:"video_list,omitempty"`
// 	NextTime  int64              `json:"next_time,omitempty"`
// }

// =================================================================================================================

type FeedResponse struct {
	Response
	*service.FeedVideoList
}

func Feed(c *fiber.Ctx) error {
	p := NewProxyFeedVideoList(c)
	token := c.Query("token")
	//无登录状态
	if token == "" {
		err := p.DoNoToken()
		if err != nil {
			p.FeedVideoListError(err.Error())
		}
	}
	//有登录状态
	err := p.DoHasToken(token)
	if err != nil {
		p.FeedVideoListError(err.Error())
	}
	return nil
}

type ProxyFeedVideoList struct {
	*fiber.Ctx
}

func NewProxyFeedVideoList(c *fiber.Ctx) *ProxyFeedVideoList {
	return &ProxyFeedVideoList{Ctx: c}
}

// DoNoToken 未登录的视频流推送处理
func (p *ProxyFeedVideoList) DoNoToken() error {
	rawTimestamp := p.Query("latest_time")
	var latestTime time.Time
	intTime, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err == nil {
		latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
	}
	videoList, err := service.QueryFeedVideoList(0, latestTime)
	if err != nil {
		return err
	}
	p.FeedVideoListOk(videoList)
	return nil
}

// DoHasToken 如果是登录状态，则生成UserId字段
func (p *ProxyFeedVideoList) DoHasToken(token string) error {
	//解析成功
	claimPtr, err := service.ParseToken(token)
	if err != nil {
		fmt.Printf("token invalid: %v\n", err)
		return p.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "token invalid"})
	}
	// //token超时
	// if time.Now().Unix() > claimPtr.ExpiresAt {
	// 	return errors.New("token超时")
	// }
	rawTimestamp := p.Query("latest_time")
	var latestTime time.Time
	intTime, err := strconv.ParseInt(rawTimestamp, 10, 64)
	if err != nil {
		latestTime = time.Unix(0, intTime*1e6) //注意：前端传来的时间戳是以ms为单位的
	}
	//调用service层接口
	uid := uint((*claimPtr).ID)
	videoList, err := service.QueryFeedVideoList(uid, latestTime)
	if err != nil {
		return err
	}
	p.FeedVideoListOk(videoList)
	return nil
}

func (p *ProxyFeedVideoList) FeedVideoListError(msg string) {
	p.Status(http.StatusOK).JSON(FeedResponse{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  msg,
		}})
}

func (p *ProxyFeedVideoList) FeedVideoListOk(videoList *service.FeedVideoList) {
	p.Status(http.StatusOK).JSON(FeedResponse{
		Response: Response{
			StatusCode: 0,
		},
		FeedVideoList: videoList,
	},
	)
}

// ==================================================================================================

// Feed same demo video list for every request
// func Feed(c *fiber.Ctx) error {
// 	// token := c.Query("token")
// 	// claims, _ := service.ParseToken(token)
// 	// fromId := uint(claims.ID)

// 	var DemoVideoList []models.VideoInfo
// 	DemoVideoInfo := models.NewVideoInfo(&DemoVideo, &models.UserInfo{}, 0, 0)
// 	DemoVideoInfo.Author = service.GenerateUserInfo(&DemoUser)
// 	// DemoVideoInfo.Author.IsFollow = service.HasRelation(fromId,DemoUser.ID)
// 	DemoVideoList = append(DemoVideoList, DemoVideoInfo)
// 	return c.Status(http.StatusOK).JSON(FeedResponse{
// 		Response:  Response{StatusCode: 0},
// 		VideoList: DemoVideoList,
// 		NextTime:  time.Now().Unix(),
// 	})
// }

//  ==================================================================================================
