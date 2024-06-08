package config

import (
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
	Port int
	Host string
}

func NewApp(v *viper.Viper) App {
	return App{
		Port: v.GetInt("app.port"),
		Host: v.GetString("app.host"),
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
		Name:         v.GetString("database.name"),
		Host:         v.GetString("database.host"),
		Port:         v.GetInt("database.port"),
		Password:     v.GetString("database.password"),
		User:         v.GetString("database.user"),
		Timezone:     v.GetString("database.timezone"),
		SslMode:      v.GetString("database.sslmode"),
		MigrationURL: v.GetString("database.migrationURL"),
	}
}

type Token struct {
	SecretKey string
	Duration  time.Duration
}

func NewToken(v *viper.Viper) Token {
	return Token{
		SecretKey: v.GetString("token.secretKey"),
		Duration:  v.GetDuration("token.duration"),
	}
}

type Cloudinary struct {
	Name      string
	ApiKey    string
	ApiSecret string
}

func NewCloudinary(v *viper.Viper) Cloudinary {
	return Cloudinary{
		Name:      v.GetString("cloudinary.name"),
		ApiKey:    v.GetString("cloudinary.apiKey"),
		ApiSecret: v.GetString("cloudinary.apiSecret"),
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
