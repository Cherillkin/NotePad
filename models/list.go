package models

import (
	"context"
	"time"
)

type List struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:text;not null"`
	Description string    `json:"description"`
	UserID      uint      `json:"user_id"`
	User        User      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}

type ListRepository interface {
	CreateList(ctx context.Context, listData *List) (*List, error)
	GetListsByUserID(ctx context.Context, userID uint) ([]List, error)
	GetListByID(ctx context.Context, id uint) (*List, error)
	DeleteList(ctx context.Context, id uint) error
	UpdateList(ctx context.Context, id uint, listData *List) (*List, error)
}

type ListService interface {
	CreateList(ctx context.Context, userID uint, listData *List) (*List, error)
	GetUserLists(ctx context.Context, userID uint) ([]List, error)
	GetList(ctx context.Context, userID, listID uint) (*List, error)
	DeleteList(ctx context.Context, userID, listID uint) error
	UpdateList(ctx context.Context, userID, listID uint, updateData *List) (*List, error)
}
