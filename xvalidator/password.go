package xvalidator

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var ComplexPwdVerifyTag = "xv_complex_pwd"
var ComplexPwdFieldErrorInfo = "密码应包含数字、大/小写字母、特殊字符中的3种, 且至少8个字符"

func RegisterComplexPwdValidation(validate *validator.Validate) error {
	return validate.RegisterValidation(ComplexPwdVerifyTag, VerifyPwd)
}

func VerifyPwd(f validator.FieldLevel) bool {
	val := f.Field().String()
	if len(val) < 8 || len(val) > 20 { // length需要通过验证
		fmt.Println("pwd length error")
		return false
	}

	pwdPattern := `^[0-9a-zA-Z!@#$%^&*~-_+]{8,20}$`
	reg, err := regexp.Compile(pwdPattern) // filter exclude chars
	if err != nil {
		return false
	}

	match := reg.MatchString(val)
	if !match {
		fmt.Println("not match error.")
		return false
	}

	var cnt int = 0          // 满足3中以上即可通过验证
	patternList := []string{ // 数字、大小写字母、特殊字符
		`[0-9]+`,
		`[a-z]+`,
		`[A-Z]+`,
		`[!@#$%^&*~-_+]+`,
	}
	for _, pattern := range patternList {
		match, _ = regexp.MatchString(pattern, val)
		if match {
			cnt++
		}
	}
	if cnt < 3 {
		fmt.Println("pwd should include at least 3 types.")
		return false
	}
	return true
}
