package main

import (
	"ARIGATOBANK/handler"
	"ARIGATOBANK/repository"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	registerHandlers(router)
	loggingSetting()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Println("Failed to gracefully shutdown:", err)
	}
	log.Println("Server shutdown")
}

func registerHandlers(r *gin.Engine) {
	balanceHandler := new(handler.Balance)
	balanceHandler.DB = repository.NewMySql()
	v1 := r.Group("/v1")
	{
		v1.POST("/balance", balanceHandler.Update)
	}
}

func loggingSetting() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logrus.SetLevel(logrus.InfoLevel)
		return
	}
	parsedLogLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
		return
	}
	logrus.SetLevel(parsedLogLevel)
}
