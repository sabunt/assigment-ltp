package main

import (
	"net/http"
	"os"
	"strconv"
	"test-assigment-ltp/internal/restapi"
	"test-assigment-ltp/internal/services"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	httpclient := http.Client{Timeout: 10 * time.Second}
	kService := services.NewKrakenService(&httpclient)

	concurrencyLimit := 5 // Default
	if val, exists := os.LookupEnv("MAX_CONCURRENCY"); exists {
		if limit, err := strconv.Atoi(val); err == nil {
			concurrencyLimit = limit
		}
	}

	ltpHanlder := restapi.NewLTPHanlder(kService, concurrencyLimit)

	e.GET("/api/v1/ltp", ltpHanlder.GetLTP)

	logrus.Info("Starting server on")

	if err := e.Start(":8080"); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}
