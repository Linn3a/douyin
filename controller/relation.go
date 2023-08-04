package controller

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *fiber.Ctx) error {
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 0})
	} else {
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

// FriendList all users have same friend list
func FriendList(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

//// RelationAction no practical effect, just check if token is valid
//func RelationAction(c *gin.Context) {
//	token := c.Query("token")
//
//	if _, exist := usersLoginInfo[token]; exist {
//		c.JSON(http.StatusOK, Response{StatusCode: 0})
//	} else {
//		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
//	}
//}
//
//// FollowList all users have same follow list
//func FollowList(c *gin.Context) {
//	c.JSON(http.StatusOK, UserListResponse{
//		Response: Response{
//			StatusCode: 0,
//		},
//		UserList: []User{DemoUser},
//	})
//}
//
//// FollowerList all users have same follower list
//func FollowerList(c *gin.Context) {
//	c.JSON(http.StatusOK, UserListResponse{
//		Response: Response{
//			StatusCode: 0,
//		},
//		UserList: []User{DemoUser},
//	})
//}
//
//// FriendList all users have same friend list
//func FriendList(c *gin.Context) {
//	c.JSON(http.StatusOK, UserListResponse{
//		Response: Response{
//			StatusCode: 0,
//		},
//		UserList: []User{DemoUser},
//	})
//}
