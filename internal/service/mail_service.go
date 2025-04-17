package service

import (
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
	"log"
	"time"
	"user-service/config"
	"user-service/internal/model"
	"user-service/internal/repository"
)

type MailService interface {
	GenerateEmailOTP(userID uint) (string, error)
	SendVerifyMail(email string, link string)
}

type mailService struct {
	repo repository.OTPRepository
}

func NewMailService(repo repository.OTPRepository) *mailService {
	return &mailService{repo}
}

func (s *mailService) GenerateEmailOTP(userID uint) (string, error) {
	token, err := uuid.NewUUID()
	if err != nil {
		return "", err // Proper error handling here
	}

	otp := model.EmailVerificationToken{
		UserID:    userID,
		Token:     token.String(),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := s.repo.Create(&otp); err != nil {
		return "", err // Return the error if the Create method fails
	}

	return token.String(), nil
}

func (s *mailService) SendVerifyMail(email string, link string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "dreamhoaihack@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Verify your email address")
	body := "Click the following link to verify your email: <a href=\"" + link + "\">Verify Email</a>"
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, config.GetString("MAIL_USERNAME"), config.GetString("MAIL_PASSWORD"))

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		// Handle the error (log or return)
		log.Println("Error sending email:", err)
	}
}
