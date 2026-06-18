package domain

import "time"

// Student 学生实体
type Student struct {
	ID               uint         // 学生ID
	Name             string       // 学生姓名
	Email            string       // 邮箱
	Password         string       // 密码
	RegisteredAt     int64        // 注册时间（Unix 时间戳）
	LastVisitAt      int64        // 最后登录时间（Unix 时间戳）
	SchoolID         uint         // 所属学校ID
	AvailableModules []uint       // 可用模块ID列表
	AvailableCourses []uint       // 可用课程ID列表
	AvailableOffers  []uint       // 可用优惠ID列表
	Verification     Verification // 邮箱验证信息
	Session          Session      // 会话信息
	Blocked          bool         // 是否被封禁
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
	now := time.Now().Unix()
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

// Verification 邮箱验证值对象
type Verification struct {
	Code     string // 验证码
	Verified bool   // 是否已验证
}

// StudentLessons 学生课时进度值对象
type StudentLessons struct {
	StudentID  uint   // 学生ID
	Finished   []uint // 已完成课时ID列表
	LastOpened uint   // 最后打开的课时ID
}

// StudentInfoShort 学生简要信息值对象
type StudentInfoShort struct {
	ID    uint   // 学生ID
	Name  string // 学生姓名
	Email string // 学生邮箱
}

// UpdateStudentInput 学生更新输入（Repository 层使用）
type UpdateStudentInput struct {
	Name      string // 姓名
	Email     string // 邮箱
	Verified  *bool  // 是否验证
	Blocked   *bool  // 是否封禁
	StudentID uint   // 学生ID（内部使用）
	SchoolID  uint   // 学校ID（内部使用）
}

// CreateStudentInput 学生创建输入（Service 层使用）
type CreateStudentInput struct {
	Name     string // 姓名
	Email    string // 邮箱
	Password string // 密码
	SchoolID uint   // 学校ID（内部使用）
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
