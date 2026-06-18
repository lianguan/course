package service

import (
	"context"

	"ultrathreads/internal/domain"
)

type OffersService struct {
	repo            OffersRepository
	modulesService  Modules
	packagesService Packages
}

func NewOffersService(repo OffersRepository, modulesService Modules, packagesService Packages) *OffersService {
	return &OffersService{repo: repo, modulesService: modulesService, packagesService: packagesService}
}

func (s *OffersService) GetById(ctx context.Context, id uint) (domain.Offer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OffersService) getByPackage(ctx context.Context, schoolID, packageID uint) ([]domain.Offer, error) {
	offers, err := s.repo.GetBySchool(ctx, schoolID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Offer, 0)

	for i := range offers {
		if offers[i].ContainsPackage(packageID) {
			result = append(result, offers[i])
		}
	}

	return result, nil
}

func (s *OffersService) GetByModule(ctx context.Context, schoolID, moduleID uint) ([]domain.Offer, error) {
	module, err := s.modulesService.GetById(ctx, moduleID)
	if err != nil {
		return nil, err
	}

	return s.getByPackage(ctx, schoolID, module.PackageID)
}

func (s *OffersService) GetByCourse(ctx context.Context, courseID uint) ([]domain.Offer, error) {
	packages, err := s.packagesService.GetByCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}

	if len(packages) == 0 {
		return []domain.Offer{}, nil
	}

	packageIDs := make([]uint, len(packages))
	for i, pkg := range packages {
		packageIDs[i] = pkg.ID
	}

	return s.repo.GetByPackages(ctx, packageIDs)
}

func (s *OffersService) Create(ctx context.Context, inp domain.CreateOfferInput) (uint, error) {
	offer := domain.NewOffer(
		inp.SchoolID,
		inp.Name,
		inp.Description,
		inp.Benefits,
		inp.Price,
		inp.PaymentMethod,
		inp.Packages,
	)

	if err := offer.ValidatePaymentMethod(); err != nil {
		return 0, err
	}

	return s.repo.Create(ctx, *offer)
}

func (s *OffersService) GetAll(ctx context.Context, schoolID uint) ([]domain.Offer, error) {
	return s.repo.GetBySchool(ctx, schoolID)
}

func (s *OffersService) Update(ctx context.Context, inp domain.UpdateOfferInput) error {
	if err := inp.ValidatePayment(); err != nil {
		return err
	}

	return s.repo.Update(ctx, inp)
}

func (s *OffersService) Delete(ctx context.Context, schoolID, id uint) error {
	return s.repo.Delete(ctx, schoolID, id)
}

func (s *OffersService) GetByIds(ctx context.Context, ids []uint) ([]domain.Offer, error) {
	return s.repo.GetByIDs(ctx, ids)
}
