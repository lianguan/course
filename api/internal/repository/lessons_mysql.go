package repository

import (
	"context"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/repository/models"
	"gorm.io/gorm"
)

type LessonContentRepo struct {
	db *gorm.DB
}

func NewLessonContentRepo(db *gorm.DB) *LessonContentRepo {
	return &LessonContentRepo{db: db}
}

func (r *LessonContentRepo) GetByLessons(ctx context.Context, lessonIDs []uint) ([]domain.LessonContent, error) {
	var contentModels []models.LessonContentModel
	err := r.db.WithContext(ctx).
		Where("lesson_id IN ?", lessonIDs).
		Find(&contentModels).Error
	if err != nil {
		return nil, err
	}
	contents := make([]domain.LessonContent, len(contentModels))
	for i, m := range contentModels {
		contents[i] = m.ToDomain()
	}
	return contents, nil
}

func (r *LessonContentRepo) GetByLesson(ctx context.Context, lessonID uint) (domain.LessonContent, error) {
	var model models.LessonContentModel
	err := r.db.WithContext(ctx).
		Where("lesson_id = ?", lessonID).
		First(&model).Error
	if err != nil {
		return domain.LessonContent{}, err
	}
	return model.ToDomain(), nil
}

func (r *LessonContentRepo) Update(ctx context.Context, schoolID, lessonID uint, content string) error {
	var model models.LessonContentModel
	err := r.db.WithContext(ctx).
		Where("lesson_id = ? AND school_id = ?", lessonID, schoolID).
		First(&model).Error

	if err == gorm.ErrRecordNotFound {
		// Create new
		model = models.LessonContentModel{
			LessonID: lessonID,
			SchoolID: schoolID,
			Content:  content,
		}
		return r.db.WithContext(ctx).Create(&model).Error
	} else if err != nil {
		return err
	}

	// Update existing
	model.Content = content
	return r.db.WithContext(ctx).Save(&model).Error
}

func (r *LessonContentRepo) DeleteContent(ctx context.Context, schoolID uint, lessonIDs []uint) error {
	return r.db.WithContext(ctx).
		Where("school_id = ? AND lesson_id IN ?", schoolID, lessonIDs).
		Delete(&models.LessonContentModel{}).Error
}
