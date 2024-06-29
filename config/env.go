package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	env := os.Getenv("ENV")
	if env == "" {
		return fmt.Errorf("ENV not set")
	}

	if env == "local" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
		return nil
	}

	if env == "localBinary" {
		envFile := os.Getenv("ENV_FILE")
		if envFile == "" {
			return fmt.Errorf("ENV_FILE not set for localBinary environment")
		}

		err := godotenv.Load(envFile)
		if err != nil {
			return fmt.Errorf("error loading env file: %v", err)
		}
		fmt.Printf("Loaded env variables from %s\n", envFile)
		return nil
	}

	return nil
}
