package utils

import "regexp"

func CheckPWD(pwd string) bool {
	if len(pwd) < 8 || len(pwd) > 16 {
		return false
	}

	checkLetter := regexp.MustCompile(`[a-zA-Z]`)
	checkNumber := regexp.MustCompilePOSIX(`[0-9]`)
	checkSpecialChar := regexp.MustCompile(`[~!@$%^&*.]`)

	letterbool := checkLetter.MatchString(pwd)
	numberbool := checkNumber.MatchString(pwd)
	specialcharbool := checkSpecialChar.MatchString(pwd)

	return (letterbool && numberbool) || (letterbool && specialcharbool) || (numberbool && specialcharbool)
}
