package models

import (
	"context"
)

type SharedList struct {
	ID     uint `gorm:"primarKey"`
	ListID uint `gorm:"not null;index"`
	UserID uint `gorm:"not null;index"`
}

type SharedListRepository interface {
	SharedList(ctx context.Context, listID, userID uint) error
	GetSharedLists(ctx context.Context, userID uint) ([]List, error)
}

type SharedListService interface {
	SharedList(ctx context.Context, listID, userID uint) error
	GetSharedLists(ctx context.Context, userID uint) ([]List, error)
}
