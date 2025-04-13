package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
	JWT   JWTConfig
	SMTP  SMTPConfig
	OTP   OTPConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port int
	URL  string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type OTPConfig struct {
	Expiration time.Duration
	Length     int
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config := &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Env:  viper.GetString("APP_ENV"),
			Port: viper.GetInt("APP_PORT"),
			URL:  viper.GetString("APP_URL"),
		},
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetInt("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetInt("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("JWT_SECRET"),
			Expiration: viper.GetDuration("JWT_EXPIRATION"),
		},
		SMTP: SMTPConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetInt("SMTP_PORT"),
			User:     viper.GetString("SMTP_USER"),
			Password: viper.GetString("SMTP_PASSWORD"),
			From:     viper.GetString("SMTP_FROM"),
		},
		OTP: OTPConfig{
			Expiration: viper.GetDuration("OTP_EXPIRATION"),
			Length:     viper.GetInt("OTP_LENGTH"),
		},
	}

	return config, nil
}
