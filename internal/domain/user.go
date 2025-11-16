package domain

import "time"

type User struct {
	ID        uint      `gorm:"primary_key;column:id;autoIncrement"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"`
	Email     string    `gorm:"column:email;type:varchar(100);uniqueIndex;not null"`
	Password  string    `gorm:"column:password;type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`

	Notes []Note `gorm:"foreignKey:UserID"`
}
