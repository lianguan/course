package service

import (
	"context"

	"ultrathreads/internal/domain"
)

type OrdersService struct {
	offersService     Offers
	promoCodesService PromoCodes
	studentsService   Students

	repo OrdersRepository
}

func NewOrdersService(repo OrdersRepository, offersService Offers, promoCodesService PromoCodes, studentsService Students) *OrdersService {
	return &OrdersService{
		repo:              repo,
		offersService:     offersService,
		promoCodesService: promoCodesService,
		studentsService:   studentsService,
	}
}

func (s *OrdersService) Create(ctx context.Context, studentID, offerID, promocodeID uint) (uint, error) {
	offer, err := s.offersService.GetById(ctx, offerID)
	if err != nil {
		return 0, err
	}

	promocode, err := s.getOrderPromocode(ctx, offer.SchoolID, promocodeID)
	if err != nil {
		return 0, err
	}

	student, err := s.studentsService.GetById(ctx, offer.SchoolID, studentID)
	if err != nil {
		return 0, err
	}

	orderAmount := domain.CalculateDiscountedPrice(offer.Price.Value, promocode.DiscountPercentage)

	promoInfo := domain.OrderPromoInfo{}
	if promocode.ID != 0 {
		promoInfo = domain.OrderPromoInfo{
			ID:   promocode.ID,
			Code: promocode.Code,
		}
	}

	order := domain.NewOrder(
		offer.SchoolID,
		domain.StudentInfoShort{
			ID:    student.ID,
			Name:  student.Name,
			Email: student.Email,
		},
		domain.OrderOfferInfo{
			ID:   offer.ID,
			Name: offer.Name,
		},
		orderAmount,
		offer.Price.Currency,
		promoInfo,
	)

	if err := s.repo.Create(ctx, *order); err != nil {
		return 0, err
	}

	return order.ID, nil
}

func (s *OrdersService) AddTransaction(ctx context.Context, id uint, transaction domain.Transaction) (domain.Order, error) {
	return s.repo.AddTransaction(ctx, id, transaction)
}

func (s *OrdersService) GetBySchool(ctx context.Context, schoolID uint, query domain.GetOrdersQuery) ([]domain.Order, int64, error) {
	return s.repo.GetBySchool(ctx, schoolID, query)
}

func (s *OrdersService) GetById(ctx context.Context, id uint) (domain.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrdersService) SetStatus(ctx context.Context, id uint, status string) error {
	return s.repo.SetStatus(ctx, id, status)
}

func (s *OrdersService) getOrderPromocode(ctx context.Context, schoolID, promocodeID uint) (domain.PromoCode, error) {
	var (
		promocode domain.PromoCode
		err       error
	)

	if promocodeID != 0 {
		promocode, err = s.promoCodesService.GetById(ctx, schoolID, promocodeID)
		if err != nil {
			return promocode, err
		}

		if promocode.IsExpired() {
			return promocode, domain.ErrPromocodeExpired
		}
	}

	return promocode, nil
}
