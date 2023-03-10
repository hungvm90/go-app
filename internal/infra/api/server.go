package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hungvm90/go-app/internal"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func StartServer(appConfig internal.AppConfig) {
	ec := echo.New()
	setupMiddleware(ec, appConfig)

	rootContext, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	appContext := internal.AppContext{Wg: wg, Context: rootContext}

	healthController := NewHealthController(appConfig.Version)
	healthController.Init(ec)

	go func() {
		if err := ec.Start(fmt.Sprintf(":%d", appConfig.Port)); err != nil && err != http.ErrServerClosed {
			log.Error().Stack().Err(err).Msgf("shutting down the server")
			ec.Logger.Fatal("shutting down the server")
		}
	}()

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interruptChan
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := ec.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msgf("failed to shutdown api server gracefully")
	}
	cancelFunc()
	appContext.Wg.Wait()
	log.Info().Msgf("shutdown api successfully!")

}

func setupMiddleware(e *echo.Echo, appConfig internal.AppConfig) {
	setUpLoggerMiddleware(e)
	e.Use(middleware.Recover())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Skipper: middleware.DefaultSkipper,
		Generator: func() string {
			return uuid.New().String()
		},
		RequestIDHandler: func(e echo.Context, s string) {
			e.Set("requestId", s)
		},
		TargetHeader: echo.HeaderXRequestID,
	}))
	e.HTTPErrorHandler = func(err error, context echo.Context) {
		log.Warn().Stack().Err(err).Msgf("process %s with error", context.Path())
		//custom here
		e.DefaultHTTPErrorHandler(err, context)
	}
}

func setUpLoggerMiddleware(e *echo.Echo) {
	logConfig := middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
		Output:           os.Stdout,
	}
	e.Use(middleware.LoggerWithConfig(logConfig))
}

func createRequestContext(c echo.Context) internal.RequestContext {
	requestId := ""
	temp := c.Get("requestId")
	if temp == nil {
		requestId = uuid.New().String()
	} else {
		requestId = temp.(string)
	}
	logger := log.With().Str("requestId", requestId).Logger()
	return internal.RequestContext{Logger: &logger, Context: c.Request().Context()}
}
