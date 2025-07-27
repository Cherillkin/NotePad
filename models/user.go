package models

import "time"

type User struct {
	ID         uint      `json:"id" gorm:"primary key"`
	Email      string    `json:"email" gorm:"text;not null"`
	Password   string    `json:"-"`
	Picture    string    `json:"picture"`
	Lists      []List    `json:"foreignKey:UserId"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
}
