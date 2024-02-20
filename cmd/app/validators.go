package main

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func (app *Config) CreateRoastValidator(next echo.HandlerFunc) echo.HandlerFunc {
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

func (app *Config) CreateReviewValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		body, err := io.ReadAll(req.Body)
		if err != nil {
			errMsg := "Failed to read request body"
			slog.Error(errMsg, "err", err)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": errMsg})
		}

		// Reset the request body so it can be read again by the main handler
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
