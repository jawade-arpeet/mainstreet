package v1

import (
	"mainstreet/internal/handler"

	"github.com/labstack/echo/v5"
)

func mountAuthRoutes(
	routerGrp *echo.Group,
	authHdlr *handler.AuthHandler,
) {
	authGrp := routerGrp.Group("/auth")

	authGrp.POST("/sign-up", authHdlr.SignUp)
	authGrp.POST("/sign-in", authHdlr.SignIn)
}
