package repository

import (
	"context"
	"errors"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/repository/models"

	"gorm.io/gorm"
)

type OffersRepo struct {
	db *gorm.DB
}

func NewOffersRepo(db *gorm.DB) *OffersRepo {
	return &OffersRepo{db: db}
}

func (r *OffersRepo) Create(ctx context.Context, offer domain.Offer) (uint, error) {
	var model models.OfferModel
	model.FromDomain(offer)
	err := r.db.WithContext(ctx).Create(&model).Error
	return model.ID, err
}

func (r *OffersRepo) GetBySchool(ctx context.Context, schoolID uint) ([]domain.Offer, error) {
	var offerModels []models.OfferModel
	err := r.db.WithContext(ctx).
		Where("school_id = ?", schoolID).
		Find(&offerModels).Error
	if err != nil {
		return nil, err
	}
	offers := make([]domain.Offer, len(offerModels))
	for i, m := range offerModels {
		offers[i] = m.ToDomain()
	}
	return offers, nil
}

func (r *OffersRepo) GetByID(ctx context.Context, id uint) (domain.Offer, error) {
	var model models.OfferModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Offer{}, domain.ErrOfferNotFound
		}
		return domain.Offer{}, err
	}
	return model.ToDomain(), nil
}

func (r *OffersRepo) GetByPackages(ctx context.Context, packageIDs []uint) ([]domain.Offer, error) {
	var offerModels []models.OfferModel
	// 使用 MySQL JSON 函数查询包含指定套餐ID的优惠
	err := r.db.WithContext(ctx).
		Where("JSON_CONTAINS(package_ids, ?)", packageIDs).
		Find(&offerModels).Error
	if err != nil {
		return nil, err
	}
	offers := make([]domain.Offer, len(offerModels))
	for i, m := range offerModels {
		offers[i] = m.ToDomain()
	}
	return offers, nil
}

func (r *OffersRepo) GetByIDs(ctx context.Context, ids []uint) ([]domain.Offer, error) {
	var offerModels []models.OfferModel
	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&offerModels).Error
	if err != nil {
		return nil, err
	}
	offers := make([]domain.Offer, len(offerModels))
	for i, m := range offerModels {
		offers[i] = m.ToDomain()
	}
	return offers, nil
}

func (r *OffersRepo) Update(ctx context.Context, inp domain.UpdateOfferInput) error {
	updates := map[string]interface{}{}

	if inp.Name != "" {
		updates["name"] = inp.Name
	}
	if inp.Description != "" {
		updates["description"] = inp.Description
	}
	if inp.Benefits != nil {
		updates["benefits"] = inp.Benefits
	}
	if inp.Price != nil {
		updates["price_value"] = inp.Price.Value
		updates["price_currency"] = inp.Price.Currency
	}
	if inp.Packages != nil {
		updates["package_ids"] = inp.Packages
	}
	if inp.PaymentMethod != nil {
		updates["payment_method_uses_provider"] = inp.PaymentMethod.UsesProvider
		updates["payment_method_provider"] = inp.PaymentMethod.Provider
	}

	return r.db.WithContext(ctx).
		Model(&models.OfferModel{}).
		Where("id = ? AND school_id = ?", inp.ID, inp.SchoolID).
		Updates(updates).Error
}

func (r *OffersRepo) Delete(ctx context.Context, schoolID, id uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&models.OfferModel{}).Error
}
