package repository

import (
	"context"
	"time"

	"ultrathreads/internal/domain"

	"gorm.io/gorm"
)

type SchoolsRepo struct {
	db *gorm.DB
}

func NewSchoolsRepo(db *gorm.DB) *SchoolsRepo {
	return &SchoolsRepo{db: db}
}

func (r *SchoolsRepo) Create(ctx context.Context, name string) (uint, error) {
	school := domain.School{
		Name:         name,
		RegisteredAt: time.Now(),
	}
	err := r.db.WithContext(ctx).Create(&school).Error
	return school.ID, err
}

func (r *SchoolsRepo) GetByDomain(ctx context.Context, domainName string) (domain.School, error) {
	var school domain.School
	// 使用 MySQL JSON 函数直接在数据库层搜索域名
	err := r.db.WithContext(ctx).
		Where("JSON_CONTAINS(settings->'$.domains', ?)", `"`+domainName+`"`).
		First(&school).Error
	if err != nil {
		return domain.School{}, err
	}

	// Load courses for this school
	if err := r.loadCourses(ctx, &school); err != nil {
		return domain.School{}, err
	}

	return school, nil
}

func (r *SchoolsRepo) GetById(ctx context.Context, id uint) (domain.School, error) {
	var school domain.School
	err := r.db.WithContext(ctx).First(&school, id).Error
	if err != nil {
		return domain.School{}, err
	}

	// Load courses for this school
	if err := r.loadCourses(ctx, &school); err != nil {
		return domain.School{}, err
	}

	return school, nil
}

func (r *SchoolsRepo) UpdateSettings(ctx context.Context, id uint, inp domain.UpdateSchoolSettingsInput) error {
	var school domain.School
	if err := r.db.WithContext(ctx).First(&school, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{}

	if inp.Name != nil {
		updates["name"] = *inp.Name
	}

	settings := school.Settings

	if inp.Color != nil {
		settings.Color = *inp.Color
	}
	if inp.Domains != nil {
		settings.Domains = inp.Domains
	}
	if inp.ShowPaymentImages != nil {
		settings.ShowPaymentImages = *inp.ShowPaymentImages
	}
	if inp.DisableRegistration != nil {
		settings.DisableRegistration = *inp.DisableRegistration
	}
	if inp.GoogleAnalyticsCode != nil {
		settings.GoogleAnalyticsCode = *inp.GoogleAnalyticsCode
	}
	if inp.LogoURL != nil {
		settings.Logo = *inp.LogoURL
	}
	if inp.ContactInfo != nil {
		if inp.ContactInfo.Address != nil {
			settings.ContactInfo.Address = *inp.ContactInfo.Address
		}
		if inp.ContactInfo.BusinessName != nil {
			settings.ContactInfo.BusinessName = *inp.ContactInfo.BusinessName
		}
		if inp.ContactInfo.Email != nil {
			settings.ContactInfo.Email = *inp.ContactInfo.Email
		}
		if inp.ContactInfo.Phone != nil {
			settings.ContactInfo.Phone = *inp.ContactInfo.Phone
		}
		if inp.ContactInfo.RegistrationNumber != nil {
			settings.ContactInfo.RegistrationNumber = *inp.ContactInfo.RegistrationNumber
		}
	}
	if inp.Pages != nil {
		if inp.Pages.Confidential != nil {
			settings.Pages.Confidential = *inp.Pages.Confidential
		}
		if inp.Pages.NewsletterConsent != nil {
			settings.Pages.NewsletterConsent = *inp.Pages.NewsletterConsent
		}
		if inp.Pages.ServiceAgreement != nil {
			settings.Pages.ServiceAgreement = *inp.Pages.ServiceAgreement
		}
	}

	updates["settings"] = settings

	return r.db.WithContext(ctx).Model(&domain.School{}).Where("id = ?", id).Updates(updates).Error
}

func (r *SchoolsRepo) SetFondyCredentials(ctx context.Context, id uint, fondy domain.Fondy) error {
	var school domain.School
	if err := r.db.WithContext(ctx).First(&school, id).Error; err != nil {
		return err
	}
	school.Settings.Fondy = fondy
	return r.db.WithContext(ctx).Model(&domain.School{}).Where("id = ?", id).Update("settings", school.Settings).Error
}

func (r *SchoolsRepo) loadCourses(ctx context.Context, school *domain.School) error {
	var courses []domain.Course
	err := r.db.WithContext(ctx).
		Joins("INNER JOIN modules ON modules.course_id = courses.id").
		Where("modules.school_id = ?", school.ID).
		Distinct().
		Find(&courses).Error
	if err != nil {
		return err
	}
	school.Courses = courses
	return nil
}
