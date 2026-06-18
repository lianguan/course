package repository

import (
	"context"
	"errors"
	"time"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/repository/models"

	"gorm.io/gorm"
)

type StudentsRepo struct {
	db *gorm.DB
}

func NewStudentsRepo(db *gorm.DB) *StudentsRepo {
	return &StudentsRepo{db: db}
}

func (r *StudentsRepo) Create(ctx context.Context, student *domain.Student) error {
	var model models.StudentModel
	model.FromDomain(*student)
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *StudentsRepo) Update(ctx context.Context, inp domain.UpdateStudentInput) error {
	updates := map[string]interface{}{}

	if inp.Name != "" {
		updates["name"] = inp.Name
	}
	if inp.Email != "" {
		updates["email"] = inp.Email
	}
	if inp.Verified != nil {
		updates["verification_verified"] = *inp.Verified
	}
	if inp.Blocked != nil {
		updates["blocked"] = *inp.Blocked
	}

	return r.db.WithContext(ctx).
		Model(&models.StudentModel{}).
		Where("id = ? AND school_id = ?", inp.StudentID, inp.SchoolID).
		Updates(updates).Error
}

func (r *StudentsRepo) Delete(ctx context.Context, schoolID, studentID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", studentID, schoolID).
		Delete(&models.StudentModel{}).Error
}

func (r *StudentsRepo) GetByCredentials(ctx context.Context, schoolID uint, email, password string) (domain.Student, error) {
	var model models.StudentModel
	err := r.db.WithContext(ctx).
		Where("email = ? AND password = ? AND school_id = ? AND verification_verified = ?", email, password, schoolID, true).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Student{}, domain.ErrUserNotFound
		}
		return domain.Student{}, err
	}
	return model.ToDomain(), nil
}

func (r *StudentsRepo) GetByRefreshToken(ctx context.Context, schoolID uint, refreshToken string) (domain.Student, error) {
	var model models.StudentModel
	err := r.db.WithContext(ctx).
		Where("session_refresh_token = ? AND school_id = ? AND session_expires_at > ?", refreshToken, schoolID, time.Now()).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Student{}, domain.ErrUserNotFound
		}
		return domain.Student{}, err
	}
	return model.ToDomain(), nil
}

func (r *StudentsRepo) GetById(ctx context.Context, schoolID, id uint) (domain.Student, error) {
	var model models.StudentModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Student{}, domain.ErrUserNotFound
		}
		return domain.Student{}, err
	}
	return model.ToDomain(), nil
}

func (r *StudentsRepo) GetBySchool(ctx context.Context, schoolID uint, query domain.GetStudentsQuery) ([]domain.Student, int64, error) {
	var models []models.StudentModel
	var count int64

	db := r.db.WithContext(ctx).Model(&models.StudentModel{}).Where("school_id = ?", schoolID)

	if query.Search != "" {
		db = db.Where("name LIKE ? OR email LIKE ?", "%"+query.Search+"%", "%"+query.Search+"%")
	}
	if query.Verified != nil {
		db = db.Where("verification_verified = ?", *query.Verified)
	}
	if query.RegisterDateFrom != "" {
		db = db.Where("registered_at >= ?", query.RegisterDateFrom)
	}
	if query.RegisterDateTo != "" {
		db = db.Where("registered_at <= ?", query.RegisterDateTo)
	}
	if query.LastVisitDateFrom != "" {
		db = db.Where("last_visit_at >= ?", query.LastVisitDateFrom)
	}
	if query.LastVisitDateTo != "" {
		db = db.Where("last_visit_at <= ?", query.LastVisitDateTo)
	}

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if query.PaginationQuery.Limit > 0 {
		db = db.Limit(int(query.PaginationQuery.Limit))
	}
	if query.PaginationQuery.Skip > 0 {
		db = db.Offset(int(query.PaginationQuery.Skip))
	}

	err := db.Order("registered_at DESC").Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	students := make([]domain.Student, len(models))
	for i, m := range models {
		students[i] = m.ToDomain()
	}
	return students, count, nil
}

func (r *StudentsRepo) SetSession(ctx context.Context, studentID uint, session domain.Session) error {
	return r.db.WithContext(ctx).
		Model(&models.StudentModel{}).
		Where("id = ?", studentID).
		Updates(map[string]interface{}{
			"session_refresh_token": session.RefreshToken,
			"session_expires_at":    time.Unix(session.ExpiresAt, 0),
			"last_visit_at":         time.Now(),
		}).Error
}

func (r *StudentsRepo) GiveAccessToModule(ctx context.Context, studentID, moduleID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model models.StudentModel
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&model, studentID).Error; err != nil {
			return err
		}

		student := model.ToDomain()
		if student.IsModuleAvailableByID(moduleID) {
			return nil
		}

		student.GrantModuleAccess(moduleID)
		model.FromDomain(student)

		return tx.Model(&models.StudentModel{}).
			Where("id = ?", studentID).
			Update("available_modules", model.AvailableModules).Error
	})
}

func (r *StudentsRepo) AttachOffer(ctx context.Context, studentID, offerID uint, moduleIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model models.StudentModel
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&model, studentID).Error; err != nil {
			return err
		}

		student := model.ToDomain()
		student.GrantOfferAccess(offerID, moduleIDs)
		model.FromDomain(student)

		return tx.Model(&models.StudentModel{}).
			Where("id = ?", studentID).
			Updates(map[string]interface{}{
				"available_modules": model.AvailableModules,
				"available_offers":  model.AvailableOffers,
			}).Error
	})
}

func (r *StudentsRepo) DetachOffer(ctx context.Context, studentID, offerID uint, moduleIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model models.StudentModel
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&model, studentID).Error; err != nil {
			return err
		}

		student := model.ToDomain()
		student.RevokeOfferAccess(offerID, moduleIDs)
		model.FromDomain(student)

		return tx.Model(&models.StudentModel{}).
			Where("id = ?", studentID).
			Updates(map[string]interface{}{
				"available_modules": model.AvailableModules,
				"available_offers":  model.AvailableOffers,
			}).Error
	})
}

func (r *StudentsRepo) Verify(ctx context.Context, code string) (domain.Student, error) {
	var model models.StudentModel
	err := r.db.WithContext(ctx).
		Where("verification_code = ?", code).
		First(&model).Error
	if err != nil {
		return domain.Student{}, err
	}

	student := model.ToDomain()
	student.MarkAsVerified()
	model.FromDomain(student)

	err = r.db.WithContext(ctx).Save(&model).Error
	return model.ToDomain(), err
}
