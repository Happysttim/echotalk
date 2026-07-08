package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI           string
	Port               int16
	GoogleClientID     string
	GoogleClientSecret string
	RedirectURL        string
}

func Load() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		slog.Warn("failed to load .env file", "error", err)
		return nil, err
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 16)
	if err != nil {
		slog.Warn("failed to parse PORT", "error", err)
		return nil, err
	}

	return &Config{
		MongoURI:           os.Getenv("MONGODB_URI"),
		Port:               int16(port),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:        os.Getenv("REDIRECT_URL"),
	}, nil
}
