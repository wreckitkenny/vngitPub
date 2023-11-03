package utils

import (
	"regexp"
	"unicode"
)

func UsernameValidator(username string) (bool) {
	logger := ConfigZap()
	matched, err := regexp.MatchString(`[a-z0-9(.)?]+\@vnpay.vn`, username)
	if err != nil {
		logger.Error(err)
	}

	if !matched {
		return matched
	}

	return matched
}

func PasswordValidator(password string) (string, bool) {
	const passwordMinLength = 12
	const passwordMaxLength = 30
	var (
        hasUpper   = false
        hasLower   = false
        hasNumber  = false
        hasSpecial = false
    )

	if len(password) < passwordMinLength {
		return "The password you entered is too short. Passwords must be 12-14 characters at least.", false
	}

	if len(password) > passwordMaxLength {
		return "The password you entered is too long. Max length of passwords is 30 characters.", false
	}

	for _, char := range password {
        switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsNumber(char):
				hasNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
        }
    }

	if hasUpper && hasLower && hasNumber && hasSpecial {
		return "", true
	}

	return "The password you entered should have a mix of upper case letters, lower case letters, numbers and special symbols.", false
}