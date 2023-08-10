package controller

import (
	"douyin/models"
	"douyin/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CommentListResponse struct {
	Response
	CommentList []models.CommentInfo `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment models.CommentInfo `json:"comment,omitempty"`
}

func CommentAction(c *fiber.Ctx) error {
	token := c.Query("token")
	var uid uint
	if claimPtr, err := service.ParseToken(token); err != nil {
		fmt.Printf("User Unauthorized: %v\n", err)
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg:  "Unauthorized"})
	} else {
		uid = uint((*claimPtr).ID)
	}
	videoId := c.Query("video_id"); vid, _ := strconv.Atoi(videoId)
	actionType := c.Query("action_type")
	if _, err := service.GetUserById(uid); err != nil {
		fmt.Printf("user don't exist: %v\n", err)
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 2, StatusMsg: "User doesn't exist"})
	}
	if actionType == "1" {
		text := c.Query("comment_text")
		comment := models.Comment{
			UserId:  uid,
			VideoId: uint(vid),
			Content: text,
		}
		if err := service.CreateComment(&comment); err != nil {
			fmt.Printf("db create comment failed: %v\n", err)
			return c.Status(http.StatusOK).JSON(Response{StatusCode: 3, StatusMsg: "create comment failed"})
		}
		if commentInfo, err :=  service.GenerateCommentInfo(&comment); err != nil {
			fmt.Printf("get user info failed: %v\n", err)
			return c.Status(http.StatusOK).JSON(Response{StatusCode: 4, StatusMsg: "get userinfo failed"})	
		} else {
			return c.Status(http.StatusOK).JSON(CommentActionResponse{
				Response: Response{StatusCode: 0},
				Comment: commentInfo,
			})
		}
	} else {
		commentId := c.Query("comment_id"); cid, _ := strconv.Atoi(commentId)
		if err := service.DeleteComment(uint(cid)); err != nil {
			fmt.Printf("db delete comment failed: %v\n", err)
			return c.Status(http.StatusOK).JSON(Response{StatusCode: 5, StatusMsg: "delete comment failed"})
		}
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 0})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *fiber.Ctx) error {
	token := c.Query("token")
	if _, err := service.ParseToken(token); err != nil {
		fmt.Printf("User Unauthorized: %v\n", err)
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 1, StatusMsg:  "Unauthorized"})
	}
	videoId := c.Query("video_id"); vid, _ := strconv.Atoi(videoId)
	if commentInfos, err := service.GetCommentInfosByVideoId(uint(vid)); err != nil {
		fmt.Printf("get commentInfos failed: %v\n", err)
		return c.Status(http.StatusOK).JSON(Response{StatusCode: 2, StatusMsg: "get commentInfos failed"})
	} else {
		if len(commentInfos) == 0 {
			return c.Status(http.StatusOK).JSON(Response{StatusCode: 3, StatusMsg: "no comments found"})
		}
		return c.Status(http.StatusOK).JSON(CommentListResponse{
				Response: Response{StatusCode: 0},
				CommentList: commentInfos,
		})
	}
}
