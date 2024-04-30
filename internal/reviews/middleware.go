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
			errMsg := "failed to read request body"
			slog.Error(errMsg, "error", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
		}

		// Reset the request body, so it can be read again by the main handler
		req.Body = io.NopCloser(bytes.NewBuffer(body))

		var reqData struct {
			UserID         string `json:"userID"`
			Comment        string `json:"comment,omitempty"`
			RoastName      string `json:"roastName"`
			ImageURL       string `json:"imageURL"`
			OverallRating  int    `json:"overallRating"`
			MeatRating     int    `json:"meatRating"`
			PotatoesRating int    `json:"potatoesRating"`
			VegRating      int    `json:"vegRating"`
			GravyRating    int    `json:"gravyRating"`
		}

		if err := json.Unmarshal(body, &reqData); err != nil {
			errMsg := "error unmarshalling json"
			slog.Error(errMsg, "error", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
		}

		ratings := []int{
			reqData.OverallRating,
			reqData.MeatRating,
			reqData.PotatoesRating,
			reqData.VegRating,
			reqData.GravyRating,
		}

		for _, rating := range ratings {
			if rating < 1 || rating > 10 {
				errMsg := "invalid rating: ratings should be between 1 and 10"
				slog.Error(errMsg)
				return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
			}
		}

		return next(c)
	}
}
