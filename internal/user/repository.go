package user

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, userID int64) (*User, error)
	GetAll(ctx context.Context) ([]User, error)
	Delete(ctx context.Context, userID int64) error
	UpdateRole(ctx context.Context, userID int64, role string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByID(ctx context.Context, userID int64) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Delete(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Delete(&User{}, userID).Error
}

func (r *repository) GetAll(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	return users, nil
}

func (r *repository) UpdateRole(ctx context.Context, userID int64, role string) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", userID).Update("role", role).Error
}
