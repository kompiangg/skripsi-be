package http

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
)

type HTTPServerConfig struct {
	Echo *echo.Echo
	Port int
	Host string
}

func Start(config HTTPServerConfig) error {
	errChan := make(chan error, 1)
	go func() {
		config.Echo.HideBanner = true
		config.Echo.HidePort = true
		if err := config.Echo.Start(fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
			errChan <- err
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Println("[ERROR] Error while starting server: ", err.Error())
		return err
	case <-signalChan:
		config.Echo.Shutdown(context.Background())
		return nil
	}
}
