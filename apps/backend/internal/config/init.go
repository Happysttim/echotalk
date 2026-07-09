package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type config struct {
	MongoURI           string
	Port               int16
	GoogleClientID     string
	GoogleClientSecret string
	RedirectURL        string
}

var Config *config

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 16)
	if err != nil {
		panic("failed to parse PORT environment variable: " + err.Error())
	}

	Config = &config{
		MongoURI:           os.Getenv("MONGODB_URI"),
		Port:               int16(port),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:        os.Getenv("REDIRECT_URL"),
	}
}
