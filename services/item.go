package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Cherillkin/Notepad/models"
	"gorm.io/gorm"
)

type ItemService struct {
	repository models.ItemRepository
}

func (s *ItemService) CreateItem(ctx context.Context, listID uint, itemData *models.Item) (*models.Item, error) {
	itemData.ListID = listID
	itemData.Created_At = time.Now()
	itemData.Updated_At = time.Now()

	newItem, err := s.repository.CreateItem(ctx, itemData)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return newItem, nil
}

func (s *ItemService) GetListItems(ctx context.Context, listID uint) ([]models.Item, error) {
	items, err := s.repository.GetListItems(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items list: %w", err)
	}

	return items, nil
}

func (s *ItemService) GetItem(ctx context.Context, listID, itemID uint) (*models.Item, error) {
	item, err := s.repository.GetItem(ctx, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("item not found")
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if item.ListID != listID {
		return nil, fmt.Errorf("list not found")
	}

	return item, nil
}

func (s *ItemService) DeleteItem(ctx context.Context, listID, itemID uint) error {
	item, err := s.repository.GetItem(ctx, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("item not found")
		}

		return fmt.Errorf("failed to get item: %w", err)
	}

	if item.ListID != listID {
		return fmt.Errorf("list not found")
	}

	if err := s.repository.DeleteItem(ctx, itemID); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

func (s *ItemService) UpdateItem(ctx context.Context, listID, itemID uint, updateData *models.Item) (*models.Item, error) {
	existingItem, err := s.repository.GetItem(ctx, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("item not found")
		}

		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if existingItem.ListID != listID {
		return nil, fmt.Errorf("list not found")
	}

	if updateData.Name != "" {
		existingItem.Name = updateData.Name
	}

	if updateData.Description != "" {
		existingItem.Description = updateData.Description
	}

	existingItem.Updated_At = time.Now()

	updateItem, err := s.repository.UpdateItem(ctx, itemID, existingItem)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	return updateItem, nil
}

func NewItemService(repository models.ItemRepository) models.ItemService {
	return &ItemService{
		repository: repository,
	}
}
