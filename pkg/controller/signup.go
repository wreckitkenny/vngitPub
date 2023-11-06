package controller

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/wreckitkenny/vngitpub/pkg/utils"
	"github.com/wreckitkenny/vngitpub/model"

	"github.com/gin-gonic/gin"
)

// SignUp handles user input for signup
func SignUp(c *gin.Context) (string, error) {
	logger := utils.ConfigZap()
	var signup model.User
	if err := c.BindJSON(&signup); err != nil {
		logger.Error(err)
	}

	username := signup.Username
	email := signup.Email
	password := signup.Password
	fullname := signup.FullName
	department := signup.Department
	role := signup.Role

	h := sha256.New()
	h.Write([]byte(password + "tobtignv"))
	hashedPass := h.Sum(nil)
	encodedPass := base64.URLEncoding.EncodeToString(hashedPass)

	_, findUsernameErr := FindUser("username", username)
	if findUsernameErr == nil {
		return "The username is taken.", findUsernameErr
	}

	_, findEmailErr := FindUser("email", email)
	if findEmailErr == nil {
		return "The email is taken.", findEmailErr
	}

	saveErr := SaveUser(username, email, encodedPass, fullname, department, role)
	if saveErr != nil {
		return "Failed to sign up a new user.", saveErr
	}

	return "Just signed up a new user.", nil
}