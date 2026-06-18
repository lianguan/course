package models

import (
	"time"

	"ultrathreads/internal/domain"
)

// CourseModel 课程 GORM 模型
type CourseModel struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"size:255;not null"`
	Code        string    `gorm:"size:100;uniqueIndex"`
	Description string    `gorm:"type:text"`
	Color       string    `gorm:"size:50"`
	ImageURL    string    `gorm:"size:500"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	Published   bool      `gorm:"not null;default:false"`
}

// ToDomain 转换为 domain 实体
func (m *CourseModel) ToDomain() domain.Course {
	return domain.Course{
		ID:          m.ID,
		Name:        m.Name,
		Code:        m.Code,
		Description: m.Description,
		Color:       m.Color,
		ImageURL:    m.ImageURL,
		CreatedAt:   m.CreatedAt.Unix(),
		UpdatedAt:   m.UpdatedAt.Unix(),
		Published:   m.Published,
	}
}

// FromDomain 从 domain 实体转换
func (m *CourseModel) FromDomain(c domain.Course) {
	m.ID = c.ID
	m.Name = c.Name
	m.Code = c.Code
	m.Description = c.Description
	m.Color = c.Color
	m.ImageURL = c.ImageURL
	m.CreatedAt = time.Unix(c.CreatedAt, 0)
	m.UpdatedAt = time.Unix(c.UpdatedAt, 0)
	m.Published = c.Published
}

// ModuleModel 模块 GORM 模型
type ModuleModel struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:255;not null"`
	Position  uint   `gorm:"not null;default:0"`
	Published bool   `gorm:"not null;default:false"`
	CourseID  uint   `gorm:"not null;index"`
	PackageID uint   `gorm:"index"`
	SchoolID  uint   `gorm:"not null;index"`
}

// ToDomain 转换为 domain 实体
func (m *ModuleModel) ToDomain() domain.Module {
	return domain.Module{
		ID:        m.ID,
		Name:      m.Name,
		Position:  m.Position,
		Published: m.Published,
		CourseID:  m.CourseID,
		PackageID: m.PackageID,
		SchoolID:  m.SchoolID,
	}
}

// FromDomain 从 domain 实体转换
func (m *ModuleModel) FromDomain(mod domain.Module) {
	m.ID = mod.ID
	m.Name = mod.Name
	m.Position = mod.Position
	m.Published = mod.Published
	m.CourseID = mod.CourseID
	m.PackageID = mod.PackageID
	m.SchoolID = mod.SchoolID
}

// LessonModel 课时 GORM 模型
type LessonModel struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"size:255;not null"`
	Position  uint   `gorm:"not null;default:0"`
	Published bool   `gorm:"not null;default:false"`
	Content   string `gorm:"type:text"`
	SchoolID  uint   `gorm:"not null;index"`
}

// ToDomain 转换为 domain 实体
func (m *LessonModel) ToDomain() domain.Lesson {
	return domain.Lesson{
		ID:        m.ID,
		Name:      m.Name,
		Position:  m.Position,
		Published: m.Published,
		Content:   m.Content,
		SchoolID:  m.SchoolID,
	}
}

// FromDomain 从 domain 实体转换
func (m *LessonModel) FromDomain(l domain.Lesson) {
	m.ID = l.ID
	m.Name = l.Name
	m.Position = l.Position
	m.Published = l.Published
	m.Content = l.Content
	m.SchoolID = l.SchoolID
}

// LessonContentModel 课时内容 GORM 模型
type LessonContentModel struct {
	LessonID uint   `gorm:"primaryKey"`
	SchoolID uint   `gorm:"not null;index"`
	Content  string `gorm:"type:text"`
}

// ToDomain 转换为 domain 值对象
func (m *LessonContentModel) ToDomain() domain.LessonContent {
	return domain.LessonContent{
		LessonID: m.LessonID,
		SchoolID: m.SchoolID,
		Content:  m.Content,
	}
}

// FromDomain 从 domain 值对象转换
func (m *LessonContentModel) FromDomain(lc domain.LessonContent) {
	m.LessonID = lc.LessonID
	m.SchoolID = lc.SchoolID
	m.Content = lc.Content
}

// PackageModel 套餐 GORM 模型
type PackageModel struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `gorm:"size:255;not null"`
	CourseID uint   `gorm:"not null;index"`
	SchoolID uint   `gorm:"not null;index"`
}

// ToDomain 转换为 domain 实体
func (m *PackageModel) ToDomain() domain.Package {
	return domain.Package{
		ID:       m.ID,
		Name:     m.Name,
		CourseID: m.CourseID,
		SchoolID: m.SchoolID,
	}
}

// FromDomain 从 domain 实体转换
func (m *PackageModel) FromDomain(p domain.Package) {
	m.ID = p.ID
	m.Name = p.Name
	m.CourseID = p.CourseID
	m.SchoolID = p.SchoolID
}
