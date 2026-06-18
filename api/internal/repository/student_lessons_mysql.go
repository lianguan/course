package repository

import (
	"context"

	"ultrathreads/internal/domain"
	"gorm.io/gorm"
)

type StudentLessonsRepo struct {
	db *gorm.DB
}

func NewStudentLessonsRepo(db *gorm.DB) *StudentLessonsRepo {
	return &StudentLessonsRepo{db: db}
}

func (r *StudentLessonsRepo) AddFinished(ctx context.Context, studentID, lessonID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var studentLessons domain.StudentLessons
		err := tx.Set("gorm:query_option", "FOR UPDATE").First(&studentLessons, "student_id = ?", studentID).Error
		
		if err == gorm.ErrRecordNotFound {
			// 记录不存在，创建新记录
			studentLessons = domain.StudentLessons{
				StudentID: studentID,
				Finished:  []uint{lessonID},
			}
			return tx.Create(&studentLessons).Error
		} else if err != nil {
			return err
		}

		// 检查是否已存在
		for _, id := range studentLessons.Finished {
			if id == lessonID {
				return nil // 已存在，无需添加
			}
		}

		// 添加新的课时 ID
		studentLessons.Finished = append(studentLessons.Finished, lessonID)
		return tx.Save(&studentLessons).Error
	})
}

func (r *StudentLessonsRepo) SetLastOpened(ctx context.Context, studentID, lessonID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var studentLessons domain.StudentLessons
		err := tx.Set("gorm:query_option", "FOR UPDATE").First(&studentLessons, "student_id = ?", studentID).Error
		
		if err == gorm.ErrRecordNotFound {
			// 记录不存在，创建新记录
			studentLessons = domain.StudentLessons{
				StudentID:  studentID,
				LastOpened: lessonID,
			}
			return tx.Create(&studentLessons).Error
		} else if err != nil {
			return err
		}

		// 更新最后打开的课时
		studentLessons.LastOpened = lessonID
		return tx.Save(&studentLessons).Error
	})
}
