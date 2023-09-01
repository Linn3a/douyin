package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)


type JWT struct {
	signingKey []byte
}

func New(signingKey []byte) *JWT {
	return &JWT{
		signingKey: signingKey,
	}
}

type BaseResponse interface {
	Set(sc int32, sm string)
}

// 1. 签发token，用于注册和登陆serv中
func (j *JWT)NewToken(uid uint) (string, error) {
	claims := jwt.MapClaims{
		"ID":  uid,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and return
	t, err := token.SignedString(j.signingKey)
	if err != nil {
		return "", err
	} else {
		return t, nil
	}
}

func parseTokenMap(token *jwt.Token) *jwt.MapClaims {
	claims := token.Claims.(jwt.MapClaims)
	return &claims
}

func (j *JWT)AuthTokenValid(c *fiber.Ctx, resp BaseResponse, uid *uint, requestToken string) error {
	token, err := jwt.ParseWithClaims(requestToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})
	if err != nil {
		fmt.Printf("token invalid: %v\n", err)
		resp.Set(3, "token invalid")
		return c.Status(fiber.StatusOK).JSON(resp)
	}
	(*uid) = (*parseTokenMap(token))["ID"].(uint)
	return nil
}

func (j *JWT)AuthCurUser(c *fiber.Ctx, resp BaseResponse, requestToken string, uid uint) error {
	token, err := jwt.ParseWithClaims(requestToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})
	if err != nil {
		fmt.Printf("token invalid: %v\n", err)
		resp.Set(3, "token invalid")
		return c.Status(fiber.StatusOK).JSON(resp)
	}
	user_id := (*parseTokenMap(token))["ID"].(uint)
	if user_id != uid {
		fmt.Printf("unauthorized: %v\n", err)
		resp.Set(4, "unauthorized")
		return c.Status(fiber.StatusOK).JSON(resp)
	}
	return nil
}
