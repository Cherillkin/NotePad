package database

import (
	"github.com/Cherillkin/Notepad/models"
	"gorm.io/gorm"
)

func DBMigrator(db *gorm.DB) error {
	return db.AutoMigrate(&models.User{}, &models.List{}, &models.Item{}, &models.SharedList{})
}
