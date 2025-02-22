package models

import (
	"time"

	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey"`
	Url       string    `gorm:"not null;uniqueIndex"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
