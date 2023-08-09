package public

import (
	"douyin/utils/jwt"
)

var (
	Jwt *jwt.JWT
)

func InitJWT() {
	// TODO: ADD jwt signing key to configuration file
	Jwt = jwt.NewJWT([]byte("test_key"))
}
