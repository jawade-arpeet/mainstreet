package handler

import (
	"mainstreet/internal/dto"
	"mainstreet/internal/service"
	"net/http"

	"github.com/labstack/echo/v5"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func newAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) SignUp(ctx *echo.Context) error {
	var req dto.SignUpRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if err := h.authSvc.SignUp(
		ctx.Request().Context(),
		req.Email,
		req.Password,
	); err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"message": "successfully created account"})
}

func (h *AuthHandler) SignIn(ctx *echo.Context) error {
	var req dto.SignInRequest
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	accID, err := h.authSvc.SignIn(
		ctx.Request().Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"access_token": accID,
		"token_type":   "Bearer",
	})
}
