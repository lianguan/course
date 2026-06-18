package domain

import "errors"

var ErrFondyIsNotConnected = errors.New("fondy is not connected")

// School 学校实体
type School struct {
	ID           uint     // 学校ID
	Name         string   // 学校名称
	Subtitle     string   // 副标题
	Description  string   // 描述
	RegisteredAt int64    // 注册时间（Unix 时间戳）
	Admins       []Admin  // 管理员列表
	Courses      []Course // 课程列表
	Settings     Settings // 学校设置
}

// Settings 学校设置值对象
type Settings struct {
	Color               string      // 主题颜色
	Domains             []string    // 域名列表
	ContactInfo         ContactInfo // 联系信息
	Pages               Pages       // 页面内容
	ShowPaymentImages   bool        // 是否显示支付图片
	Logo                string      // Logo URL
	GoogleAnalyticsCode string      // Google Analytics 代码
	Fondy               Fondy       // Fondy 支付配置
	SendPulse           SendPulse   // SendPulse 邮件配置
	DisableRegistration bool        // 是否禁用注册
}

func (s Settings) GetDomain() string {
	if len(s.Domains) == 0 {
		return ""
	}
	return s.Domains[0]
}

// Fondy Fondy 支付配置值对象
type Fondy struct {
	MerchantID       string // 商户ID
	MerchantPassword string // 商户密码
	Connected        bool   // 是否已连接
}

// SendPulse SendPulse 邮件配置值对象
type SendPulse struct {
	ID        string // SendPulse ID
	Secret    string // SendPulse Secret
	ListID    string // 邮件列表ID
	Connected bool   // 是否已连接
}

// ContactInfo 联系信息值对象
type ContactInfo struct {
	BusinessName       string // 企业名称
	RegistrationNumber string // 注册号
	Address            string // 地址
	Email              string // 联系邮箱
	Phone              string // 联系电话
}

// Pages 页面内容值对象
type Pages struct {
	Confidential      string // 隐私政策
	ServiceAgreement  string // 服务协议
	NewsletterConsent string // 邮件订阅同意条款
}

// Admin 管理员实体
type Admin struct {
	ID       uint    // 管理员ID
	Name     string  // 管理员姓名
	Email    string  // 邮箱
	Password string  // 密码
	SchoolID uint    // 所属学校ID
	Session  Session // 会话信息
}

type UpdateSchoolSettingsInput struct {
	Name                *string
	Color               *string
	Domains             []string
	Email               *string
	ContactInfo         *UpdateSchoolSettingsContactInfo
	Pages               *UpdateSchoolSettingsPages
	ShowPaymentImages   *bool
	DisableRegistration *bool
	GoogleAnalyticsCode *string
	LogoURL             *string
}

type UpdateSchoolSettingsPages struct {
	Confidential      *string // 隐私政策
	ServiceAgreement  *string // 服务协议
	NewsletterConsent *string // 邮件订阅同意条款
}

type UpdateSchoolSettingsContactInfo struct {
	BusinessName       *string // 企业名称
	RegistrationNumber *string // 注册号
	Address            *string // 地址
	Email              *string // 联系邮箱
	Phone              *string // 联系电话
}

// ConnectFondyInput Fondy 支付连接输入（Service 层使用）
type ConnectFondyInput struct {
	SchoolID         uint
	MerchantID       string
	MerchantPassword string
}

// ConnectSendPulseInput SendPulse 邮件服务连接输入（Service 层使用）
type ConnectSendPulseInput struct {
	SchoolID uint
	ID       string
	Secret   string
	ListID   string
}
