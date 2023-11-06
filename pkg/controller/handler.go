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

    r.POST("/publish", h.HandlePublish)
	r.GET("/loadState", h.HandleLoadState)
	r.POST("/login", h.HandleLogin)
	r.POST("/signup", h.HandleSignup)
	r.GET("/statistic", h.HandleStatistic)

	// Healthcheck
	r.GET("/healthz", h.Healthz)
	r.GET("/livez", h.Livez)
	r.GET("/readyz", h.Readyz)
}

// HandlePublish handles incoming requests
func (h *Handler) HandlePublish(c *gin.Context) {
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

// HandleLoadState returns all stored data in DB
func (h *Handler) HandleLoadState(c *gin.Context) {
	bearer := c.Request.Header["Authorization"]
	isValid, err := utils.IsValidToken(bearer)
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error":err})
	} else {
		c.JSON(http.StatusOK, LoadState())
	}
}

// HandleSignup authenticates users
func (h *Handler) HandleSignup(c *gin.Context) {
	logger := utils.ConfigZap()
	msg, err := SignUp(c)
	if err != nil {
		logger.Warn(msg)
		logger.Errorf("Signing up a new user...FAILED: %s", err)
		c.JSON(http.StatusOK, map[string]string{
			"Msg": msg,
		})
	}

	c.JSON(http.StatusOK, map[string]string{
		"Msg": msg,
	})
}

// HandleLogin authenticates users
func (h *Handler) HandleLogin(c *gin.Context) {
	token := Login(c)
	if token != "" {
		c.JSON(http.StatusOK, map[string]string{
			"token": token,
		})
	} else {
		c.JSON(http.StatusOK, map[string]string{
			"msg": "Incorrect login credentials. Please try again.",
		})
	}
}

// HandleStatistic handles today's
func (h *Handler) HandleStatistic(c *gin.Context) {
	bearer := c.Request.Header["Authorization"]
	isValid, err := utils.IsValidToken(bearer)
	if !isValid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, map[string]interface{} {
			"total": Total(),
			"today": Today(),
			"graph": Graph(),
		})
	}
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