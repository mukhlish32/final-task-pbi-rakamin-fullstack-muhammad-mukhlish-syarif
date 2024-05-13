package models

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID        uint      `gorm:"primaryKey" form:"id"`
	Title     string    `gorm:"not null" form:"title"`
	Caption   string    `gorm:"not null" form:"caption"`
	PhotoURL  string    `gorm:"not null" form:"photo_url"`
	UserID    uint      `gorm:"not null" form:"user_id"`
	CreatedAt time.Time `form:"created_at"`
	UpdatedAt time.Time `form:"updated_at"`
}

func (photo *Photo) BeforeCreate(tx *gorm.DB) (err error) {
	photo.CreatedAt = time.Now()
	photo.UpdatedAt = time.Now()
	return
}

func (photo *Photo) BeforeUpdate(tx *gorm.DB) (err error) {
	photo.UpdatedAt = time.Now()
	return
}
