package main

import (
	"fmt"
	"net/http"
	"os"

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
	r.GET("/healthz", func (c *gin.Context)  {
		rabbitMQRunning := controller.ValidateRabbitMQConnection()

		if !rabbitMQRunning {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot connect to RabbitMQ"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
	r.GET("/livez", func (c *gin.Context)  {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
	r.GET("/readyz", func (c *gin.Context)  {
		rabbitMQRunning := controller.ValidateRabbitMQConnection()

		if !rabbitMQRunning {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot connect to RabbitMQ"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	if controller.ValidateRabbitMQConnection() {
		logger.Infof("[*] Listening on %s", info)
		r.Run(info)
	}
}
