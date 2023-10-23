package controller

import (
	"regexp"

	"github.com/wreckitkenny/vngitpub/pkg/utils"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `"json:username"`
	Password string `"json:password"`
}

func Login(c *gin.Context) (string, string) {
	logger := utils.ConfigZap()
	var user User
	if err := c.BindJSON(&user); err != nil {
		logger.Error(err)
	}

	username := user.Username
	// password := user.Password

	matched := usernameValidator(username)
	if !matched {
		logger.Warn("Username must be ended with @vnpay.vn")
		return "", "Username must be ended with @vnpay.vn"
	}

	// if err := passwordValidator(password); err != nil {
	// 	logger.Errorf("Validating password...FAILED: %s", err)
	// }

	return "abcxyz", ""
}

func usernameValidator(username string) (bool) {
	logger := utils.ConfigZap()
	matched, err := regexp.MatchString(`[a-z0-9(.)?]+\@vnpay.vn`, username)
	if err != nil {
		logger.Error(err)
	}

	if !matched {
		return matched
	}

	return matched
}

// func passwordValidator(username string) (re bool, err error) {
	
// }