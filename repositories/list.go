package repositories

import (
	"context"

	"github.com/Cherillkin/Notepad/models"
	"gorm.io/gorm"
)

type ListRepository struct {
	db *gorm.DB
}

func (r *ListRepository) CreateList(ctx context.Context, listData *models.List) (*models.List, error) {
	if err := r.db.WithContext(ctx).Create(listData).Error; err != nil {
		return nil, err
	}

	return listData, nil
}

func (r *ListRepository) GetListsByUserID(ctx context.Context, userID uint) ([]models.List, error) {
	var lists []models.List
	if err := r.db.WithContext(ctx).Where("userId = ?", userID).Find(&lists).Error; err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *ListRepository) GetListByID(ctx context.Context, id uint) (*models.List, error) {
	var list models.List
	if err := r.db.WithContext(ctx).First(&list, id).Error; err != nil {
		return nil, err
	}

	return &list, nil
}

func (r *ListRepository) DeleteList(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.List{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *ListRepository) UpdateList(ctx context.Context, id uint, listData *models.List) (*models.List, error) {
	if err := r.db.WithContext(ctx).Model(&models.List{}).Where("id = ?", id).Updates(listData).Error; err != nil {
		return nil, err
	}

	return listData, nil
}

func NewListRepository(db *gorm.DB) models.ListRepository {
	return &ListRepository{
		db: db,
	}
}
