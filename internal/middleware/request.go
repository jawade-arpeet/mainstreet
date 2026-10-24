package middleware

import (
	"errors"
	"mainstreet/internal/errs"

	"github.com/labstack/echo/v5"
)

type RequestMiddleware struct{}

func newRequestMiddleware() *RequestMiddleware {
	return &RequestMiddleware{}
}

func (m *RequestMiddleware) HandleHTTPErr(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx *echo.Context) error {
		err := next(ctx)
		if err == nil {
			return nil
		}

		var httpErr *errs.HTTPError
		if !errors.As(err, &httpErr) {
			return errs.Internal(err)
		}

		return ctx.JSON(
			httpErr.Status,
			map[string]any{
				"error": httpErr,
			},
		)
	}
}
