package repository

import (
	"context"
	"errors"
	"time"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/repository/models"
	"gorm.io/gorm"
)

type AdminsRepo struct {
	db *gorm.DB
}

func NewAdminsRepo(db *gorm.DB) *AdminsRepo {
	return &AdminsRepo{db: db}
}

func (r *AdminsRepo) GetByCredentials(ctx context.Context, schoolID uint, email, password string) (domain.Admin, error) {
	var model models.AdminModel
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND email = ? AND password = ?", schoolID, email, password).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Admin{}, domain.ErrUserNotFound
		}
		return domain.Admin{}, err
	}
	return model.ToDomain(), nil
}

func (r *AdminsRepo) GetByRefreshToken(ctx context.Context, schoolID uint, refreshToken string) (domain.Admin, error) {
	var model models.AdminModel
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND session_refresh_token = ? AND session_expires_at > ?", schoolID, refreshToken, time.Now()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Admin{}, domain.ErrUserNotFound
		}
		return domain.Admin{}, err
	}
	return model.ToDomain(), nil
}

func (r *AdminsRepo) SetSession(ctx context.Context, id uint, session domain.Session) error {
	return r.db.WithContext(ctx).
		Model(&models.AdminModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"session_refresh_token": session.RefreshToken,
			"session_expires_at":    time.Unix(session.ExpiresAt, 0),
		}).Error
}

func (r *AdminsRepo) GetById(ctx context.Context, id uint) (domain.Admin, error) {
	var model models.AdminModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		return domain.Admin{}, err
	}
	return model.ToDomain(), nil
}
