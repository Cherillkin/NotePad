package database

import (
	"fmt"

	"github.com/Cherillkin/Notepad/config"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(config *config.EnvConfig, DBMigrate func(db *gorm.DB) error) *gorm.DB {
	uri := fmt.Sprintf(`host=%s port=%s user=%s dbname=%s password=%s sslmode=%s`,
		config.DBHost, config.DBPORT, config.DBUser, config.DBName, config.DBPassword, config.DBSSLMode)

	db, err := gorm.Open(postgres.Open(uri), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	log.Info("Connected to database!")

	if err := DBMigrate(db); err != nil {
		log.Fatalf("Unable to migrate tables: %v", err)
	}

	return db
}
