package models

import (
	"time"

	"ultrathreads/internal/domain"
)

// UserModel 用户 GORM 模型
type UserModel struct {
	ID               uint              `gorm:"primaryKey;autoIncrement"`
	Name             string            `gorm:"size:255;not null"`
	Email            string            `gorm:"size:255;not null;uniqueIndex"`
	Phone            string            `gorm:"size:50"`
	Password         string            `gorm:"size:255;not null"`
	RegisteredAt     time.Time         `gorm:"not null"`
	LastVisitAt      time.Time         `gorm:"not null"`
	VerificationCode string            `gorm:"size:50"`
	VerificationDone bool              `gorm:"not null;default:false"`
	Schools          JSONSlice[uint]   `gorm:"type:json"`
}

// ToDomain 转换为 domain 实体
func (m *UserModel) ToDomain() domain.User {
	return domain.User{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		Phone:        m.Phone,
		Password:     m.Password,
		RegisteredAt: m.RegisteredAt.Unix(),
		LastVisitAt:  m.LastVisitAt.Unix(),
		Verification: domain.Verification{
			Code:     m.VerificationCode,
			Verified: m.VerificationDone,
		},
		Schools: m.Schools,
	}
}

// FromDomain 从 domain 实体转换
func (m *UserModel) FromDomain(u domain.User) {
	m.ID = u.ID
	m.Name = u.Name
	m.Email = u.Email
	m.Phone = u.Phone
	m.Password = u.Password
	m.RegisteredAt = time.Unix(u.RegisteredAt, 0)
	m.LastVisitAt = time.Unix(u.LastVisitAt, 0)
	m.VerificationCode = u.Verification.Code
	m.VerificationDone = u.Verification.Verified
	m.Schools = u.Schools
}

// StudentModel 学生 GORM 模型
type StudentModel struct {
	ID                   uint              `gorm:"primaryKey;autoIncrement"`
	Name                 string            `gorm:"size:255;not null"`
	Email                string            `gorm:"size:255;not null;uniqueIndex:idx_school_email"`
	Password             string            `gorm:"size:255;not null"`
	RegisteredAt         time.Time         `gorm:"not null"`
	LastVisitAt          time.Time         `gorm:"not null"`
	SchoolID             uint              `gorm:"not null;uniqueIndex:idx_school_email"`
	AvailableModules     JSONSlice[uint]   `gorm:"type:json"`
	AvailableCourses     JSONSlice[uint]   `gorm:"type:json"`
	AvailableOffers      JSONSlice[uint]   `gorm:"type:json"`
	VerificationCode     string            `gorm:"size:50"`
	VerificationDone     bool              `gorm:"not null;default:false"`
	SessionRefreshToken  string            `gorm:"size:500"`
	SessionExpiresAt     time.Time
	Blocked              bool              `gorm:"not null;default:false"`
}

// ToDomain 转换为 domain 实体
func (m *StudentModel) ToDomain() domain.Student {
	return domain.Student{
		ID:               m.ID,
		Name:             m.Name,
		Email:            m.Email,
		Password:         m.Password,
		RegisteredAt:     m.RegisteredAt.Unix(),
		LastVisitAt:      m.LastVisitAt.Unix(),
		SchoolID:         m.SchoolID,
		AvailableModules: m.AvailableModules,
		AvailableCourses: m.AvailableCourses,
		AvailableOffers:  m.AvailableOffers,
		Verification: domain.Verification{
			Code:     m.VerificationCode,
			Verified: m.VerificationDone,
		},
		Session: domain.Session{
			RefreshToken: m.SessionRefreshToken,
			ExpiresAt:    m.SessionExpiresAt.Unix(),
		},
		Blocked: m.Blocked,
	}
}

