package utils

import "regexp"

var (
	emailReg = regexp.MustCompile("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$")
	// 至少6个字符
	weakPasswdReg = regexp.MustCompile("^[\\w_-]{6,16}$")
)

func IsEmail(email string) bool {
	return emailReg.MatchString(email)
}

func IsStrongPasswd(passwd string) bool {
	return weakPasswdReg.MatchString(passwd)
}
