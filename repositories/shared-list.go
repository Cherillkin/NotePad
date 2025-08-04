package repositories

import (
	"context"

	"github.com/Cherillkin/Notepad/models"
	"gorm.io/gorm"
)

type SharedListRepository struct {
	db *gorm.DB
}

func (r *SharedListRepository) SharedList(ctx context.Context, listID, userID uint) error {
	shared := models.SharedList{ListID: listID, UserID: userID}
	return r.db.WithContext(ctx).Create(&shared).Error
}

func (r *SharedListRepository) GetSharedLists(ctx context.Context, userID uint) ([]models.List, error) {
	var lists []models.List

	err := r.db.WithContext(ctx).Joins("JOIN shared_lists ON shared_lists.list_id = list_id").
		Where("shared_lists.user_id = ?", userID).Find(&lists).Error

	return lists, err
}

func NewSharedListRepository(db *gorm.DB) models.SharedListRepository {
	return &SharedListRepository{
		db: db,
	}
}
