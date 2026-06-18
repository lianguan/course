package repository

import (
	"context"

	"ultrathreads/internal/domain"

	"gorm.io/gorm"
)

type ModulesRepo struct {
	db *gorm.DB
}

func NewModulesRepo(db *gorm.DB) *ModulesRepo {
	return &ModulesRepo{db: db}
}

func (r *ModulesRepo) Create(ctx context.Context, module domain.Module) (uint, error) {
	err := r.db.WithContext(ctx).Create(&module).Error
	return module.ID, err
}

func (r *ModulesRepo) GetPublishedByCourseID(ctx context.Context, courseID uint) ([]domain.Module, error) {
	var modules []domain.Module
	err := r.db.WithContext(ctx).
		Where("course_id = ? AND published = ?", courseID, true).
		Order("position ASC").
		Find(&modules).Error
	return modules, err
}

func (r *ModulesRepo) GetByCourseID(ctx context.Context, courseID uint) ([]domain.Module, error) {
	var modules []domain.Module
	err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("position ASC").
		Find(&modules).Error
	return modules, err
}

func (r *ModulesRepo) GetPublishedByID(ctx context.Context, moduleID uint) (domain.Module, error) {
	var module domain.Module
	err := r.db.WithContext(ctx).
		Where("id = ? AND published = ?", moduleID, true).
		First(&module).Error
	return module, err
}

func (r *ModulesRepo) GetByID(ctx context.Context, moduleID uint) (domain.Module, error) {
	var module domain.Module
	err := r.db.WithContext(ctx).Where("id = ?", moduleID).First(&module).Error
	return module, err
}

func (r *ModulesRepo) GetByPackages(ctx context.Context, packageIDs []uint) ([]domain.Module, error) {
	var modules []domain.Module
	err := r.db.WithContext(ctx).
		Where("package_id IN ?", packageIDs).
		Order("position ASC").
		Find(&modules).Error
	return modules, err
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
		Model(&domain.Module{}).
		Where("id = ? AND school_id = ?", inp.ID, inp.SchoolID).
		Updates(updates).Error
}

func (r *ModulesRepo) Delete(ctx context.Context, schoolID, id uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		Delete(&domain.Module{}).Error
}

func (r *ModulesRepo) AddLesson(ctx context.Context, schoolID, id uint, lesson domain.Lesson) error {
	var module domain.Module
	if err := r.db.WithContext(ctx).Where("id = ? AND school_id = ?", id, schoolID).First(&module).Error; err != nil {
		return err
	}

	module.Lessons = append(module.Lessons, lesson)
	return r.db.WithContext(ctx).Save(&module).Error
}

func (r *ModulesRepo) GetByLesson(ctx context.Context, lessonID uint) (domain.Module, error) {
	var module domain.Module
	// 使用 MySQL JSON 函数正确搜索 JSON 数组中的 ID
	err := r.db.WithContext(ctx).
		Where("JSON_CONTAINS(lessons, JSON_OBJECT('id', ?))", lessonID).
		First(&module).Error
	return module, err
}

func (r *ModulesRepo) UpdateLesson(ctx context.Context, inp domain.UpdateLessonInput) error {
	var module domain.Module
	if err := r.db.WithContext(ctx).
		Where("school_id = ?", inp.SchoolID).
		First(&module).Error; err != nil {
		return err
	}

	for i, lesson := range module.Lessons {
		if lesson.ID == inp.ID {
			if inp.Name != "" {
				module.Lessons[i].Name = inp.Name
			}
			if inp.Position != nil {
				module.Lessons[i].Position = *inp.Position
			}
			if inp.Published != nil {
				module.Lessons[i].Published = *inp.Published
			}
			break
		}
	}

	return r.db.WithContext(ctx).Save(&module).Error
}

func (r *ModulesRepo) DeleteLesson(ctx context.Context, schoolID, id uint) error {
	var module domain.Module
	if err := r.db.WithContext(ctx).
		Where("school_id = ?", schoolID).
		First(&module).Error; err != nil {
		return err
	}

	var filtered []domain.Lesson
	for _, lesson := range module.Lessons {
		if lesson.ID != id {
			filtered = append(filtered, lesson)
		}
	}
	module.Lessons = filtered

	return r.db.WithContext(ctx).Save(&module).Error
}

func (r *ModulesRepo) DetachPackageFromAll(ctx context.Context, schoolID, packageID uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Module{}).
		Where("school_id = ? AND package_id = ?", schoolID, packageID).
		Update("package_id", 0).Error
}

func (r *ModulesRepo) AttachPackage(ctx context.Context, schoolID, packageID uint, modules []uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Module{}).
		Where("id IN ? AND school_id = ?", modules, schoolID).
		Update("package_id", packageID).Error
}

func (r *ModulesRepo) DeleteByCourse(ctx context.Context, schoolID, courseID uint) error {
	return r.db.WithContext(ctx).
		Where("course_id = ? AND school_id = ?", courseID, schoolID).
		Delete(&domain.Module{}).Error
}

func (r *ModulesRepo) AttachSurvey(ctx context.Context, schoolID, id uint, survey domain.Survey) error {
	return r.db.WithContext(ctx).
		Model(&domain.Module{}).
		Where("id = ? AND school_id = ?", id, schoolID).
		Update("survey", survey).Error
}

func (r *ModulesRepo) DetachSurvey(ctx context.Context, schoolID, id uint) error {
	return r.db.WithContext(ctx).
		Model(&domain.Module{}).
		Where("id = ? AND school_id = ?", id, schoolID).
		Update("survey", nil).Error
}
