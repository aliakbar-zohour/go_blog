// service/auth_service: Registration (email code), verification, and login with JWT.
package service

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/aliakbar-zohour/go_blog/internal/config"
	"github.com/aliakbar-zohour/go_blog/internal/mail"
	"github.com/aliakbar-zohour/go_blog/internal/model"
	"github.com/aliakbar-zohour/go_blog/internal/repository"
	"github.com/aliakbar-zohour/go_blog/pkg/auth"
	"gorm.io/gorm"
)

const codeLength = 6
const codeExpiryMinutes = 15

type AuthService struct {
	authorRepo   *repository.AuthorRepository
	evRepo       *repository.EmailVerificationRepository
	cfg          *config.Config
}

func NewAuthService(authorRepo *repository.AuthorRepository, evRepo *repository.EmailVerificationRepository, cfg *config.Config) *AuthService {
	return &AuthService{authorRepo: authorRepo, evRepo: evRepo, cfg: cfg}
}

// RequestVerification sends a 6-digit code to the email. Returns the code when SMTP is not configured (for dev/testing so you can use it in Swagger).
func (s *AuthService) RequestVerification(ctx context.Context, email string) (devCode string, err error) {
	email = normalizeEmail(email)
	if email == "" {
		return "", errors.New("email is required")
	}
	code, err := generateCode(codeLength)
	if err != nil {
		return "", err
	}
	_ = s.evRepo.DeleteByEmail(ctx, email)
	ev := &model.EmailVerification{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(codeExpiryMinutes * time.Minute),
	}
	if err := s.evRepo.Create(ctx, ev); err != nil {
		return "", err
	}
	if s.cfg.SMTPHost == "" {
		log.Printf("[auth] SMTP not configured; verification code for %s: %s", email, code)
		return code, nil
	}
	_ = mail.SendVerificationCode(email, code, s.cfg.SMTPHost, s.cfg.SMTPPort, s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPFrom)
	return "", nil
}

// VerifyAndRegister checks the code, creates the author with name/password, and returns author + JWT.
func (s *AuthService) VerifyAndRegister(ctx context.Context, email, code, name, password string) (*model.Author, string, error) {
	email = normalizeEmail(email)
	name = strings.TrimSpace(name)
	if email == "" || code == "" || name == "" || password == "" {
		return nil, "", errors.New("email, code, name and password are required")
	}
	if len(password) < 8 {
		return nil, "", errors.New("password must be at least 8 characters")
	}
	_, err := s.evRepo.FindValid(ctx, email, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid or expired code")
		}
		return nil, "", err
	}
	_ = s.evRepo.DeleteByEmail(ctx, email)
	existing, _ := s.authorRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, "", errors.New("email already registered")
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, "", err
	}
	now := time.Now()
	a := &model.Author{
		Name:            name,
		Email:           &email,
		PasswordHash:    hash,
		EmailVerifiedAt: &now,
	}
	if err := s.authorRepo.Create(ctx, a); err != nil {
		return nil, "", err
	}
	token, err := auth.NewToken(a.ID, s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
	if err != nil {
		return a, "", err
	}
	return a, token, nil
}

// Login returns author and JWT if email/password are valid.
func (s *AuthService) Login(ctx context.Context, email, password string) (*model.Author, string, error) {
	email = normalizeEmail(email)
	if email == "" || password == "" {
		return nil, "", errors.New("email and password are required")
	}
	a, err := s.authorRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("invalid email or password")
		}
		return nil, "", err
	}
	if a.PasswordHash == "" {
		return nil, "", errors.New("invalid email or password")
	}
	if !auth.CheckPassword(a.PasswordHash, password) {
		return nil, "", errors.New("invalid email or password")
	}
	token, err := auth.NewToken(a.ID, s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
	if err != nil {
		return nil, "", err
	}
	return a, token, nil
}

func normalizeEmail(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func generateCode(length int) (string, error) {
	const digits = "0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		b[i] = digits[n.Int64()]
	}
	return string(b), nil
}
