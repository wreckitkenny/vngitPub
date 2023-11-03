package controller

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/wreckitkenny/vngitpub/model"
	"github.com/wreckitkenny/vngitpub/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Login handles user input for login
func Login(c *gin.Context) (string) {
	logger := utils.ConfigZap()
	var user model.User
	if err := c.BindJSON(&user); err != nil {
		logger.Error(err)
	}

	username := user.Username
	password := user.Password

	userResult, err := FindUser("username", username)
	if err != nil {
		logger.Errorf("Account with username %s not found", username)
		return ""
	}

	// Password check
	h := sha256.New()
	h.Write([]byte(password + "tobtignv"))
	hashedPass := h.Sum(nil)
	encodedPass := base64.URLEncoding.EncodeToString(hashedPass)

	if encodedPass != userResult.Password {
		logger.Errorf("Password of username %s is not matched", username)
		return ""
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = userResult.Email
	claims["fullname"] = userResult.FullName
	claims["role"] = userResult.Role
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("nekottobtignv"))
	if err != nil {
		logger.Errorf("Failed to generate an user token: %s", err)
		return ""
	}

	return t
}