package repository

import (
	"context"
	"errors"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/repository/models"
	"gorm.io/gorm"
)

type PromocodesRepo struct {
	db *gorm.DB
}

func NewPromocodeRepo(db *gorm.DB) *PromocodesRepo {
	return &PromocodesRepo{db: db}
}

func (r *PromocodesRepo) Create(ctx context.Context, promocode domain.PromoCode) (uint, error) {
	var model models.PromoCodeModel
	model.FromDomain(promocode)
	err := r.db.WithContext(ctx).Create(&model).Error
	return model.ID, err
}

func (r *PromocodesRepo) Update(ctx context.Context, inp domain.UpdatePromoCodeInput) error {
	updates := map[string]interface{}{}

	if inp.Code != "" {
		updates["code"] = inp.Code
	}
	if inp.DiscountPercentage != 0 {
		updates["discount_percentage"] = inp.DiscountPercentage
	}
	if inp.ExpiresAt != 0 {
		updates["expires_at"] = inp.ExpiresAt
	}
	if inp.OfferIDs != nil {
		updates["offer_ids"] = inp.OfferIDs
	}

	return r.db.WithContext(ctx).
		Model(&models.PromoCodeModel{}).
		Where("id = ? AND school_id = ?", inp.ID, inp.SchoolID).
		Updates(updates).Error
}

func (r *PromocodesRepo) Delete(ctx context.Context, schoolID, id uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&models.PromoCodeModel{}).Error
}

func (r *PromocodesRepo) GetByCode(ctx context.Context, schoolID uint, code string) (domain.PromoCode, error) {
	var model models.PromoCodeModel
	err := r.db.WithContext(ctx).
		Where("school_id = ? AND code = ?", schoolID, code).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.PromoCode{}, domain.ErrPromoNotFound
		}
		return domain.PromoCode{}, err
	}
	return model.ToDomain(), nil
}

func (r *PromocodesRepo) GetByID(ctx context.Context, schoolID, id uint) (domain.PromoCode, error) {
	var model models.PromoCodeModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.PromoCode{}, domain.ErrPromoNotFound
		}
		return domain.PromoCode{}, err
	}
	return model.ToDomain(), nil
}

func (r *PromocodesRepo) GetBySchool(ctx context.Context, schoolID uint) ([]domain.PromoCode, error) {
	var promoModels []models.PromoCodeModel
	err := r.db.WithContext(ctx).
		Where("school_id = ?", schoolID).
		Find(&promoModels).Error
	if err != nil {
		return nil, err
	}
	promocodes := make([]domain.PromoCode, len(promoModels))
	for i, m := range promoModels {
		promocodes[i] = m.ToDomain()
	}
	return promocodes, nil
}
