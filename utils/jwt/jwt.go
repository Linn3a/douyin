package jwt

import (
	"fmt"

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

type CustomClaim struct {
	ID uint
	jwt.StandardClaims
}

// 1. 签发token，用于注册和登陆serv中
func (j *JWT) NewToken(uid uint) (string, error) {
	claims := CustomClaim{
		ID: uid,
	}

	// Generate encoded token and return
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.signingKey)
	if err != nil {
		return "", err
	} else {
		return token, nil
	}
}

func parseTokenMap(token *jwt.Token) *CustomClaim {
	claims := token.Claims.(*CustomClaim)
	return claims
}

func (j *JWT) AuthTokenValid(c *fiber.Ctx, resp BaseResponse, uid *uint, requestToken string) (error, error) {
	token, err := jwt.ParseWithClaims(requestToken, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})
	if err != nil {
		fmt.Printf("token invalid: %v\n", err)
		resp.Set(3, "token invalid")
		return fmt.Errorf("token invalid: %v", err), c.Status(fiber.StatusOK).JSON(resp)
	}
	(*uid) = (*parseTokenMap(token)).ID
	return nil, nil
}

func (j *JWT) AuthCurUser(c *fiber.Ctx, resp BaseResponse, requestToken string, uid uint) (error, error) {
	token, err := jwt.ParseWithClaims(requestToken, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})
	if err != nil {
		fmt.Printf("token invalid: %v\n", err)
		resp.Set(3, "token invalid")
		return fmt.Errorf("token invalid: %v", err), c.Status(fiber.StatusOK).JSON(resp)
	}
	user_id := (*parseTokenMap(token)).ID
	fmt.Print(user_id, uid)
	if user_id != uid {
		fmt.Print(user_id, uid)
		fmt.Printf("unauthorized: %v\n", err)
		resp.Set(4, "unauthorized")
		return fmt.Errorf("unauthorized: %v", err), c.Status(fiber.StatusOK).JSON(resp)
	}
	return nil, nil
}
