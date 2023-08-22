package bootstrap

import (
	"douyin/config"
	"douyin/controller"
	"douyin/models"
	"douyin/public"
	"douyin/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"douyin/middleware/rabbitmq"
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
	err = service.InitOSS()
	if err != nil {
		return nil, err
	}
	err = rabbitmq.InitRabbitMQ()
	if err != nil {
		return nil, err
	}
	public.InitJWT()

	app := fiber.New()
	controller.RegisterRoutes(app)

	// Initialize default config
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		//	AllowHeaders: "Origin, Content-Type, Accept",
	}))

	return app, err
}
