package util

import "regexp"

var (
	emailRegexp = regexp.MustCompile("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$")
)

func CheckEmail(email string) bool {
	return emailRegexp.MatchString(email)
}
