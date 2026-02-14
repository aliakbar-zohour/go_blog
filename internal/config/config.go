// config: Loads settings from environment variables and default values.
package config

import (
	"log"
	"os"
	"strconv"
)

const (
	DefaultJWTSecret = "change-me-in-production"
	DefaultBodyLimit = 32 << 20  // 32MB max request body (multipart posts)
	DefaultAuthRate  = 10        // requests per minute per IP for auth
	DefaultListLimit = 20
	MaxListLimit     = 100
)

type Config struct {
	ServerPort      string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	DBName          string
	DBSSL           string
	UploadDir       string
	MaxFileMB       int
	JWTSecret       string
	JWTExpiryHours  int
	CORSOrigins     string
	BodyLimitBytes  int64
	AuthRatePerMin  int
	SMTPHost        string
	SMTPPort        string
	SMTPUser        string
	SMTPPass        string
	SMTPFrom        string
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
	bodyLimit := int64(DefaultBodyLimit)
	if b, _ := strconv.ParseInt(getEnv("BODY_LIMIT_BYTES", ""), 10, 64); b > 0 {
		bodyLimit = b
	}
	authRate, _ := strconv.Atoi(getEnv("AUTH_RATE_PER_MIN", "10"))
	if authRate <= 0 {
		authRate = DefaultAuthRate
	}
	jwtSecret := getEnv("JWT_SECRET", DefaultJWTSecret)
	if jwtSecret == DefaultJWTSecret {
		log.Printf("warning: JWT_SECRET is default; set a strong secret in production")
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
		JWTSecret:      jwtSecret,
		JWTExpiryHours: jwtHours,
		CORSOrigins:    getEnv("CORS_ORIGINS", "*"),
		BodyLimitBytes: bodyLimit,
		AuthRatePerMin: authRate,
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
