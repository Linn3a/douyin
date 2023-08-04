package controller

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

// FavoriteAction no practical effect, just check if token is valid

func FavoriteAction(c *fiber.Ctx) error {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 0})
	} else {
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})

}

//func FavoriteAction(c *gin.Context) {
//	token := c.Query("token")
//
//	if _, exist := usersLoginInfo[token]; exist {
//		c.JSON(http.StatusOK, Response{StatusCode: 0})
//	} else {
//		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
//	}
//}
//
//// FavoriteList all users have same favorite video list
//func FavoriteList(c *gin.Context) {
//	c.JSON(http.StatusOK, VideoListResponse{
//		Response: Response{
//			StatusCode: 0,
//		},
//		VideoList: DemoVideos,
//	})
//}
