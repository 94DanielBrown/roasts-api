package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	env := os.Getenv("ENV")
	switch env {
	case "local":
		err := godotenv.Load()
		if err != nil {
			return err
		}
		return nil
	case "dev":
		return nil
	default:
		return fmt.Errorf("ENV variable not recognized")
	}
}
