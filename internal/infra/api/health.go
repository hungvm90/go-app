package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type HealthController struct {
	version string
}

func NewHealthController(version string) *HealthController {
	return &HealthController{version: version}
}

func (handler *HealthController) Init(ec *echo.Echo) {
	ec.GET("/", handler.Ping)
}

func (handler *HealthController) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"version": handler.version,
	})
}
