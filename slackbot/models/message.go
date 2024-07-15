package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        uint   `gorm:"primary_key"`
	Text      string `gorm:"type:text"`
	CreatedAt time.Time
}

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&Message{})
}
