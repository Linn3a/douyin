package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var RequiredIf validator.Func = func(fl validator.FieldLevel) bool {

	/* 用法示例
		Type        string	`json:"type" validate:"required,oneof=flat_off percent_off"`
		MaxValue	uint	`json:"max_value" validate:"required_if=Type percent_off"`
	*/

	otherFieldName := strings.Split(fl.Param(), " ")[0]
	otherFieldValCheck := strings.Split(fl.Param(), " ")[1]

	var otherFieldVal reflect.Value
	if fl.Parent().Kind() == reflect.Ptr {
		otherFieldVal = fl.Parent().Elem().FieldByName(otherFieldName)
	} else {
		otherFieldVal = fl.Parent().FieldByName(otherFieldName)
	}
	fmt.Printf("other filed name=%s, other filed value=%v\n value check=%v\n", otherFieldName, otherFieldVal, otherFieldValCheck)

	if otherFieldValCheck == otherFieldVal.String() {
		return !isNilOrZeroValue(fl.Field())
	}
	return true
}

func isNilOrZeroValue(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsZero()
	}
}
