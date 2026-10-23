package v1

import (
	"mainstreet/internal/handler"

	"github.com/labstack/echo/v5"
)

func mountHealthRouter(
	routerGrp *echo.Group,
	hdlr *handler.HealthHandler,
) {
	healthGrp := routerGrp.Group("/health")

	healthGrp.GET("/check", hdlr.Check)
}
