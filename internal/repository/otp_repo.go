package repository

import (
	"gorm.io/gorm"
	"user-service/internal/model"
)

type OTPRepository struct {
	db *gorm.DB
}

type OtpRepository interface {
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return OTPRepository{db}
}

func (r *OTPRepository) Create(otp *model.EmailVerificationToken) error {
	return r.db.Create(otp).Error
}
