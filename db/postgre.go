package db

import (
	"fmt"

	"github.com/fmelihh/product-hunt-graph-visualize/config"
	"github.com/fmelihh/product-hunt-graph-visualize/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreDb struct {
	GormDB *gorm.DB
}

func NewPostgreSqlDb(cfg *config.Config) (*PostgreDb, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Europe/Istanbul",
		cfg.PostgreHost, cfg.PostgreUser, cfg.PostgrePassword, cfg.PostgreDbName, cfg.PostgrePort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PostgreDb{GormDB: db}, nil
}

func (d *PostgreDb) MigrateDatabaseModels() {
	err := d.GormDB.AutoMigrate(&models.BaseUrl{}, &models.EntityUrl{})
	if err != nil {
		panic(fmt.Errorf("failed to auto-migrate database schema: %w", err))
	}
}

func (d *PostgreDb) Close() error {
	sqlDB, err := d.GormDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}
	return sqlDB.Close()
}
