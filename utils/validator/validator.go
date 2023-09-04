package validator

import (
	"douyin/utils/log"
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

func (v customValidator) ValidateQuery(c *fiber.Ctx, resp BaseResponse, request interface{}) (error, error) {
	if err := c.QueryParser(request); err != nil {
		log.FieldLog("validator", "error", fmt.Sprintf("query parse failed: %v", err))
		resp.Set(1, "request type wrong "+err.Error())
		return fmt.Errorf("request type wrong: %v", err), c.Status(fiber.StatusOK).JSON(resp)
	}
	// 测试非空等限制
	if err := v.Struct(request); err != nil {
		log.FieldLog("validator", "error", fmt.Sprintf("request invalid: %v", err))
		resp.Set(2, "request invalid "+err.Error())
		return fmt.Errorf("request invalid: %v", err), c.Status(fiber.StatusOK).JSON(resp)
	}
	return nil, nil
}

func (v customValidator) ValidateStruct(c *fiber.Ctx, resp BaseResponse, request interface{}) (error, error) {
	if err := v.Struct(request); err != nil {
		log.FieldLog("validator", "error", fmt.Sprintf("request invalid: %v", err))
		resp.Set(2, "request invalid "+err.Error())
		return fmt.Errorf("request invalid: %v", err), c.Status(fiber.StatusOK).JSON(resp)
	}
	return nil, nil
}
