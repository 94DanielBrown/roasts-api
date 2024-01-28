package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func CorrelationIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		correlationID := generateCorrelationID()
		c.Set("correlationID", correlationID)
		return next(c)
	}
}

func generateCorrelationID() string {
	return uuid.NewString()
}
