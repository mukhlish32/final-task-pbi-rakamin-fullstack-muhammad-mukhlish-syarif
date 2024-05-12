package models

import (
	"time"
)

type Photo struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Caption   string
	PhotoURL  string
	UserID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
