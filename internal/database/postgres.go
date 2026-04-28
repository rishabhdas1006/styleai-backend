package database

import (
	"fmt"
	"log"

	"styleai-backend/internal/config"
	"styleai-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func EnsureDatabaseExists(cfg *config.Config) {
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
		log.Fatal("Failed to connect to default DB:", err)
	}

	var exists bool

	query := `
		SELECT EXISTS (
			SELECT FROM pg_database WHERE datname = ?
		)
	`

	if err := db.Raw(query, cfg.Database.DBName).Scan(&exists).Error; err != nil {
		log.Fatal("Failed to check database existence:", err)
	}

	if !exists {
		createDBSQL := fmt.Sprintf("CREATE DATABASE %s", cfg.Database.DBName)
		if err := db.Exec(createDBSQL).Error; err != nil {
			log.Fatal("Failed to create database:", err)
		}
		log.Println("Database created:", cfg.Database.DBName)
	} else {
		log.Println("Database already exists:", cfg.Database.DBName)
	}
}

func ConnectDB(cfg *config.Config) {

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
		log.Fatal("Failed to connect database:", err)
	}

	DB = db

	log.Println("PostgreSQL connected")

	runMigrations()
}

func runMigrations() {

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
		log.Fatal("Migration failed:", err)
	}

	log.Println("Database migrated")
}
