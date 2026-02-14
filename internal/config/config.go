// config: Loads settings from environment variables and default values.
package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort    string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPass        string
	DBName        string
	DBSSL         string
	UploadDir     string
	MaxFileMB     int
	JWTSecret     string
	JWTExpiryHours int
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPass      string
	SMTPFrom      string
}

func Load() *Config {
	maxMB, _ := strconv.Atoi(getEnv("MAX_UPLOAD_MB", "50"))
	if maxMB <= 0 {
		maxMB = 50
	}
	jwtHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "72"))
	if jwtHours <= 0 {
		jwtHours = 72
	}
	return &Config{
		ServerPort:     getEnv("PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPass:         getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "go_blog"),
		DBSSL:          getEnv("DB_SSLMODE", "disable"),
		UploadDir:      getEnv("UPLOAD_DIR", "uploads"),
		MaxFileMB:      maxMB,
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiryHours: jwtHours,
		SMTPHost:       getEnv("SMTP_HOST", ""),
		SMTPPort:       getEnv("SMTP_PORT", "587"),
		SMTPUser:       getEnv("SMTP_USER", ""),
		SMTPPass:       getEnv("SMTP_PASS", ""),
		SMTPFrom:       getEnv("SMTP_FROM", "noreply@go-blog.local"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
