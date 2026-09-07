package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	CollectionUser        = "users"
	CollectionSession     = "sessions"
	CollectionAnswer      = "answers"
	CollectionSurvey      = "surveys"
	CollectionRateUp      = "rateups"
	CollectionHashCounter = "hashcounters"
)

type config struct {
	WebURL             string
	Host               string
	MongoURI           string
	MongoDB            string
	Port               int16
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	SecretKey          string
	SmtpHost           string
	SmtpPort           int16
	MailId             string
	MailPassword       string
	MailFrom           string
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

	smtpPort, err := strconv.ParseInt(os.Getenv("SMTP_PORT"), 10, 16)
	if err != nil {
		panic("failed to parse SMTP_PORT environment variable: " + err.Error())
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic("failed to parse REDIS_DB environment variable: " + err.Error())
	}

	Config = &config{
		WebURL:             os.Getenv("WEB_URL"),
		Host:               os.Getenv("HOST"),
		MongoURI:           os.Getenv("MONGODB_URI"),
		MongoDB:            os.Getenv("MONGODB"),
		Port:               int16(port),
		RedisAddr:          os.Getenv("REDIS_ADDR"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		RedisDB:            redisDB,
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		SecretKey:          os.Getenv("SECRET_KEY"),
		SmtpHost:           os.Getenv("SMTP_HOST"),
		SmtpPort:           int16(smtpPort),
		MailId:             os.Getenv("MAIL_ID"),
		MailPassword:       os.Getenv("MAIL_PASSWORD"),
		MailFrom:           os.Getenv("MAIL_FROM"),
	}
}
