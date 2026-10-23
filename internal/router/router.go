package router

import (
	"mainstreet/internal/handler"
	"mainstreet/internal/middleware"
	v1 "mainstreet/internal/router/v1"

	"github.com/labstack/echo/v5"
)

func New(mw *middleware.Middleware, hdlr *handler.Handler) *echo.Echo {
	e := echo.New()

	apiGrp := e.Group("/api")

	v1.MountV1Router(apiGrp, hdlr)

	return e
}
