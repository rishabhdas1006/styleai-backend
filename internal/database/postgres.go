package database

import (
	"fmt"

	"styleai-backend/internal/config"
	"styleai-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(cfg *config.Config) error {

	if cfg.Database.AutoCreate {
		if err := EnsureDatabaseExists(cfg); err != nil {
			return err
		}
	}

	return ConnectDB(cfg)
}

func EnsureDatabaseExists(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%d sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect default database: %w", err)
	}

	var exists bool

	query := `
		SELECT EXISTS (
			SELECT FROM pg_database WHERE datname = ?
		)
	`

	if err := db.Raw(query, cfg.Database.DBName).Scan(&exists).Error; err != nil {
		return fmt.Errorf("check database existence: %w", err)
	}

	if !exists {
		createDBSQL := fmt.Sprintf("CREATE DATABASE %s", cfg.Database.DBName)

		if err := db.Exec(createDBSQL).Error; err != nil {
			return fmt.Errorf("create database: %w", err)
		}
	}

	return nil
}

func ConnectDB(cfg *config.Config) error {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	DB = db

	if err := runMigrations(); err != nil {
		return err
	}

	return nil
}

func runMigrations() error {

	err := DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.ProductVariant{},
		&models.VariantImage{},
		&models.Category{},
		&models.Cart{},
		&models.CartItem{},
	)

	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

func Close() error {

	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
