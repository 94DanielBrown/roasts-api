package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"io"
	"strings"
)

func (app *Config) CreateRoastValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		// Read the body to a variable
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Error("Failed to read request body: ", err)
			return err
		}

		req.Body = io.NopCloser(strings.NewReader(string(body)))

		var bodyBytes []byte
		if c.Request().Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request().Body)
			c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			// parse json
			reqData := struct {
				RoastID    string `json:"roastId"`
				PriceRange string `dynamodbav:"SK" json:"priceRange"`
				Name       string `dynamodbav:"Name" json:"name"`
				ImageUrl   string `dynamodbav:"ImageUrl" json:"imageUrl"`
			}{}
			err := json.Unmarshal(bodyBytes, &reqData)
			if err != nil {
				return c.JSON(400, "error json")
			}
			fmt.Println(reqData.RoastID)
		}

		// Call handler
		return next(c)
	}
}
