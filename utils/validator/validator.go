package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)


type BaseResponse interface {
	Set(sc int32, sm string)
}

type customValidator struct {
	*validator.Validate
}

func New() customValidator {
	validate := validator.New()
	validate.RegisterValidation("required_if", RequiredIf)
	return customValidator{validate}
}

func (v customValidator)ValidateQuery(c *fiber.Ctx, resp BaseResponse, request interface{}) error {
	if err := c.QueryParser(request); err != nil {
		fmt.Printf("request type wrong: %v\n", err)
		resp.Set(1, "request type wrong")
		return c.Status(fiber.StatusOK).JSON(resp)
	}
	// 测试非空等限制
	if err := v.Struct(request); err != nil {
		fmt.Printf("request invalid: %v\n", err)
		resp.Set(2, "request invalid")
		return c.Status(fiber.StatusOK).JSON(resp)
	}
	return nil
}

func (v customValidator)ValidateStruct(c *fiber.Ctx, resp BaseResponse, request interface{}) error {
	if err := v.Struct(request); err != nil {
		fmt.Printf("request invalid: %v\n", err)
		resp.Set(2, "request invalid")
		return c.Status(fiber.StatusOK).JSON(resp)
	}
	return nil
}
