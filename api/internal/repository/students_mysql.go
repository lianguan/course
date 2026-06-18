package repository

import (
	"context"
	"errors"
	"time"

	"ultrathreads/internal/domain"

	"gorm.io/gorm"
)

type StudentsRepo struct {
	db *gorm.DB
}

func NewStudentsRepo(db *gorm.DB) *StudentsRepo {
	return &StudentsRepo{db: db}
}

func (r *StudentsRepo) Create(ctx context.Context, student *domain.Student) error {
	return r.db.WithContext(ctx).Create(student).Error
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
		Model(&domain.Student{}).
		Where("id = ? AND school_id = ?", inp.StudentID, inp.SchoolID).
		Updates(updates).Error
}

func (r *StudentsRepo) Delete(ctx context.Context, schoolID, studentID uint) error {
	return r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", studentID, schoolID).
		Delete(&domain.Student{}).Error
}

func (r *StudentsRepo) GetByCredentials(ctx context.Context, schoolID uint, email, password string) (domain.Student, error) {
	var student domain.Student
	err := r.db.WithContext(ctx).
		Where("email = ? AND password = ? AND school_id = ? AND verification_verified = ?", email, password, schoolID, true).
		First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Student{}, domain.ErrUserNotFound
		}
		return domain.Student{}, err
	}
	return student, nil
}

func (r *StudentsRepo) GetByRefreshToken(ctx context.Context, schoolID uint, refreshToken string) (domain.Student, error) {
	var student domain.Student
	err := r.db.WithContext(ctx).
		Where("session_refresh_token = ? AND school_id = ? AND session_expires_at > ?", refreshToken, schoolID, time.Now()).
		First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Student{}, domain.ErrUserNotFound
		}
		return domain.Student{}, err
	}
	return student, nil
}

func (r *StudentsRepo) GetById(ctx context.Context, schoolID, id uint) (domain.Student, error) {
	var student domain.Student
	err := r.db.WithContext(ctx).
		Where("id = ? AND school_id = ?", id, schoolID).
		First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Student{}, domain.ErrUserNotFound
		}
		return domain.Student{}, err
	}
	return student, nil
}

func (r *StudentsRepo) GetBySchool(ctx context.Context, schoolID uint, query domain.GetStudentsQuery) ([]domain.Student, int64, error) {
	var students []domain.Student
	var count int64

	db := r.db.WithContext(ctx).Model(&domain.Student{}).Where("school_id = ?", schoolID)

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

	err := db.Order("registered_at DESC").Find(&students).Error
	return students, count, err
}

func (r *StudentsRepo) SetSession(ctx context.Context, studentID uint, session domain.Session) error {
	return r.db.WithContext(ctx).
		Model(&domain.Student{}).
		Where("id = ?", studentID).
		Updates(map[string]interface{}{
			"session_refresh_token": session.RefreshToken,
			"session_expires_at":    session.ExpiresAt,
			"last_visit_at":         time.Now(),
		}).Error
}

func (r *StudentsRepo) GiveAccessToModule(ctx context.Context, studentID, moduleID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var student domain.Student
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&student, studentID).Error; err != nil {
			return err
		}

		// 检查是否已存在
		for _, m := range student.AvailableModules {
			if m == moduleID {
				return nil
			}
		}

		modules := append(student.AvailableModules, moduleID)
		return tx.Model(&domain.Student{}).
			Where("id = ?", studentID).
			Update("available_modules", modules).Error
	})
}

func (r *StudentsRepo) AttachOffer(ctx context.Context, studentID, offerID uint, moduleIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var student domain.Student
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&student, studentID).Error; err != nil {
			return err
		}

		// 添加模块（去重）
		modules := student.AvailableModules
		for _, mid := range moduleIDs {
			found := false
			for _, m := range modules {
				if m == mid {
					found = true
					break
				}
			}
			if !found {
				modules = append(modules, mid)
			}
		}

		// 添加优惠（去重）
		offers := student.AvailableOffers
		found := false
		for _, o := range offers {
			if o == offerID {
				found = true
				break
			}
		}
		if !found {
			offers = append(offers, offerID)
		}

		return tx.Model(&domain.Student{}).
			Where("id = ?", studentID).
			Updates(map[string]interface{}{
				"available_modules": modules,
				"available_offers":  offers,
			}).Error
	})
}

func (r *StudentsRepo) DetachOffer(ctx context.Context, studentID, offerID uint, moduleIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var student domain.Student
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&student, studentID).Error; err != nil {
			return err
		}

		// 移除模块
		removeSet := make(map[uint]bool)
		for _, mid := range moduleIDs {
			removeSet[mid] = true
		}
		var filtered []uint
		for _, m := range student.AvailableModules {
			if !removeSet[m] {
				filtered = append(filtered, m)
			}
		}

		// 移除优惠
		var filteredOffers []uint
		for _, o := range student.AvailableOffers {
			if o != offerID {
				filteredOffers = append(filteredOffers, o)
			}
		}

		return tx.Model(&domain.Student{}).
			Where("id = ?", studentID).
			Updates(map[string]interface{}{
				"available_modules": filtered,
				"available_offers":  filteredOffers,
			}).Error
	})
}

func (r *StudentsRepo) Verify(ctx context.Context, code string) (domain.Student, error) {
	var student domain.Student
	err := r.db.WithContext(ctx).
		Where("verification_code = ?", code).
		First(&student).Error
	if err != nil {
		return domain.Student{}, err
	}

	student.Verification.Verified = true
	student.Verification.Code = ""
	err = r.db.WithContext(ctx).Save(&student).Error

	return student, err
}
