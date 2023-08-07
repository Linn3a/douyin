package main

import (
	// "douyin/service"
	"douyin/public"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	public.InitDatabase()
	public.InitJWT()
	// go service.RunMessageServer()
	app := fiber.New()
	app.Use(logger.New())
	initRouter(app)
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		Output: os.Stdout,
	}))
	app.Listen(":8080")
}
