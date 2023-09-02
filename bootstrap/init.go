package bootstrap

import (
	"douyin/config"
	"douyin/controller"
	"douyin/models"
	"douyin/service"
	"douyin/utils/jwt"
	"douyin/utils/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"douyin/middleware/rabbitmq"
)

func Init() (*fiber.App, error) {
	err := config.InitConfig()
	if err != nil {
		return nil, err
	}
	if err := models.InitDB(); err != nil {
		return nil, err
	}
	if err := models.InitRedis(); err != nil {
		return nil, err
	}
	if err := service.InitOSS(); err != nil {
		return nil, err
	}
	if err := service.Init2Redis(); err != nil {
		return nil, err
	}
	jwt.InitJWT()
	validator.InitValidator()
	err = rabbitmq.InitRabbitMQ()
	if err != nil {
		return nil, err
	}

	app := fiber.New()
	controller.RegisterRoutes(app)

	// Initialize default config
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		//	AllowHeaders: "Origin, Content-Type, Accept",
	}))

	return app, nil
}
