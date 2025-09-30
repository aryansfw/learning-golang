package main

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host     string
	Port     int
	DBDriver string
	DBName   string
	DBUser   string
	DBPass   string
	DBHost   string
	DBPort   int
}

func InitConfig() Config {
	appPort, _ := strconv.Atoi(os.Getenv("APP_PORT"))
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	return Config{
		Host:     os.Getenv("APP_HOST"),
		Port:     appPort,
		DBDriver: os.Getenv("DB_DRIVER"),
		DBName:   os.Getenv("DB_NAME"),
		DBUser:   os.Getenv("DB_USER"),
		DBPass:   os.Getenv("DB_PASS"),
		DBHost:   os.Getenv("DB_HOST"),
		DBPort:   dbPort,
	}
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func (c Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUser, c.DBPass, c.DBName)
}
