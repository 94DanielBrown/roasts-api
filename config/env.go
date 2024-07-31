package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	TableName   string
	ImageBucket string
	WebPort     int
}

func LoadEnvVariables() (Env, error) {
	err := getEnv()
	if err != nil {
		return Env{}, err
	}

	webPort, err := strconv.Atoi(os.Getenv("WEB_PORT"))
	if err != nil {
		webPort = 8000
	}

	return Env{
		TableName:   os.Getenv("TABLE_NAME"),
		ImageBucket: os.Getenv("IMAGE_BUCKET"),
		WebPort:     webPort,
	}, nil
}

func getEnv() error {
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
