package roasts

import (
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

		//TODO - add validation for the request body
		//var reqData struct {
		//	RoastID    string `json:"roastId"`
		//	PriceRange string `json:"priceRange"`
		//	Name       string `json:"name"`
		//	ImageUrl   string `json:"imageUrl"`
		//}
		//slog.Info("Request Body: ", "body", req.Body)
		//
		//if err := json.Unmarshal(body, &reqData); err != nil {
		//	errMsg := "error creating roast due to unmarshalling of json"
		//	slog.Error(errMsg, "error", err)
		//	return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
		//}

		return next(c)
	}
}
