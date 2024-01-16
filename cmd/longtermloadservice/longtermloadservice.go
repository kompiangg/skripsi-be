package longtermloadservice

import (
	"context"

	"skripsi-be/config"
	"skripsi-be/pkg/http"
	"skripsi-be/service"

	"skripsi-be/cmd/longtermloadservice/handler"
	inmiddleware "skripsi-be/cmd/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func Init(
	service service.Service,
	config config.LongTermLoadService,
	mw inmiddleware.Middleware,
) error {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: config.WhiteListAllowOrigin,
		},
	))

	handler.Init(e, service, mw)

	log.Info().Msgf("Starting Auth Service HTTP server on %s:%d", config.Host, config.Port)
	err := http.Start(http.HTTPServerConfig{
		Echo: e,
		Port: config.Port,
		Host: config.Host,
	})
	if err != nil {
		return err
	}

	log.Info().Msg("Starting graceful shutdown HTTP Server...")

	err = e.Shutdown(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Error while shutting down HTTP server")
		return err
	}

	log.Info().Msg("HTTP Server shutdown gracefully, RIP 🙏")

	return nil
}
