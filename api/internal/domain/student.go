package domain

import (
	"time"
)

type Student struct {
	ID               uint         `gorm:"primaryKey;autoIncrement" json:"id"`                        // 学生ID
	Name             string       `gorm:"size:255;not null" json:"name"`                             // 学生姓名
	Email            string       `gorm:"size:255;not null;uniqueIndex" json:"email"`                // 邮箱
	Password         string       `gorm:"size:255;not null" json:"password"`                         // 密码
	RegisteredAt     time.Time    `gorm:"not null" json:"registeredAt"`                              // 注册时间
	LastVisitAt      time.Time    `gorm:"not null" json:"lastVisitAt"`                               // 最后登录时间
	SchoolID         uint         `gorm:"not null;index" json:"schoolId"`                            // 所属学校ID
	AvailableModules []uint       `gorm:"serializer:json" json:"availableModules"`                   // 可用模块ID列表
	AvailableCourses []uint       `gorm:"serializer:json" json:"availableCourses"`                   // 可用课程ID列表
	AvailableOffers  []uint       `gorm:"serializer:json" json:"availableOffers"`                    // 可用优惠ID列表
	Verification     Verification `gorm:"embedded;embeddedPrefix:verification_" json:"verification"` // 邮箱验证信息
	Session          Session      `gorm:"embedded;embeddedPrefix:session_" json:"session"`           // 会话信息
	Blocked          bool         `gorm:"not null;default:false" json:"blocked"`                     // 是否被封禁
}

// IsModuleAvailable 检查学生是否可以访问指定模块
func (s *Student) IsModuleAvailable(m Module) bool {
	for _, id := range s.AvailableModules {
		if m.ID == id {
			return true
		}
	}
	return false
}

// IsModuleAvailableByID 检查学生是否可以访问指定模块（通过ID）
func (s *Student) IsModuleAvailableByID(moduleID uint) bool {
	for _, id := range s.AvailableModules {
		if id == moduleID {
			return true
		}
	}
	return false
}

// IsBlocked 检查学生是否被封禁
func (s *Student) IsBlocked() bool {
	return s.Blocked
}

// IsEmailVerified 检查学生邮箱是否已验证
func (s *Student) IsEmailVerified() bool {
	return s.Verification.Verified
}

// SetVerificationCode 设置邮箱验证码
func (s *Student) SetVerificationCode(code string) {
	s.Verification.Code = code
	s.Verification.Verified = false
}

// MarkAsVerified 标记邮箱为已验证
func (s *Student) MarkAsVerified() {
	s.Verification.Code = ""
	s.Verification.Verified = true
}

// SetPassword 设置密码哈希
func (s *Student) SetPassword(hash string) {
	s.Password = hash
}

// SetSession 设置会话信息
func (s *Student) SetSession(session Session) {
	s.Session = session
}

// GrantModuleAccess 授予模块访问权限
func (s *Student) GrantModuleAccess(moduleID uint) {
	if !s.IsModuleAvailableByID(moduleID) {
		s.AvailableModules = append(s.AvailableModules, moduleID)
	}
}

// RevokeModuleAccess 撤销模块访问权限
func (s *Student) RevokeModuleAccess(moduleID uint) {
	for i, id := range s.AvailableModules {
		if id == moduleID {
			s.AvailableModules = append(s.AvailableModules[:i], s.AvailableModules[i+1:]...)
			return
		}
	}
}

// GrantOfferAccess 授予优惠访问权限（包括相关模块）
func (s *Student) GrantOfferAccess(offerID uint, moduleIDs []uint) {
	// 添加优惠到可用列表
	found := false
	for _, id := range s.AvailableOffers {
		if id == offerID {
			found = true
			break
		}
	}
	if !found {
		s.AvailableOffers = append(s.AvailableOffers, offerID)
	}

	// 添加相关模块
	for _, moduleID := range moduleIDs {
		s.GrantModuleAccess(moduleID)
	}
}

// RevokeOfferAccess 撤销优惠访问权限（包括相关模块）
func (s *Student) RevokeOfferAccess(offerID uint, moduleIDs []uint) {
	// 移除优惠
	for i, id := range s.AvailableOffers {
		if id == offerID {
			s.AvailableOffers = append(s.AvailableOffers[:i], s.AvailableOffers[i+1:]...)
			break
		}
	}

	// 移除相关模块
	for _, moduleID := range moduleIDs {
		s.RevokeModuleAccess(moduleID)
	}
}

// NewStudent 创建新学生实例
func NewStudent(name, email, password string, schoolID uint) *Student {
	now := time.Now()
	return &Student{
		Name:         name,
		Email:        email,
		Password:     password,
		RegisteredAt: now,
		LastVisitAt:  now,
		SchoolID:     schoolID,
		Blocked:      false,
	}
}

type Verification struct {
	Code     string `gorm:"size:50" json:"code"`                    // 验证码
	Verified bool   `gorm:"not null;default:false" json:"verified"` // 是否已验证
}

type StudentLessons struct {
	StudentID  uint   `gorm:"primaryKey" json:"studentId"`     // 学生ID
	Finished   []uint `gorm:"serializer:json" json:"finished"` // 已完成课时ID列表
	LastOpened uint   `json:"lastOpened"`                      // 最后打开的课时ID
}

type StudentInfoShort struct {
	ID    uint   `json:"id"`    // 学生ID
	Name  string `json:"name"`  // 学生姓名
	Email string `json:"email"` // 学生邮箱
}

type UpdateStudentInput struct {
	Name      string `json:"name"`     // 姓名
	Email     string `json:"email"`    // 邮箱
	Verified  *bool  `json:"verified"` // 是否验证
	Blocked   *bool  `json:"blocked"`  // 是否封禁
	StudentID uint   `json:"-"`        // 学生ID（内部使用）
	SchoolID  uint   `json:"-"`        // 学校ID（内部使用）
}

type CreateStudentInput struct {
	Name     string `json:"name" binding:"required,min=2"`     // 姓名
	Email    string `json:"email" binding:"required,email"`    // 邮箱
	Password string `json:"password" binding:"required,min=6"` // 密码
	SchoolID uint   `json:"-"`                                 // 学校ID（内部使用）
}

// StudentSignUpInput 学生注册输入（Service 层使用）
type StudentSignUpInput struct {
	Name         string
	Email        string
	Password     string
	SchoolID     uint
	SchoolDomain string
	Verified     bool
}

// SchoolSignInInput 学校登录输入（Service 层使用）
type SchoolSignInInput struct {
	Email    string
	Password string
	SchoolID uint
}
