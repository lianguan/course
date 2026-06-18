package service

import (
	"context"

	"ultrathreads/internal/domain"
)

type PackagesService struct {
	repo        PackagesRepository
	modulesRepo ModulesRepository
}

func NewPackagesService(repo PackagesRepository, modulesRepo ModulesRepository) *PackagesService {
	return &PackagesService{repo: repo, modulesRepo: modulesRepo}
}

func (s *PackagesService) Create(ctx context.Context, inp domain.CreatePackageInput) (uint, error) {
	id, err := s.repo.Create(ctx, domain.Package{
		CourseID: inp.CourseID,
		SchoolID: inp.SchoolID,
		Name:     inp.Name,
	})

	if err != nil {
		return 0, err
	}

	if inp.Modules != nil {
		if err := s.modulesRepo.AttachPackage(ctx, inp.SchoolID, id, inp.Modules); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (s *PackagesService) GetByCourse(ctx context.Context, courseID uint) ([]domain.Package, error) {
	pkgs, err := s.repo.GetByCourse(ctx, courseID)
	if err != nil {
		return nil, err
	}

	for i := range pkgs {
		modules, err := s.modulesRepo.GetByPackages(ctx, []uint{pkgs[i].ID})
		if err != nil {
			return nil, err
		}

		pkgs[i].Modules = modules
	}

	return pkgs, nil
}

func (s *PackagesService) GetById(ctx context.Context, id uint) (domain.Package, error) {
	pkg, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return pkg, err
	}

	modules, err := s.modulesRepo.GetByPackages(ctx, []uint{pkg.ID})
	if err != nil {
		return pkg, err
	}

	pkg.Modules = modules

	return pkg, nil
}

func (s *PackagesService) GetByIds(ctx context.Context, ids []uint) ([]domain.Package, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	return s.repo.GetByIDs(ctx, ids)
}

func (s *PackagesService) Update(ctx context.Context, inp domain.UpdatePackageInput) error {
	if inp.Name != "" {
		if err := s.repo.Update(ctx, inp); err != nil {
			return err
		}
	}

	if inp.Modules != nil {
		if err := s.modulesRepo.DetachPackageFromAll(ctx, inp.SchoolID, inp.ID); err != nil {
			return err
		}

		if err := s.modulesRepo.AttachPackage(ctx, inp.SchoolID, inp.ID, inp.Modules); err != nil {
			return err
		}
	}

	return nil
}

func (s *PackagesService) Delete(ctx context.Context, schoolID, id uint) error {
	return s.repo.Delete(ctx, schoolID, id)
}
