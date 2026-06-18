package domain

import "time"

type Course struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`      // 课程ID
	Name        string    `gorm:"size:255;not null" json:"name"`           // 课程名称
	Code        string    `gorm:"size:100;uniqueIndex" json:"code"`        // 课程编码
	Description string    `gorm:"type:text" json:"description"`            // 课程描述
	Color       string    `gorm:"size:50" json:"color"`                    // 主题颜色
	ImageURL    string    `gorm:"size:500" json:"imageUrl"`                // 封面图片URL
	CreatedAt   time.Time `gorm:"not null" json:"createdAt"`               // 创建时间
	UpdatedAt   time.Time `gorm:"not null" json:"updatedAt"`               // 更新时间
	Published   bool      `gorm:"not null;default:false" json:"published"` // 是否已发布
}

type Module struct {
	ID        uint     `gorm:"primaryKey;autoIncrement" json:"id"`       // 模块ID
	Name      string   `gorm:"size:255;not null" json:"name"`            // 模块名称
	Position  uint     `gorm:"not null;default:0" json:"position"`       // 排序位置
	Published bool     `gorm:"not null;default:false" json:"published"`  // 是否已发布
	CourseID  uint     `gorm:"not null;index" json:"courseId"`           // 所属课程ID
	PackageID uint     `gorm:"index" json:"packageId,omitempty"`         // 所属套餐ID
	SchoolID  uint     `gorm:"not null;index" json:"schoolId"`           // 所属学校ID
	Lessons   []Lesson `gorm:"serializer:json" json:"lessons,omitempty"` // 课时列表
	Survey    Survey   `gorm:"serializer:json" json:"survey,omitempty"`  // 调查问卷
}

type Lesson struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`      // 课时ID
	Name      string `gorm:"size:255;not null" json:"name"`           // 课时名称
	Position  uint   `gorm:"not null;default:0" json:"position"`      // 排序位置
	Published bool   `gorm:"not null;default:false" json:"published"` // 是否已发布
	Content   string `gorm:"type:text" json:"content,omitempty"`      // 课时内容
	SchoolID  uint   `gorm:"not null;index" json:"schoolId"`          // 所属学校ID
}

type LessonContent struct {
	LessonID uint   `gorm:"primaryKey" json:"lessonId"`     // 课时ID
	SchoolID uint   `gorm:"not null;index" json:"schoolId"` // 所属学校ID
	Content  string `gorm:"type:text" json:"content"`       // 课时内容
}

type Package struct {
	ID       uint     `gorm:"primaryKey;autoIncrement" json:"id"` // 套餐ID
	Name     string   `gorm:"size:255;not null" json:"name"`      // 套餐名称
	CourseID uint     `gorm:"not null;index" json:"courseId"`     // 所属课程ID
	SchoolID uint     `gorm:"not null;index" json:"schoolId"`     // 所属学校ID
	Modules  []Module `gorm:"-" json:"modules"`                   // 模块列表
}

type ModuleContent struct {
	Lessons []Lesson `gorm:"serializer:json" json:"lessons"` // 课时列表
	Survey  Survey   `gorm:"serializer:json" json:"survey"`  // 调查问卷
}

// IsPublished 检查模块是否已发布
func (m *Module) IsPublished() bool {
	return m.Published
}

// GetPublishedLessons 获取已发布的课时
func (m *Module) GetPublishedLessons() []Lesson {
	publishedLessons := make([]Lesson, 0)
	for _, lesson := range m.Lessons {
		if lesson.Published {
			publishedLessons = append(publishedLessons, lesson)
		}
	}
	return publishedLessons
}

// GetLessonIDs 获取所有课时ID
func (m *Module) GetLessonIDs() []uint {
	lessonIDs := make([]uint, 0, len(m.Lessons))
	for _, lesson := range m.Lessons {
		lessonIDs = append(lessonIDs, lesson.ID)
	}
	return lessonIDs
}

// GetPublishedLessonIDs 获取已发布课时的ID
func (m *Module) GetPublishedLessonIDs() []uint {
	lessonIDs := make([]uint, 0)
	for _, lesson := range m.Lessons {
		if lesson.Published {
			lessonIDs = append(lessonIDs, lesson.ID)
		}
	}
	return lessonIDs
}

// SortLessons 按位置排序课时
func (m *Module) SortLessons() {
	for i := 0; i < len(m.Lessons); i++ {
		for j := i + 1; j < len(m.Lessons); j++ {
			if m.Lessons[i].Position > m.Lessons[j].Position {
				m.Lessons[i], m.Lessons[j] = m.Lessons[j], m.Lessons[i]
			}
		}
	}
}

// IsPublished 检查课时是否已发布
func (l *Lesson) IsPublished() bool {
	return l.Published
}

// NewModule 创建新模块
func NewModule(name string, position uint, courseID, schoolID uint) *Module {
	return &Module{
		Name:      name,
		Position:  position,
		CourseID:  courseID,
		SchoolID:  schoolID,
		Published: false,
	}
}

// NewLesson 创建新课时
func NewLesson(name string, position uint, schoolID uint) *Lesson {
	return &Lesson{
		Name:      name,
		Position:  position,
		SchoolID:  schoolID,
		Published: false,
	}
}

// UpdateCourseInput 课程更新输入（Repository 层使用）
type UpdateCourseInput struct {
	ID          uint
	SchoolID    uint
	Name        *string
	ImageURL    *string
	Description *string
	Color       *string
	Published   *bool
}

// UpdateModuleInput 模块更新输入（Repository 层使用）
type UpdateModuleInput struct {
	ID        uint
	SchoolID  uint
	Name      string
	Position  *uint
	Published *bool
}

// UpdateLessonInput 课时更新输入（Repository 层使用）
type UpdateLessonInput struct {
	ID        uint
	SchoolID  uint
	Name      string
	Content   string
	Position  *uint
	Published *bool
}

// UpdatePackageInput 套餐更新输入（Repository 层使用）
type UpdatePackageInput struct {
	ID       uint
	SchoolID uint
	Name     string
	Modules  []uint
}

// UpdateOfferInput 优惠更新输入（Repository 层使用）
type UpdateOfferInput struct {
	ID            uint
	SchoolID      uint
	Name          string
	Description   string
	Benefits      []string
	Price         *Price
	Packages      []uint
	PaymentMethod *PaymentMethod
}

// CreateModuleInput 模块创建输入（Service 层使用）
type CreateModuleInput struct {
	SchoolID uint
	CourseID uint
	Name     string
	Position uint
}

// AddLessonInput 课时创建输入（Service 层使用）
type AddLessonInput struct {
	ModuleID uint
	SchoolID uint
	Name     string
	Position uint
}

// CreatePackageInput 套餐创建输入（Service 层使用）
type CreatePackageInput struct {
	CourseID uint
	SchoolID uint
	Name     string
	Modules  []uint
}
