package main

import (
	"douyin/controller"
	"github.com/gofiber/fiber/v2"
)

func initRouter(app *fiber.App) {
	// public directory is used to serve static resources
	app.Static("/static", "./public")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./templates/index.html")
	})
	apiRouter := app.Group("/douyin")

	// basic apis
	apiRouter.Get("/feed/", controller.Feed)
	apiRouter.Get("/user/", controller.UserInfo)
	apiRouter.Post("/user/register/", controller.Register)
	apiRouter.Post("/user/login/", controller.Login)
	apiRouter.Post("/publish/action/", controller.Publish)
	apiRouter.Get("/publish/list/", controller.PublishList)

	// extra apis - I
	apiRouter.Post("/favorite/action/", controller.FavoriteAction)
	apiRouter.Get("/favorite/list/", controller.FavoriteList)
	apiRouter.Post("/comment/action/", controller.CommentAction)
	apiRouter.Get("/comment/list/", controller.CommentList)

	// extra apis - II
	apiRouter.Post("/relation/action/", controller.RelationAction)
	apiRouter.Get("/relation/follow/list/", controller.FollowList)
	apiRouter.Get("/relation/follower/list/", controller.FollowerList)
	apiRouter.Get("/relation/friend/list/", controller.FriendList)
	apiRouter.Get("/message/chat/", controller.MessageChat)
	apiRouter.Post("/message/action/", controller.MessageAction)
}
