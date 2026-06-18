package repository

import (
	"context"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/repository/models"
	"gorm.io/gorm"
)

type PackagesRepo struct {
	db *gorm.DB
}

func NewPackagesRepo(db *gorm.DB) *PackagesRepo {
	return &PackagesRepo{db: db}
}

func (r *PackagesRepo) Create(ctx context.Context, pkg domain.Package) (uint, error) {
	var model models.PackageModel
	model.FromDomain(pkg)
	err := r.db.WithContext(ctx).Create(&model).Error
	return model.ID, err
}

func (r *PackagesRepo) GetByCourse(ctx context.Context, courseID uint) ([]domain.Package, error) {
	var pkgModels []models.PackageModel
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Find(&pkgModels).Error
	if err != nil {
		return nil, err
	}
	packages := make([]domain.Package, len(pkgModels))
	for i, m := range pkgModels {
		packages[i] = m.ToDomain()
	}
	return packages, nil
}

func (r *PackagesRepo) GetByID(ctx context.Context, id uint) (domain.Package, error) {
	var model models.PackageModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		return domain.Package{}, err
	}
	return model.ToDomain(), nil
}

func (r *PackagesRepo) GetByIDs(ctx context.Context, ids []uint) ([]domain.Package, error) {
	var pkgModels []models.PackageModel
	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&pkgModels).Error
	if err != nil {
		return nil, err
	}
	packages := make([]domain.Package, len(pkgModels))
	for i, m := range pkgModels {
		packages[i] = m.ToDomain()
	}
	return packages, nil
}

func (r *PackagesRepo) Update(ctx context.Context, inp domain.UpdatePackageInput) error {
	updates := map[string]interface{}{}

	if inp.Name != "" {
		updates["name"] = inp.Name
	}

	return r.db.WithContext(ctx).
		Model(&models.PackageModel{}).
		Where("id = ? AND school_id = ?", inp.ID, inp.SchoolID).
		Updates(updates).Error
}

func (r *PackagesRepo) Delete(ctx context.Context, schoolID, id uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&models.PackageModel{}).Error
}
