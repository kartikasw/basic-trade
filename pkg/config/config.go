package config

import (
	"encoding/base64"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App        App
	Database   Database
	Token      Token
	Cloudinary Cloudinary
}

func NewConfig(v *viper.Viper) Config {
	return Config{
		App:        NewApp(v),
		Database:   NewDatabase(v),
		Token:      NewToken(v),
		Cloudinary: NewCloudinary(v),
	}
}

type App struct {
	Port    int
	Host    string
	Timeout time.Duration
}

func NewApp(v *viper.Viper) App {
	return App{
		Port:    v.GetInt("APP_PORT"),
		Host:    v.GetString("APP_HOST"),
		Timeout: v.GetDuration("APP_TIMEOUT"),
	}
}

type Database struct {
	Name         string
	Host         string
	Port         int
	Password     string
	User         string
	Timezone     string
	SslMode      string
	MigrationURL string
}

func NewDatabase(v *viper.Viper) Database {
	return Database{
		Name:         v.GetString("DB_NAME"),
		Host:         v.GetString("DB_HOST"),
		Port:         v.GetInt("DB_PORT"),
		Password:     v.GetString("DB_PASSWORD"),
		User:         v.GetString("DB_USER"),
		Timezone:     v.GetString("DB_TIMEZONE"),
		SslMode:      v.GetString("DB_SSLMODE"),
		MigrationURL: v.GetString("DB_MIGRATION_URL"),
	}
}

type Token struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
	PublicKey            string
	PrivateKey           string
}

func NewToken(v *viper.Viper) Token {
	encodedPublicKey := v.GetString("TOKEN_PUBLIC_KEY")

	publicKey, err := base64.StdEncoding.DecodeString(encodedPublicKey)
	if err != nil {
		log.Fatal("Couldn't encode public key")
	}

	encodedPrivateKey := v.GetString("TOKEN_PRIVATE_KEY")
	privateKey, err := base64.StdEncoding.DecodeString(encodedPrivateKey)
	if err != nil {
		log.Fatal("Couldn't encode private key")
	}

	return Token{
		AccessTokenDuration:  v.GetDuration("TOKEN_ACCESS_TOKEN_DURATION"),
		RefreshTokenDuration: v.GetDuration("TOKEN_REFRESH_TOKEN_DURATION"),
		PublicKey:            string(publicKey),
		PrivateKey:           string(privateKey),
	}
}

type Cloudinary struct {
	Name      string
	ApiKey    string
	ApiSecret string
}

func NewCloudinary(v *viper.Viper) Cloudinary {
	return Cloudinary{
		Name:      v.GetString("CLOUDINARY_NAME"),
		ApiKey:    v.GetString("CLOUDINARY_API_KEY"),
		ApiSecret: v.GetString("CLOUDINARY_API_SECRET"),
	}
}

func LoadConfig(path string) Config {
	v := viper.New()
	v.SetConfigFile(path)

	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("Load config error: ", err.Error())
	}

	return NewConfig(v)
}

func LoadTestConfig(path string) Config {
	v := viper.New()
	v.SetConfigFile(path)

	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("Load config error: ", err.Error())
	}

	return NewConfig(v)
}
