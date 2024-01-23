package database

import (
	"fmt"
	"github.com/aliazizii/IMDB-In-Golang/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DSN(cfg config.Database) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
}

func New(cfg config.Database) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DSN(cfg)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
