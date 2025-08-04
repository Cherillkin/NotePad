package services

import (
	"context"

	"github.com/Cherillkin/Notepad/models"
)

type SharedListService struct {
	repository models.SharedListRepository
}

func (s *SharedListService) SharedList(ctx context.Context, listID, userID uint) error {
	return s.repository.SharedList(ctx, listID, userID)
}

func (s *SharedListService) GetSharedLists(ctx context.Context, userID uint) ([]models.List, error) {
	return s.repository.GetSharedLists(ctx, userID)
}

func NewSharedListService(repository models.SharedListRepository) models.SharedListService {
	return &SharedListService{
		repository: repository,
	}
}
