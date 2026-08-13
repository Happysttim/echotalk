package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	CollectionUser    = "users"
	CollectionSession = "sessions"
	CollectionAnswer  = "answers"
	CollectionSurvey  = "survey"
)

type config struct {
	MongoURI           string
	DBName             string
	Port               int16
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	SecretKey          string
}

var Config *config

func init() {
	if Config != nil {
		return
	}

	if err := godotenv.Load(".env"); err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	port, err := strconv.ParseInt(os.Getenv("PORT"), 10, 16)
	if err != nil {
		panic("failed to parse PORT environment variable: " + err.Error())
	}

	Config = &config{
		MongoURI:           os.Getenv("MONGODB_URI"),
		DBName:             os.Getenv("DB_NAME"),
		Port:               int16(port),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		SecretKey:          os.Getenv("SECRET_KEY"),
	}
}
