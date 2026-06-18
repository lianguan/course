package repository

import (
	"context"
	"errors"
	"time"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/repository/models"

	"gorm.io/gorm"
)

type UsersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) *UsersRepo {
	return &UsersRepo{db: db}
}

func (r *UsersRepo) Create(ctx context.Context, user domain.User) error {
	var model models.UserModel
	model.FromDomain(user)
	res := r.db.WithContext(ctx).Create(&model)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *UsersRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var model models.UserModel
	err := r.db.WithContext(ctx).Where("email = ? AND password = ?", email, password).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return model.ToDomain(), nil
}

func (r *UsersRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var model models.UserModel
	err := r.db.WithContext(ctx).
		Where("session_refresh_token = ? AND session_expires_at > ?", refreshToken, time.Now()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return model.ToDomain(), nil
}

func (r *UsersRepo) Verify(ctx context.Context, userID uint, code string) error {
	res := r.db.WithContext(ctx).
		Model(&models.UserModel{}).
		Where("id = ? AND verification_code = ?", userID, code).
		Updates(map[string]interface{}{
			"verification_verified": true,
			"verification_code":     "",
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrVerificationCodeInvalid
	}
	return nil
}

func (r *UsersRepo) SetSession(ctx context.Context, userID uint, session domain.Session) error {
	return r.db.WithContext(ctx).
		Model(&models.UserModel{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"session_refresh_token": session.RefreshToken,
			"session_expires_at":    time.Unix(session.ExpiresAt, 0),
			"last_visit_at":         time.Now(),
		}).Error
}

func (r *UsersRepo) AttachSchool(ctx context.Context, userID, schoolID uint) error {
	var model models.UserModel
	if err := r.db.WithContext(ctx).First(&model, userID).Error; err != nil {
		return err
	}

	schools := model.Schools
	schools = append(schools, schoolID)

	return r.db.WithContext(ctx).
		Model(&models.UserModel{}).
		Where("id = ?", userID).
		Update("schools", schools).Error
}
