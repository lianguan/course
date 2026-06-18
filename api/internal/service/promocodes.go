package service

import (
	"context"
	"errors"

	"ultrathreads/internal/domain"
)

type PromoCodeService struct {
	repo PromoCodesRepository
}

func NewPromoCodeService(repo PromoCodesRepository) *PromoCodeService {
	return &PromoCodeService{repo: repo}
}

func (s *PromoCodeService) Create(ctx context.Context, inp domain.CreatePromoCodeInput) (uint, error) {
	return s.repo.Create(ctx, domain.PromoCode{
		SchoolID:           inp.SchoolID,
		Code:               inp.Code,
		DiscountPercentage: inp.DiscountPercentage,
		ExpiresAt:          inp.ExpiresAt,
		OfferIDs:           inp.OfferIDs,
	})
}

func (s *PromoCodeService) Update(ctx context.Context, inp domain.UpdatePromoCodeInput) error {
	return s.repo.Update(ctx, inp)
}

func (s *PromoCodeService) Delete(ctx context.Context, schoolID, id uint) error {
	return s.repo.Delete(ctx, schoolID, id)
}

func (s *PromoCodeService) GetByCode(ctx context.Context, schoolID uint, code string) (domain.PromoCode, error) {
	promo, err := s.repo.GetByCode(ctx, schoolID, code)
	if err != nil {
		if errors.Is(err, domain.ErrPromoNotFound) {
			return domain.PromoCode{}, err
		}

		return domain.PromoCode{}, err
	}

	return promo, nil
}

func (s *PromoCodeService) GetById(ctx context.Context, schoolID, id uint) (domain.PromoCode, error) {
	promo, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		if errors.Is(err, domain.ErrPromoNotFound) {
			return domain.PromoCode{}, err
		}

		return domain.PromoCode{}, err
	}

	return promo, nil
}

func (s *PromoCodeService) GetBySchool(ctx context.Context, schoolID uint) ([]domain.PromoCode, error) {
	return s.repo.GetBySchool(ctx, schoolID)
}
