package bootstrap

import (
	"douyin/config"
	"douyin/controller"
	"douyin/middleware"

	"douyin/models"
	"douyin/service"
	"douyin/utils/jwt"
	"douyin/utils/log"
	"douyin/utils/validator"

	"douyin/middleware/rabbitmq"
	"github.com/gofiber/fiber/v2"
)

func Init() (*fiber.App, error) {
	log.Init()
	err := config.InitConfig()
	if err != nil {
		return nil, err
	}
	if err = models.InitDB(); err != nil {
		return nil, err
	}
	if err = models.InitRedis(); err != nil {
		return nil, err
	}
	if err = service.InitOSS(); err != nil {
		return nil, err
	}
	if err = service.Init2Redis(); err != nil {
		return nil, err
	}
	jwt.InitJWT()
	validator.InitValidator()
	err = rabbitmq.InitRabbitMQ()
	if err != nil {
		return nil, err
	}

	app := fiber.New()

	// Initialize default config
	app.Use(middleware.Logging())

	controller.RegisterRoutes(app)

	return app, nil
}
