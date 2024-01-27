package xvalidator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var CnMobileVerifyTag = "xv_cn_mobile"

func RegisterCnMobileValidation(validate *validator.Validate) error {
	return validate.RegisterValidation(CnMobileVerifyTag, VerifyCnMobile)
}

func VerifyCnMobile(f validator.FieldLevel) bool {
	val := f.Field().String()
	cnMobilePattern := `^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\d{8}$`
	reg, err := regexp.Compile(cnMobilePattern) // filter exclude chars
	if err != nil {
		return false
	}

	match := reg.MatchString(val)
	if !match {
		fmt.Println("not match error.")
		return false
	}

	return true
}
