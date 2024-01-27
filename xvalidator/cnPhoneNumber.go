package xvalidator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var CnPhoneNumberVerifyTag = "xv_cn_phone"

func RegisterCnPhoneNumberValidation(validate *validator.Validate) error {
	return validate.RegisterValidation(CnPhoneNumberVerifyTag, VerifyCnPhoneNumber)
}

func VerifyCnPhoneNumber(f validator.FieldLevel) bool {
	val := f.Field().String()
	cnPhoneNumberPattern := `^(0\d{2,3})-?(\d{7,8})$`
	reg, err := regexp.Compile(cnPhoneNumberPattern) // filter exclude chars
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
