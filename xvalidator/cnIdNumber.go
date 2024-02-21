package xvalidator

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var CnIdNumberVerifyTag = "xv_cn_id_number"

func RegisterCnIdNumberValidation(validate *validator.Validate, trans ut.Translator) error {
	err := validate.RegisterValidation(CnIdNumberVerifyTag, VerifyCnIdNumberNumber)
	if err != nil || trans == nil {
		return err
	}

	return validate.RegisterTranslation(CnIdNumberVerifyTag, trans, registerTranslator(CnIdNumberVerifyTag, "{0}身份证号码格式不正确"), translate)
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
