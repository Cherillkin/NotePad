package models

import (
	"context"
	"time"
)

type Item struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:text;not null"`
	Description string    `json:"description"`
	ListID      uint      `json:"list_id"`
	List        List      `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Created_At  time.Time `json:"created_at"`
	Updated_At  time.Time `json:"updated_at"`
}

type ItemRepository interface {
	CreateItem(ctx context.Context, itemData *Item) (*Item, error)
	GetListItems(ctx context.Context, listID uint) ([]Item, error)
	GetItem(ctx context.Context, id uint) (*Item, error)
	DeleteItem(ctx context.Context, id uint) error
	UpdateItem(ctx context.Context, id uint, listData *Item) (*Item, error)
}

type ItemService interface {
	CreateItem(ctx context.Context, listID uint, itemData *Item) (*Item, error)
	GetListItems(ctx context.Context, listID uint) ([]Item, error)
	GetItem(ctx context.Context, listID, itemID uint) (*Item, error)
	DeleteItem(ctx context.Context, listID, itemID uint) error
	UpdateItem(ctx context.Context, listID, itemID uint, updateData *Item) (*Item, error)
}
