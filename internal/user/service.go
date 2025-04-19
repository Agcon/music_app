package user

import (
	"context"
	"errors"
	"fmt"
	"music_app/pkg/auth"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, username, email, password string) error
	Login(ctx context.Context, email, password string) (int64, string, error)
	Logout(ctx context.Context, token string) error
}

type service struct {
	repo    Repository
	session auth.SessionManager
}

func NewService(repo Repository, session auth.SessionManager) Service {
	return &service{repo: repo, session: session}
}

func (s *service) Register(ctx context.Context, username, email, password string) error {
	existing, _ := s.repo.GetByEmail(ctx, email)
	if existing != nil {
		return fmt.Errorf("пользователь с таким email уже существует")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}

	return s.repo.Create(ctx, user)
}

func (s *service) Login(ctx context.Context, email, password string) (int64, string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return 0, "", errors.New("user not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return 0, "", errors.New("wrong password")
	}

	token, err := s.session.Create(ctx, user.ID, 7*24*time.Hour)
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}

func (s *service) Logout(ctx context.Context, token string) error {
	return s.session.Delete(ctx, token)
}
