package controller

import (
	"douyin/models"
	"douyin/service"
	"douyin/utils/jwt"
	"douyin/utils/validator"
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
	request := CommentActionRequest{}
	emptyResponse := CommentActionResponse{}
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return httpErr
	}
	var uid uint
	if err, httpErr := jwt.JwtClient.AuthTokenValid(c, &emptyResponse, &uid, request.Token); err != nil {
		return httpErr
	}
	vid, _ := strconv.Atoi(request.VideoID)
	actionType := request.ActionType
	if actionType == "1" {
		text := request.CommentText
		comment, err := service.CreateComment(uid, uint(vid), text)
		if err != nil {
			fmt.Printf("db create comment failed: %v\n", err)
			return c.Status(fiber.StatusOK).JSON(CommentActionResponse{Response: Response{StatusCode: 4, StatusMsg: "create comment failed"}})
		}
		commentInfo := service.GenerateCommentInfo(comment)
		userInfo, err := service.GetUserInfoById(comment.UserId)
		if err != nil {
			fmt.Printf("get user info failed: %v\n", err)
			return c.Status(fiber.StatusOK).JSON(CommentActionResponse{Response: Response{StatusCode: 5, StatusMsg: "get userinfo failed"}})
		}
		commentInfo.User = &userInfo
		// 填充is follow 信息
		service.GetUserIsFollow(commentInfo.User, uid)
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

func CommentList(c *fiber.Ctx) error {
	request := CommentListRequest{}
	emptyResponse := CommentListResponse{}
	if err, httpErr := validator.ValidateClient.ValidateQuery(c, &emptyResponse, &request); err != nil {
		return httpErr
	}
	var uid uint
	if err, httpErr := jwt.JwtClient.AuthTokenValid(c, &emptyResponse, &uid, request.Token); err != nil {
		return httpErr
	}
	vid, _ := strconv.Atoi(request.VideoID)
	cids, err := service.GetCommentIdsByVideoId(uint(vid))
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(CommentListResponse{Response: Response{StatusCode: 5, StatusMsg: "redis get comments error"}})
	}
	if len(cids) == 0 {
		return c.Status(fiber.StatusOK).JSON(CommentListResponse{Response: Response{StatusCode: 0, StatusMsg: "暂无评论信息"}, CommentList: []models.CommentInfo{}})
	}
	comments, err := service.GetCommentsByIds(cids)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(CommentListResponse{Response: Response{StatusCode: 6, StatusMsg: "sql get comments error"}})
	}
	uids := make([]uint, len(cids))
	for i, c := range comments {
		uids[i] = c.UserId
	}
	userInfos, err := service.GetUserInfoMapByIds(uids)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(CommentListResponse{Response: Response{StatusCode: 7, StatusMsg: "get user infos error"}})
	}
	commentInfos := make([]models.CommentInfo, len(cids))
	for i, c := range comments {
		commentInfos[i] = service.GenerateCommentInfo(&c)
		userInfo := userInfos[c.UserId]
		commentInfos[i].User = &userInfo
		// 填充is follow信息
		service.GetUserIsFollow(commentInfos[i].User, uid)
	}
	return c.Status(fiber.StatusOK).JSON(CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentInfos,
	})
}
