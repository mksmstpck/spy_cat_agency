package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Config struct {
	PostgregUrl  string
	TheCatApiUrl string
	Port         int
}

func NewConfig() Config {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		logrus.Error(err)
	}
	return Config{
		PostgregUrl:  os.Getenv("POSTGRES_URL"),
		TheCatApiUrl: os.Getenv("THE_CAT_API_URL"),
		Port:         port,
	}
}
