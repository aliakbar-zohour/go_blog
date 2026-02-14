// config: Loads settings from environment variables and default values.
package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	DBSSL     string
	UploadDir  string
	MaxFileMB  int
}

func Load() *Config {
	maxMB, _ := strconv.Atoi(getEnv("MAX_UPLOAD_MB", "50"))
	if maxMB <= 0 {
		maxMB = 50
	}
	return &Config{
		ServerPort: getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPass:     getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "go_blog"),
		DBSSL:      getEnv("DB_SSLMODE", "disable"),
		UploadDir:  getEnv("UPLOAD_DIR", "uploads"),
		MaxFileMB:  maxMB,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
