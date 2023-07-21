package config

import "os"

type Config struct {
	Port  string
	DbUrl string
}

func NewConfig() *Config {
	port := os.Getenv("PORT")
	db := os.Getenv("DB")

	return &Config{
		Port:  port,
		DbUrl: db,
	}
}
