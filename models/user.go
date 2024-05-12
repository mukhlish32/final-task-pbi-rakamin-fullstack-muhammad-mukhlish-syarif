package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"not null" json:"username" valid:"required"`
	Email     string    `gorm:"unique;not null" json:"email" valid:"email,required"`
	Password  string    `gorm:"not null" json:"password" valid:"required,length(6|255)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Photos    []Photo   `gorm:"foreignKey:UserID" json:"photos"`
}
