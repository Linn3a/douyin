package controller

import (
	"douyin/models"
	"douyin/service"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CommentListRequest struct {
	Token   string `query:"token" validate:"required"`    // 用户鉴权token
	VideoID string `query:"video_id" validate:"required"` // 视频id
}

type CommentActionRequest struct {
	ActionType  string `query:"action_type" validate:"required,oneof=1 2"`        // 1-发布评论，2-删除评论
	CommentID   string `query:"comment_id" validate:"required_if=ActionType 2"`   // 要删除的评论id，在action_type=2的时候使用
	CommentText string `query:"comment_text" validate:"required_if=ActionType 1"` // 用户填写的评论内容，在action_type=1的时候使用
	Token       string `query:"token" validate:"required"`                        // 用户鉴权token
	VideoID     string `query:"video_id" validate:"required"`                     // 视频id
}

type CommentListResponse struct {
	Response
	CommentList []models.CommentInfo `json:"comment_list"`
}

type CommentActionResponse struct {
	Response
	Comment *models.CommentInfo `json:"comment"`
}

func CommentAction(c *fiber.Ctx) error {
	request := new(CommentActionRequest)
	if err := c.QueryParser(request); err != nil {
		fmt.Printf("request type wrong: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(CommentActionResponse{Response: Response{StatusCode: 1, StatusMsg: "request type wrong " + err.Error()}})
	}
	if err := ValidateStruct(*request); err != nil {
		fmt.Printf("request invalid: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(CommentActionResponse{Response: Response{StatusCode: 2, StatusMsg: "request invalid " + err.Error()}})
	}
	token := request.Token
	claimPtr, err := service.ParseToken(token)
	if err != nil {
		fmt.Printf("token invalid: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(CommentActionResponse{Response: Response{StatusCode: 3, StatusMsg: "token invalid"}}) //不加null写法
	}
	uid := uint((*claimPtr).ID)
	vid, _ := strconv.Atoi(request.VideoID)
	// TODO: vid, uid不存在是否报错
	// if _, err := service.GetUserById(uid); err != nil {
	// 	fmt.Printf("user don't exist: %v\n", err)
	// 	return c.Status(fiber.StatusOK).JSON(CommentActionResponse{Response:Response{StatusCode: 2, StatusMsg: "User doesn't exist"}})
	// }
	actionType := request.ActionType
	if actionType == "1" {
		text := request.CommentText
		comment := models.Comment{
			UserId:  uid,
			VideoId: uint(vid),
			Content: text,
		}
		if err := service.CreateComment(&comment); err != nil {
			fmt.Printf("db create comment failed: %v\n", err)
			return c.Status(fiber.StatusOK).JSON(CommentActionResponse{Response: Response{StatusCode: 4, StatusMsg: "create comment failed"}})
		}
		commentInfo, err := service.GenerateCommentInfo(&comment)
		if err != nil {
			fmt.Printf("get user info failed: %v\n", err)
			return c.Status(fiber.StatusOK).JSON(CommentActionResponse{Response: Response{StatusCode: 5, StatusMsg: "get userinfo failed"}})
		}
		return c.Status(fiber.StatusOK).JSON(CommentActionResponse{
			Response: Response{StatusCode: 0},
			Comment:  &commentInfo,
		})
	} else {
		commentId := request.CommentID
		cid, _ := strconv.Atoi(commentId)
		if err := service.DeleteComment(uint(cid)); err != nil {
			fmt.Printf("db delete comment failed: %v\n", err)
			return c.Status(fiber.StatusOK).JSON(CommentActionResponse{Response: Response{StatusCode: 6, StatusMsg: "delete comment failed"}})
		}
		return c.Status(fiber.StatusOK).JSON(CommentActionResponse{Response: Response{StatusCode: 0}})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *fiber.Ctx) error {
	request := new(CommentListRequest)
	if err := c.QueryParser(request); err != nil {
		fmt.Printf("request type wrong: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(CommentListResponse{Response: Response{StatusCode: 1, StatusMsg: "request type wrong " + err.Error()}})
	}
	if err := ValidateStruct(*request); err != nil {
		fmt.Printf("request invalid: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(CommentListResponse{Response: Response{StatusCode: 2, StatusMsg: "request invalid " + err.Error()}})
	}
	token := request.Token
	if _, err := service.ParseToken(token); err != nil {
		fmt.Printf("token invalid: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(CommentListResponse{Response: Response{StatusCode: 3, StatusMsg: "token invalid"}})
	}
	vid, _ := strconv.Atoi(request.VideoID)
	// TODO: vid不存在是否报错
	commentInfos, err := service.GetCommentInfosByVideoId(uint(vid))
	if err != nil {
		fmt.Printf("get commentInfos failed: %v\n", err)
		return c.Status(fiber.StatusOK).JSON(CommentListResponse{Response: Response{StatusCode: 4, StatusMsg: "get commentInfos failed"}})
	}
	if len(commentInfos) == 0 {
		return c.Status(fiber.StatusOK).JSON(CommentListResponse{Response: Response{StatusCode: 0, StatusMsg: "no comments found"}, CommentList: []models.CommentInfo{}})
	}
	return c.Status(fiber.StatusOK).JSON(CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentInfos,
	})
}
