package controller

import (
	"log"
	"os"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	log.Println("init router")
	// public directory is used to serve static resources
	app.Static("/static", "./public")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./templates/index.html")
	})
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		Output: os.Stdout,
	}))
	apiRouter := app.Group("/douyin")

	// basic apis
	apiRouter.Get("/feed/", Feed)
	apiRouter.Get("/user/", UserInfo)
	apiRouter.Post("/user/register/", Register)
	apiRouter.Post("/user/login/", Login)
	apiRouter.Post("/publish/action/", Publish)
	apiRouter.Get("/publish/list/", PublishList)

	// extra apis - I
	apiRouter.Post("/favorite/action/", FavoriteAction)
	apiRouter.Get("/favorite/list/", FavoriteList)
	apiRouter.Post("/comment/action/", CommentAction)
	apiRouter.Get("/comment/list/", CommentList)

	// extra apis - II
	apiRouter.Post("/relation/action/", RelationAction)
	apiRouter.Get("/relation/follow/list/", FollowList)
	apiRouter.Get("/relation/follower/list/", FollowerList)
	apiRouter.Get("/relation/friend/list/", FriendList)
	apiRouter.Get("/message/chat/", MessageChat)
	apiRouter.Post("/message/action/", MessageAction)
}
