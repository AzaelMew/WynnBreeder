package main

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port       string
	DBPath     string
	SessionTTL time.Duration
}

func loadConfig() Config {
	port := getEnv("WYNNBREEDER_PORT", "8080")
	dbPath := getEnv("WYNNBREEDER_DB", "./wynnbreeder.db")

	ttlDays := 30
	if v := os.Getenv("WYNNBREEDER_SESSION_DAYS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			ttlDays = n
		}
	}

	return Config{
		Port:       port,
		DBPath:     dbPath,
		SessionTTL: time.Duration(ttlDays) * 24 * time.Hour,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
