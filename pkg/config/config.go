package config

import (
	"basic-trade/common"
	"encoding/base64"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App        App
	Database   Database
	Token      Token
	Cloudinary Cloudinary
}

func NewConfig() Config {
	return Config{
		App:        NewApp(),
		Database:   NewDatabase(),
		Token:      NewToken(),
		Cloudinary: NewCloudinary(),
	}
}

type App struct {
	Port    string
	Timeout time.Duration
	GinMode string
}

func NewApp() App {
	timeoutStr := os.Getenv("APP_TIMEOUT")
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatal("Couldn't parse Timeout")
	}

	return App{
		Port:    os.Getenv("PORT"),
		Timeout: timeout,
		GinMode: os.Getenv("APP_GIN_MODE"),
	}
}

type Database struct {
	Name         string
	Host         string
	Port         string
	Password     string
	User         string
	Timezone     string
	SslMode      string
	MigrationURL string
}

func NewDatabase() Database {
	return Database{
		Name:         os.Getenv("DB_NAME"),
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		Password:     os.Getenv("DB_PASSWORD"),
		User:         os.Getenv("DB_USER"),
		Timezone:     os.Getenv("DB_TIMEZONE"),
		SslMode:      os.Getenv("DB_SSLMODE"),
		MigrationURL: os.Getenv("DB_MIGRATION_URL"),
	}
}

type Token struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	PublicKey            string
	PrivateKey           string
}

func NewToken() Token {
	encodedPublicKey := os.Getenv("TOKEN_PUBLIC_KEY")

	publicKey, err := base64.StdEncoding.DecodeString(encodedPublicKey)
	if err != nil {
		log.Fatal("Couldn't encode public key")
	}

	encodedPrivateKey := os.Getenv("TOKEN_PRIVATE_KEY")
	privateKey, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		log.Fatal("Couldn't encode private key")
	}

	accessDurationStr := os.Getenv("TOKEN_ACCESS_TOKEN_DURATION")
	accessDuration, err := time.ParseDuration(accessDurationStr)
	if err != nil {
		log.Fatal("Couldn't parse Access Token Duration")
	}

	refreshDurationStr := os.Getenv("TOKEN_REFRESH_TOKEN_DURATION")
	refreshDuration, err := time.ParseDuration(refreshDurationStr)
	if err != nil {
		log.Fatal("Couldn't parse Access Token Duration")
	}

	return Token{
		AccessTokenDuration:  accessDuration,
		RefreshTokenDuration: refreshDuration,
		PublicKey:            string(publicKey),
		PrivateKey:           string(privateKey),
	}
}

type Cloudinary struct {
	Name      string
	ApiKey    string
	ApiSecret string
}

func NewCloudinary() Cloudinary {
	return Cloudinary{
		Name:      os.Getenv("CLOUDINARY_NAME"),
		ApiKey:    os.Getenv("CLOUDINARY_API_KEY"),
		ApiSecret: os.Getenv("CLOUDINARY_API_SECRET"),
	}
}

func LoadConfig() Config {
	if common.DEVELOPMENT {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return NewConfig()
}
