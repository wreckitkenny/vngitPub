package utils

import (
	"errors"
	"regexp"
	"strings"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
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

func IsValidToken(bearer []string) (bool,string) {
	logger := ConfigZap()
	tokenString :=  strings.Split(bearer[0], " ")[1]
	if bearer != nil {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("nekottobtignv"), nil
		})

		if token.Valid {
			return true,""
		} else if errors.Is(err, jwt.ErrTokenMalformed) {
			logger.Error("Parsing user's token...Wrong Token Format")
			logger.Debugf("Parsing user's token...Wrong token format: %s", err)
			return false,"Wrong token format"
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			logger.Error("Parsing user's token...Signature Invalid")
			logger.Debugf("Parsing user's token...Signature Invalid: %s", err)
			return false,"Signature Invalid"
		} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
			logger.Error("Parsing user's token...Token Expired")
			logger.Debugf("Parsing user's token...Token Expired: %s", err)
			return false,"Token Expired"
		} else {
			logger.Error("Parsing user's token...Failed to handle token")
			logger.Debugf("Parsing user's token...Failed to handle token: %s", err)
			return false,"Failed to handle token"
		}
	}

	logger.Error("Parsing user's token...Token empty")
	return false,"Token empty"
}