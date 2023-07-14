package main

import (
	"fmt"
	"net/http"
	"os"

	"vngitPub/pkg"
	"vngitPub/pkg/controller"
	"vngitPub/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := utils.ConfigZap()
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")
	info := fmt.Sprintf("%s:%s", addr, port)

	//Configure GIN
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.SetTrustedProxies(nil)

	version, err := os.ReadFile("VERSION")
	if err != nil {
		logger.Errorf("Loading version...FAILED: %s", err)
	} else {
		logger.Infof("Loading version...%s", version)
	}

	//Handle request methods
	r.GET("/ping", pkg.Ping)
	r.POST("/publish", func(c *gin.Context) {
		body, err := c.GetRawData()
		if err != nil {
			logger.Errorf("Reading message body...failed: %s", err)
		} else {
			logger.Debug("Reading message body...ok")
		}
		controller.PostHandler(body)
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	//Run server
	logger.Infof("[*] Listening on %s", info)
	r.Run(info)
}
