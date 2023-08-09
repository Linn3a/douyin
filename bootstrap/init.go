package bootstrap

import (
	"douyin/config"
	"douyin/controller"
	"douyin/models"
	"douyin/public"
	"github.com/gofiber/fiber/v2"
)

func Init() (*fiber.App, error) {
	err := config.InitConfig()
	if err != nil {
		return nil, err
	}
	err = models.InitDB()
	if err != nil {
		return nil, err
	}
	public.InitJWT()

	app := fiber.New()
	controller.RegisterRoutes(app)

	return app, err
}
