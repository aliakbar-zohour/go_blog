// database: PostgreSQL connection pool and auto-migration of models.
package database

import (
	"fmt"
	"log"
	"time"

	"github.com/aliakbar-zohour/go_blog/internal/config"
	"github.com/aliakbar-zohour/go_blog/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSL,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	if err := db.AutoMigrate(&model.Post{}, &model.Media{}, &model.Author{}, &model.Category{}, &model.Comment{}, &model.EmailVerification{}); err != nil {
		log.Printf("warning: automigrate: %v", err)
	}
	return db, nil
}
