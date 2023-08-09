package jwt

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// JWT signing Key
type JWT struct {
	SigningKey []byte
}

// private claims, share information between parties that agree on using them
type CustomClaims struct {
	ID int64
	// TODO: add feature for different authority of different user
	// AuthorityId int64
	jwt.StandardClaims
}

func NewJWT(SigningKey []byte) *JWT {
	return &JWT{
		SigningKey,
	}
}

// func (j *JWT) CreateClaim(audience string, issuer string, validTime int64, id uint) (jwt.StandardClaims) {
// 	return jwt.StandardClaims{
// 		Audience: audience,
// 		ExpiresAt: time.Now().Unix() + validTime,
// 		Id: strconv.FormatInt(id64, 10),
// 		IssuedAt: time.Now().Unix(),
// 		Issuer: "douyin",
// 		NotBefore: time.Now().Unix(),
// 		Subject: "token",
// 	}
// }

// 利用siningkey加密包裹customclaims返回加密字符串
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// zap.S().Debugf(token.SigningString())
	return token.SignedString(j.SigningKey)
}

// 利用singingkey 对token有效性进行验证 有效则返回CustomClaims
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("token malformed")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, fmt.Errorf("token not validate yet")
			} else {
				return nil, fmt.Errorf("token not valid")
			}

		}
	}
	// verify the token claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token not valid")
}
