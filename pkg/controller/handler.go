package controller

import (
	"net/http"

	"github.com/wreckitkenny/vngitpub/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Handler struct holds required services for handler to function
type Handler struct{}

// Config will hold services that will eventually be injected into this
// handler layer on handler initialization
type Config struct {
    R *gin.Engine
}

// NewHandler initializes the handler with required injected services along with http routes
// Does not return as it deals directly with a reference to the gin Engine
func NewHandler(c *Config) {
  	h := &Handler{}

	r := c.R

    r.POST("/publish", h.Publish)
	r.GET("/loadState", h.LoadState)
	r.GET("/healthz", h.Healthz)
	r.GET("/livez", h.Livez)
	r.GET("/readyz", h.Readyz)
}

// Publish handles incoming requests
func (h *Handler) Publish(c *gin.Context) {
	logger := utils.ConfigZap()
	body, err := c.GetRawData()
	if err != nil {
		logger.Errorf("Reading message body...failed: %s", err)
	} else {
		logger.Debug("Reading message body...ok")
	}
	SendToQueue(body)
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *Handler) LoadState(c *gin.Context) {
	states := LoadState()

	c.JSON(http.StatusOK, states)
}

// Healthz returns /healthz status
func (h *Handler) Healthz(c *gin.Context)  {
	rabbitMQRunning := ValidateRabbitMQConnection()

	if !rabbitMQRunning {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot connect to RabbitMQ"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

// Livez returns /livez status
func (h *Handler) Livez(c *gin.Context)  {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

// Readyz returns /readyz status
func (h *Handler) Readyz(c *gin.Context)  {
	rabbitMQRunning := ValidateRabbitMQConnection()

	if !rabbitMQRunning {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Cannot connect to RabbitMQ"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}