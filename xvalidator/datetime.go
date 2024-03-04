package xvalidator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/orinchen/xlib/xtime"
)

var DatetimeVerifyTag = "xv_datetime"
var DatetimeErrorInfo = "{0}:日期格式不正确"

func RegisterDatetimeValidation(validate *validator.Validate, trans ut.Translator) error {
	err := validate.RegisterValidation(DatetimeVerifyTag, VerifyDatetime)
	if trans == nil || err != nil {
		return err
	}

	return validate.RegisterTranslation(DatetimeVerifyTag, trans, registerTranslator(DatetimeVerifyTag, DatetimeErrorInfo), translate)
}

func VerifyDatetime(f validator.FieldLevel) bool {
	val := f.Field().String()
	_, err := xtime.AutoParse(val)
	return err == nil
}
