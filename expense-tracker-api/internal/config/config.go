package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	App struct {
		Host string
		Port int
	}
	DB struct {
		Host string
		Port int
		Name string
		User string
		Pass string
	}
	JWTSecret string
}

func NewConfig() (*Config, error) {
	appPort, err := strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		return nil, err
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	cfg.App.Host = os.Getenv("APP_HOST")
	cfg.App.Port = appPort

	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = dbPort
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Pass = os.Getenv("DB_PASS")

	cfg.JWTSecret = os.Getenv("JWT_SECRET")

	return cfg, nil
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.App.Host, c.App.Port)
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.DB.Host, c.DB.Port, c.DB.User, c.DB.Pass, c.DB.Name)
}
