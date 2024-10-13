package models

import "time"

type Event struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedBy   string    `json:"created_by" gorm:"not null"`
	User        User      `json:"-" gorm:"foreignKey:CreatedBy"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"   gorm:"default:null"`
	EndTime     time.Time `json:"end_time"   gorm:"default:null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoUpdateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
