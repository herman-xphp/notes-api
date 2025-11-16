package domain

import "time"

type Note struct {
	ID        uint   `gorm:"primary_key;autoIncrement"`
	UserID    uint   `gorm:"not null"`
	Title     string `gorm:"type:varchar(150);not null"`
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
