package repository

import (
	"context"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/repository/models"

	"gorm.io/gorm"
)

type ModulesRepo struct {
	db *gorm.DB
}

func NewModulesRepo(db *gorm.DB) *ModulesRepo {
	return &ModulesRepo{db: db}
}

func (r *ModulesRepo) Create(ctx context.Context, module domain.Module) (uint, error) {
	var model models.ModuleModel
	model.FromDomain(module)
	err := r.db.WithContext(ctx).Create(&model).Error
	return model.ID, err
}

func (r *ModulesRepo) GetPublishedByCourseID(ctx context.Context, courseID uint) ([]domain.Module, error) {
	var moduleModels []models.ModuleModel
	err := r.db.WithContext(ctx).
		Where("course_id = ? AND published = ?", courseID, true).
		Order("position ASC").
		Find(&moduleModels).Error
	if err != nil {
		return nil, err
	}
	modules := make([]domain.Module, len(moduleModels))
	for i, m := range moduleModels {
		modules[i] = m.ToDomain()
	}
	return modules, nil
}

func (r *ModulesRepo) GetByCourseID(ctx context.Context, courseID uint) ([]domain.Module, error) {
	var moduleModels []models.ModuleModel
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("position ASC").
		Find(&moduleModels).Error
	if err != nil {
		return nil, err
	}
	modules := make([]domain.Module, len(moduleModels))
	for i, m := range moduleModels {
		modules[i] = m.ToDomain()
	}
	return modules, nil
}

func (r *ModulesRepo) GetPublishedByID(ctx context.Context, moduleID uint) (domain.Module, error) {
	var model models.ModuleModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND published = ?", moduleID, true).
		First(&model).Error
	if err != nil {
		return domain.Module{}, err
	}
	return model.ToDomain(), nil
}

func (r *ModulesRepo) GetByID(ctx context.Context, moduleID uint) (domain.Module, error) {
	var model models.ModuleModel
	err := r.db.WithContext(ctx).Where("id = ?", moduleID).First(&model).Error
	if err != nil {
		return domain.Module{}, err
	}
	return model.ToDomain(), nil
}

func (r *ModulesRepo) GetByPackages(ctx context.Context, packageIDs []uint) ([]domain.Module, error) {
	var moduleModels []models.ModuleModel
	err := r.db.WithContext(ctx).
		Where("package_id IN ?", packageIDs).
		Order("position ASC").
		Find(&moduleModels).Error
	if err != nil {
		return nil, err
	}
	modules := make([]domain.Module, len(moduleModels))
	for i, m := range moduleModels {
		modules[i] = m.ToDomain()
	}
	return modules, nil
}

func (r *ModulesRepo) Update(ctx context.Context, inp domain.UpdateModuleInput) error {
	updates := map[string]interface{}{}

	if inp.Name != "" {
		updates["name"] = inp.Name
	}
	if inp.Position != nil {
		updates["position"] = *inp.Position
	}
	if inp.Published != nil {
		updates["published"] = *inp.Published
	}

	return r.db.WithContext(ctx).
		Model(&models.ModuleModel{}).
		Where("id = ? AND school_id = ?", inp.ID, inp.SchoolID).
		Updates(updates).Error
}

func (r *ModulesRepo) Delete(ctx context.Context, schoolID, id uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&models.ModuleModel{}).Error
}

func (r *ModulesRepo) DetachPackageFromAll(ctx context.Context, schoolID, packageID uint) error {
	return r.db.WithContext(ctx).
		Model(&models.ModuleModel{}).
		Where("school_id = ? AND package_id = ?", schoolID, packageID).
		Update("package_id", 0).Error
}

func (r *ModulesRepo) AttachPackage(ctx context.Context, schoolID, packageID uint, modules []uint) error {
	return r.db.WithContext(ctx).
		Model(&models.ModuleModel{}).
		Where("id IN ? AND school_id = ?", modules, schoolID).
		Update("package_id", packageID).Error
}

func (r *ModulesRepo) DeleteByCourse(ctx context.Context, schoolID, courseID uint) error {
	return r.db.WithContext(ctx).
		Where("course_id = ? AND school_id = ?", courseID, schoolID).
		Delete(&models.ModuleModel{}).Error
}
