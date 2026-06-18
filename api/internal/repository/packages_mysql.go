package repository

import (
	"context"

	"ultrathreads/internal/domain"
	"gorm.io/gorm"
)

type PackagesRepo struct {
	db *gorm.DB
}

func NewPackagesRepo(db *gorm.DB) *PackagesRepo {
	return &PackagesRepo{db: db}
}

func (r *PackagesRepo) Create(ctx context.Context, pkg domain.Package) (uint, error) {
	err := r.db.WithContext(ctx).Create(&pkg).Error
	return pkg.ID, err
}

func (r *PackagesRepo) GetByCourse(ctx context.Context, courseID uint) ([]domain.Package, error) {
	var packages []domain.Package
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Find(&packages).Error
	return packages, err
}

func (r *PackagesRepo) GetByID(ctx context.Context, id uint) (domain.Package, error) {
	var pkg domain.Package
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&pkg).Error
	return pkg, err
}

func (r *PackagesRepo) GetByIDs(ctx context.Context, ids []uint) ([]domain.Package, error) {
	var packages []domain.Package
	err := r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&packages).Error
	return packages, err
}

func (r *PackagesRepo) Update(ctx context.Context, inp domain.UpdatePackageInput) error {
	updates := map[string]interface{}{}

	if inp.Name != "" {
		updates["name"] = inp.Name
	}

	return r.db.WithContext(ctx).
		Model(&domain.Package{}).
		Where("id = ? AND school_id = ?", inp.ID, inp.SchoolID).
		Updates(updates).Error
}

func (r *PackagesRepo) Delete(ctx context.Context, schoolID, id uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&domain.Package{}).Error
}
