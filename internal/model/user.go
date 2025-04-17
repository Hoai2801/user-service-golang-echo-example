package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Name            string     `json:"name" gorm:"size:100;not null"`
	Email           string     `json:"email" gorm:"size:100;uniqueIndex;not null"`
	Password        string     `json:"-" gorm:"not null"`                      // hide password from JSON
	Role            string     `json:"role" gorm:"default:'user'"`             // "user", "admin", etc.
	Active          bool       `json:"active" gorm:"default:true"`             // account status
	LoginWithGoogle bool       `json:"login_with_google" gorm:"default:false"` // signed up via Google
	GoogleID        string     `json:"google_id,omitempty" gorm:"size:255"`    // optional, for Google OAuth linking
	EmailVerified   bool       `json:"email_verified" gorm:"default:false"`    // email confirmation status
	LastLogin       *time.Time `json:"last_login,omitempty"`                   // pointer allows null
	AvatarURL       string     `json:"avatar_url,omitempty" gorm:"size:255"`   // optional profile image

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // soft delete
}
