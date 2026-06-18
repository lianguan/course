package service

import (
	"context"

	"ultrathreads/internal/domain"
	"ultrathreads/pkg/cache"
)

type SchoolsService struct {
	repo  SchoolsRepository
	cache cache.Cache
	ttl   int64
}

func NewSchoolsService(repo SchoolsRepository, cache cache.Cache, ttl int64) *SchoolsService {
	return &SchoolsService{repo: repo, cache: cache, ttl: ttl}
}

func (s *SchoolsService) Create(ctx context.Context, name string) (uint, error) {
	return s.repo.Create(ctx, name)
}

func (s *SchoolsService) GetByDomain(ctx context.Context, domainName string) (domain.School, error) {
	if value, err := s.cache.Get(domainName); err == nil {
		return value.(domain.School), nil
	}

	school, err := s.repo.GetByDomain(ctx, domainName)
	if err != nil {
		return domain.School{}, err
	}

	err = s.cache.Set(domainName, school, s.ttl)

	return school, err
}

func (s *SchoolsService) GetById(ctx context.Context, id uint) (domain.School, error) {
	return s.repo.GetById(ctx, id)
}

func (s *SchoolsService) UpdateSettings(ctx context.Context, schoolID uint, inp domain.UpdateSchoolSettingsInput) error {
	return s.repo.UpdateSettings(ctx, schoolID, inp)
}

func (s *SchoolsService) ConnectFondy(ctx context.Context, input domain.ConnectFondyInput) error {
	// TODO: implement fondy connection test
	creds := domain.Fondy{
		MerchantPassword: input.MerchantPassword,
		MerchantID:       input.MerchantID,
		Connected:        true,
	}

	return s.repo.SetFondyCredentials(ctx, input.SchoolID, creds)
}

func (s *SchoolsService) ConnectSendPulse(ctx context.Context, input domain.ConnectSendPulseInput) error {
	// todo
	return nil
}
