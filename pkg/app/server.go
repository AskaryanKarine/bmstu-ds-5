package app

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const gracefulShutdownDeadline = 10 * time.Second

func Run(echo *echo.Echo, port int) {
	portStr := fmt.Sprintf(":%d", port)
	go func() {
		log.Info("server starting on", "port", portStr)
		if err := echo.Start(portStr); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownDeadline)
	defer cancel()

	log.Info("server shutting down")
	if err := echo.Server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

func AddHealthCheck(e *echo.Echo) {
	e.GET("/manage/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{})
	})
}
