package domain

import "time"

type Note struct {
	ID        uint      `gorm:"primary_key;autoIncrement"`
	UserID    uint      `gorm:"not null;index:idx_user_created"` // Index for faster queries by user
	Title     string    `gorm:"type:varchar(150);not null"`
	Content   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"index:idx_user_created"` // Composite index for user + created_at queries
	UpdatedAt time.Time
}
