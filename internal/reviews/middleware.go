package reviews

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateReviewValidator checks for errors when creating a new roast
func CreateReviewValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			errMsg := "Failed to read request body"
			slog.Error(errMsg, "err", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
		}

		// Reset the request body, so it can be read again by the main handler
		req.Body = io.NopCloser(bytes.NewBuffer(body))

		var reqData struct {
			UserID    string `json:"userID"`
			Rating    int    `json:"rating"`
			Comment   string `json:"comment,omitempty"`
			RoastName string `json:"roastName"`
			ImageUrl  string `json:"imageUrl"`
		}

		if err := json.Unmarshal(body, &reqData); err != nil {
			errMsg := "Error unmarshalling JSON"
			slog.Error(errMsg, "err", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
		}

		if reqData.Rating < 1 || reqData.Rating > 10 {
			errMsg := "Invalid rating"
			slog.Error(errMsg, "err", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
		}

		return next(c)
	}
}
