package main

import (
	"os"
	"test-assigment-ltp/internal/cache"
	"test-assigment-ltp/internal/restapi"
	"test-assigment-ltp/internal/services"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	redisAddr := os.Getenv("REDIS_ADDRESS")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redis := cache.NewRedis(redisAddr)

	kService := services.NewKrakenService(redis)

	ltpHanlder := restapi.NewLTPHanlder(kService)

	e.GET("/api/v1/ltp", ltpHanlder.GetLTP)

	logrus.Info("Starting server on")

	if err := e.Start(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
