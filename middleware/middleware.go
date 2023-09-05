package middleware

import (
	//"log"
	"douyin/utils/log"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func Logging() fiber.Handler {
	log.FieldLog("path", "info", "init middleware")
	return func(c *fiber.Ctx) error {
		// 创建中间件
		start := time.Now()

		defer func() {
			log.FieldLog("path", "info", fmt.Sprintf("from: %s, visit: %s, status: %v, spend: %v",
				c.IP(),
				c.Path(),
				c.Response().StatusCode(),
				time.Since(start)))
		}()

		return c.Next()
	}
}
