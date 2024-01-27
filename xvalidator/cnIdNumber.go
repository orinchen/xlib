package xvalidator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var CnIdNumberVerifyTag = "xv_cn_id_number"

func RegisterCnIdNumberValidation(validate *validator.Validate) error {
	return validate.RegisterValidation(CnIdNumberVerifyTag, VerifyCnIdNumberNumber)
}

func VerifyCnIdNumberNumber(f validator.FieldLevel) bool {
	val := f.Field().String()

	var pattern string
	if len(val) == 18 {
		pattern = `^[1-9]\d{5}(19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$`
	} else if len(val) == 15 {
		pattern = `^[1-9]\d{5}\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{2}[0-9Xx]$`
	} else {
		return false
	}
	reg, err := regexp.Compile(pattern) // filter exclude chars
	if err != nil {
		return false
	}
	match := reg.MatchString(val)
	if !match {
		return false
	}
	return true
}
