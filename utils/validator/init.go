package validator

var ValidateClient customValidator

func InitValidator() {
	ValidateClient = New()
}