package roasts

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// CreateRoastValidator checks for errors when creating a new roast
func CreateRoastValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			errMsg := "Failed to read request body: "
			slog.Error(errMsg, "err", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
		}

		req.Body = io.NopCloser(strings.NewReader(string(body)))

		var reqData struct {
			RoastID    string `json:"roastId"`
			PriceRange string `json:"priceRange"`
			Name       string `json:"name"`
			ImageUrl   string `json:"imageUrl"`
		}

		if err := json.Unmarshal(body, &reqData); err != nil {
			errMsg := "Error unmarshalling JSON"
			slog.Error(errMsg, "err", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
		}

		return next(c)
	}
}
