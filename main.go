package main

import (
	"fmt"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wreckitkenny/vngitpub/pkg/controller"
	"github.com/wreckitkenny/vngitpub/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := utils.ConfigZap()

	//Configure GIN
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(controller.CORSMiddleware())
	utils.GetVersion()
	controller.ValidateMongoConnection()

	controller.NewHandler(&controller.Config{
		R: router,
	})

	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")
	info := fmt.Sprintf("%s:%s", addr, port)

	srv := &http.Server{
	Addr:    info,
	Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Initializing server...FAILED: %v", err)
		}
	}()

	logger.Infof("[*] Listening on %s", srv.Addr)

	quit := make(chan os.Signal)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// This blocks until a signal is passed into the quit channel
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %s", err)
	}
}
