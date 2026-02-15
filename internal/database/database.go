// database: PostgreSQL connection and auto-migration of models.
package database

import (
	"fmt"
	"log"

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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}
	if err := db.AutoMigrate(&model.Author{}, &model.Category{}, &model.Post{}, &model.Media{}, &model.Comment{}); err != nil {
		log.Printf("warning: automigrate: %v", err)
	}
	return db, nil
}
