package model

import "time"

type EmailVerificationToken struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index"`
	Token     string `gorm:"uniqueIndex not null;size:255;type:varchar(255)"`
	ExpiresAt time.Time
	CreatedAt time.Time
}
