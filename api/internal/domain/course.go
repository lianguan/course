package domain

// Course 课程实体
type Course struct {
	ID          uint   // 课程ID
	Name        string // 课程名称
	Code        string // 课程编码
	Description string // 课程描述
	Color       string // 主题颜色
	ImageURL    string // 封面图片URL
	CreatedAt   int64  // 创建时间（Unix 时间戳）
	UpdatedAt   int64  // 更新时间（Unix 时间戳）
	Published   bool   // 是否已发布
}

// Module 模块实体
type Module struct {
	ID        uint     // 模块ID
	Name      string   // 模块名称
	Position  uint     // 排序位置
	Published bool     // 是否已发布
	CourseID  uint     // 所属课程ID
	PackageID uint     // 所属套餐ID
	SchoolID  uint     // 所属学校ID
	Lessons   []Lesson // 课时列表
	Survey    Survey   // 调查问卷
}

// Lesson 课时实体
type Lesson struct {
	ID        uint   // 课时ID
	Name      string // 课时名称
	Position  uint   // 排序位置
	Published bool   // 是否已发布
	Content   string // 课时内容
	SchoolID  uint   // 所属学校ID
}

// LessonContent 课时内容值对象
type LessonContent struct {
	LessonID uint   // 课时ID
	SchoolID uint   // 所属学校ID
	Content  string // 课时内容
}

// Package 套餐实体
type Package struct {
	ID       uint     // 套餐ID
	Name     string   // 套餐名称
	CourseID uint     // 所属课程ID
	SchoolID uint     // 所属学校ID
	Modules  []Module // 模块列表
}

// ModuleContent 模块内容值对象
type ModuleContent struct {
	Lessons []Lesson // 课时列表
	Survey  Survey   // 调查问卷
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
