package models

import (
	"time"

	"gorm.io/gorm"
)

type EntityUrl struct {
	gorm.Model
	ID            uint      `gorm:"primaryKey"`
	Url           string    `gorm:"not null;uniqueIndex"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	BaseEntityUrl string    `gorm:"type:varchar(100);not null"`
}

func NewEntityUrl(url string, baseEntityUrl string) *EntityUrl {
	return &EntityUrl{
		Url:           url,
		BaseEntityUrl: baseEntityUrl,
	}
}
