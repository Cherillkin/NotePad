package repositories

import (
	"context"

	"github.com/Cherillkin/Notepad/models"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
}

func (r *ItemRepository) CreateItem(ctx context.Context, itemData *models.Item) (*models.Item, error) {
	if err := r.db.WithContext(ctx).Create(itemData).Error; err != nil {
		return nil, err
	}

	return itemData, nil
}

func (r *ItemRepository) GetListItems(ctx context.Context, listID uint) ([]models.Item, error) {
	var items []models.Item
	if err := r.db.WithContext(ctx).Where("list_id = ?", listID).Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) GetItem(ctx context.Context, id uint) (*models.Item, error) {
	var item models.Item
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *ItemRepository) DeleteItem(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Item{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) UpdateItem(ctx context.Context, id uint, itemData *models.Item) (*models.Item, error) {
	if err := r.db.WithContext(ctx).Model(&models.Item{}).Where("id = ?", id).Updates(itemData).Error; err != nil {
		return nil, err
	}

	return itemData, nil
}

func NewItemRepository(db *gorm.DB) models.ItemRepository {
	return &ItemRepository{
		db: db,
	}
}
