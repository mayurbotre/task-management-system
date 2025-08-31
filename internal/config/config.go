package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	Port string

	DatabaseDSN string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBParams    string
}

func Load() *Config {
	cfg := &Config{
		Port: getEnv("PORT", "8080"),

		DatabaseDSN: os.Getenv("DATABASE_DSN"),

		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3307"),
		DBUser:     getEnv("DB_USER", "appuser"),
		DBPassword: getEnv("DB_PASSWORD", "apppass"),
		DBName:     getEnv("DB_NAME", "tasksdb"),
		DBParams:   getEnv("DB_PARAMS", "parseTime=true&charset=utf8mb4&loc=Local"),
	}

	log.Printf("[config] PORT=%s DB_HOST=%s DB_PORT=%s DB_NAME=%s DSN_SET=%v",
		cfg.Port, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DatabaseDSN != "")

	return cfg
}

func (c *Config) MySQLDSN() string {
	if c.DatabaseDSN != "" {
		return c.DatabaseDSN
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBParams,
	)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
