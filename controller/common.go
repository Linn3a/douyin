package controller

import (
	"github.com/go-playground/validator/v10"
	custom_validator "douyin/utils/validator"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

func ValidateStruct(q interface{}) error {
	validate := validator.New()
	validate.RegisterValidation("required_if", custom_validator.RequiredIf)
	if err := validate.Struct(q); err != nil {
		return err
	}
	return nil
}