// FromDomain 从 domain 实体转换
func (m *StudentModel) FromDomain(s domain.Student) {
	m.ID = s.ID
	m.Name = s.Name
	m.Email = s.Email
	m.Password = s.Password
	m.RegisteredAt = time.Unix(s.RegisteredAt, 0)
	m.LastVisitAt = time.Unix(s.LastVisitAt, 0)
	m.SchoolID = s.SchoolID
	m.AvailableModules = s.AvailableModules
	m.AvailableCourses = s.AvailableCourses
	m.AvailableOffers = s.AvailableOffers
	m.VerificationCode = s.Verification.Code
	m.VerificationDone = s.Verification.Verified
	m.SessionRefreshToken = s.Session.RefreshToken
	m.SessionExpiresAt = time.Unix(s.Session.ExpiresAt, 0)
	m.Blocked = s.Blocked
}

// AdminModel 管理员 GORM 模型
type AdminModel struct {
	ID                  uint      `gorm:"primaryKey;autoIncrement"`
	Name                string    `gorm:"size:255;not null"`
	Email               string    `gorm:"size:255;not null;uniqueIndex"`
	Password            string    `gorm:"size:255;not null"`
	SchoolID            uint      `gorm:"not null;index"`
	SessionRefreshToken string    `gorm:"size:500"`
	SessionExpiresAt    time.Time
}

// ToDomain 转换为 domain 实体
func (m *AdminModel) ToDomain() domain.Admin {
	return domain.Admin{
		ID:       m.ID,
		Name:     m.Name,
		Email:    m.Email,
		Password: m.Password,
		SchoolID: m.SchoolID,
		Session: domain.Session{
			RefreshToken: m.SessionRefreshToken,
			ExpiresAt:    m.SessionExpiresAt.Unix(),
		},
	}
}

// FromDomain 从 domain 实体转换
func (m *AdminModel) FromDomain(a domain.Admin) {
	m.ID = a.ID
	m.Name = a.Name
	m.Email = a.Email
	m.Password = a.Password
	m.SchoolID = a.SchoolID
	m.SessionRefreshToken = a.Session.RefreshToken
	m.SessionExpiresAt = time.Unix(a.Session.ExpiresAt, 0)
}

// SchoolModel 学校 GORM 模型
type SchoolModel struct {
	ID           uint           `gorm:"primaryKey;autoIncrement"`
	Name         string         `gorm:"size:255;not null"`
	Subtitle     string         `gorm:"size:255"`
	Description  string         `gorm:"type:text"`
	RegisteredAt time.Time      `gorm:"not null"`
	Settings     JSONSettings   `gorm:"type:json"`
}

// ToDomain 转换为 domain 实体
func (m *SchoolModel) ToDomain() domain.School {
	return domain.School{
		ID:           m.ID,
		Name:         m.Name,
		Subtitle:     m.Subtitle,
		Description:  m.Description,
		RegisteredAt: m.RegisteredAt.Unix(),
		Settings:     m.Settings.ToDomain(),
	}
}

// FromDomain 从 domain 实体转换
func (m *SchoolModel) FromDomain(s domain.School) {
	m.ID = s.ID
	m.Name = s.Name
	m.Subtitle = s.Subtitle
	m.Description = s.Description
	m.RegisteredAt = time.Unix(s.RegisteredAt, 0)
	m.Settings.FromDomain(s.Settings)
}

// StudentLessonsModel 学生课时进度 GORM 模型
type StudentLessonsModel struct {
	StudentID  uint          `gorm:"primaryKey"`
	Finished   JSONSlice[uint] `gorm:"type:json"`
	LastOpened uint
}

// ToDomain 转换为 domain 值对象
func (m *StudentLessonsModel) ToDomain() domain.StudentLessons {
	return domain.StudentLessons{
		StudentID:  m.StudentID,
		Finished:   m.Finished,
		LastOpened: m.LastOpened,
	}
}

// FromDomain 从 domain 值对象转换
func (m *StudentLessonsModel) FromDomain(sl domain.StudentLessons) {
	m.StudentID = sl.StudentID
	m.Finished = sl.Finished
	m.LastOpened = sl.LastOpened
}
