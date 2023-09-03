package jwt

var JwtClient *JWT

func InitJWT() {
	JwtClient = New([]byte("test_key"))
}