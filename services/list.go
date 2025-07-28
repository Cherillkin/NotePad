package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Cherillkin/Notepad/models"
	"gorm.io/gorm"
)

type ListService struct {
	repository models.ListRepository
}

func (s *ListService) CreateList(ctx context.Context, userID uint, listData *models.List) (*models.List, error) {
	listData.UserID = userID
	listData.Created_At = time.Now()
	listData.Updated_At = time.Now()

	newList, err := s.repository.CreateList(ctx, listData)
	if err != nil {
		return nil, fmt.Errorf("failed to create list: %w", err)
	}

	return newList, nil
}

func (s *ListService) GetUserLists(ctx context.Context, userID uint) ([]models.List, error) {
	lists, err := s.repository.GetListsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user lists: %w", err)
	}

	return lists, nil
}

func (s *ListService) GetList(ctx context.Context, userID, listID uint) (*models.List, error) {
	list, err := s.repository.GetListByID(ctx, listID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("list not found")
		}
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	if list.UserID != userID {
		return nil, fmt.Errorf("unauthorized to access this list")
	}

	return list, nil
}

func (s *ListService) DeleteList(ctx context.Context, userID, listID uint) error {
	list, err := s.repository.GetListByID(ctx, listID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("list not found")
		}
		return fmt.Errorf("failed to get list: %w", err)
	}

	if list.UserID != userID {
		return fmt.Errorf("unauthorized to access this list")
	}

	if err := s.repository.DeleteList(ctx, listID); err != nil {
		return fmt.Errorf("unauthorized to access this list: %w", err)
	}

	return nil
}

func (s *ListService) UpdateList(ctx context.Context, userID, listID uint, updateData *models.List) (*models.List, error) {
	existingList, err := s.repository.GetListByID(ctx, listID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("list not found")
		}
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	if existingList.UserID != userID {
		return nil, fmt.Errorf("unauthorized to update this list")
	}

	if updateData.Name != "" {
		existingList.Name = updateData.Name
	}

	if updateData.Description != "" {
		existingList.Description = updateData.Description
	}

	existingList.Updated_At = time.Now()

	updateList, err := s.repository.UpdateList(ctx, listID, existingList)
	if err != nil {
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	return updateList, nil
}

func NewListService(repository models.ListRepository) models.ListService {
	return &ListService{
		repository: repository,
	}
}
