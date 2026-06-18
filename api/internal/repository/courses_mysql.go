package repository

import (
	"context"
	"time"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/repository/models"
	"gorm.io/gorm"
)

type CoursesRepo struct {
	db *gorm.DB
}

func NewCoursesRepo(db *gorm.DB) *CoursesRepo {
	return &CoursesRepo{db: db}
}

func (r *CoursesRepo) Create(ctx context.Context, schoolID uint, course domain.Course) (uint, error) {
	var model models.CourseModel
	model.FromDomain(course)
	err := r.db.WithContext(ctx).Create(&model).Error
	return model.ID, err
}

func (r *CoursesRepo) Update(ctx context.Context, inp domain.UpdateCourseInput) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if inp.Name != nil {
		updates["name"] = *inp.Name
	}
	if inp.Description != nil {
		updates["description"] = *inp.Description
	}
	if inp.ImageURL != nil {
		updates["image_url"] = *inp.ImageURL
	}
	if inp.Color != nil {
		updates["color"] = *inp.Color
	}
	if inp.Published != nil {
		updates["published"] = *inp.Published
	}

	return r.db.WithContext(ctx).
		Model(&models.CourseModel{}).
		Where("id = ?", inp.ID).
		Updates(updates).Error
}

func (r *CoursesRepo) Delete(ctx context.Context, schoolID, courseID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ?", courseID).
		Delete(&models.CourseModel{}).Error
}
