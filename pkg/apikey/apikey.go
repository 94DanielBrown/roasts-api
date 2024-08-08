package apikey

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// APIKeyMiddleware checks for a valid API key in the request header
func Validate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("X-API-Key")
			if apiKey == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "API key is missing"})
			}

			if apiKey != os.Getenv("API_KEY") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid API key"})
			}

			return next(c)
		}
	}
}
