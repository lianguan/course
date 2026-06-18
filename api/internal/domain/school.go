package domain

import (
	"errors"
	"time"
)

var ErrFondyIsNotConnected = errors.New("fondy is not connected")

type School struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"` // 学校ID
	Name         string    `gorm:"size:255;not null" json:"name"`      // 学校名称
	Subtitle     string    `gorm:"size:255" json:"subtitle"`           // 副标题
	Description  string    `gorm:"type:text" json:"description"`       // 描述
	RegisteredAt time.Time `gorm:"not null" json:"registeredAt"`       // 注册时间
	Admins       []Admin   `gorm:"-" json:"admins,omitempty"`          // 管理员列表
	Courses      []Course  `gorm:"-" json:"courses,omitempty"`         // 课程列表
	Settings     Settings  `gorm:"serializer:json" json:"settings"`    // 学校设置
}

type Settings struct {
	Color               string      `json:"color"`               // 主题颜色
	Domains             []string    `json:"domains"`             // 域名列表
	ContactInfo         ContactInfo `json:"contactInfo"`         // 联系信息
	Pages               Pages       `json:"pages"`               // 页面内容
	ShowPaymentImages   bool        `json:"showPaymentImages"`   // 是否显示支付图片
	Logo                string      `json:"logo"`                // Logo URL
	GoogleAnalyticsCode string      `json:"googleAnalyticsCode"` // Google Analytics 代码
	Fondy               Fondy       `json:"fondy"`               // Fondy 支付配置
	SendPulse           SendPulse   `json:"sendpulse"`           // SendPulse 邮件配置
	DisableRegistration bool        `json:"disableRegistration"` // 是否禁用注册
}

func (s Settings) GetDomain() string {
	if len(s.Domains) == 0 {
		return ""
	}
	return s.Domains[0]
}

type Fondy struct {
	MerchantID       string `json:"merchantId"`       // 商户ID
	MerchantPassword string `json:"merchantPassword"` // 商户密码
	Connected        bool   `json:"connected"`        // 是否已连接
}

type SendPulse struct {
	ID        string `json:"id"`        // SendPulse ID
	Secret    string `json:"secret"`    // SendPulse Secret
	ListID    string `json:"listId"`    // 邮件列表ID
	Connected bool   `json:"connected"` // 是否已连接
}

type ContactInfo struct {
	BusinessName       string `json:"businessName"`       // 企业名称
	RegistrationNumber string `json:"registrationNumber"` // 注册号
	Address            string `json:"address"`            // 地址
	Email              string `json:"email"`              // 联系邮箱
	Phone              string `json:"phone"`              // 联系电话
}

type Pages struct {
	Confidential      string `json:"confidential"`      // 隐私政策
	ServiceAgreement  string `json:"serviceAgreement"`  // 服务协议
	NewsletterConsent string `json:"newsletterConsent"` // 邮件订阅同意条款
}

type Admin struct {
	ID       uint    `gorm:"primaryKey;autoIncrement" json:"id"`              // 管理员ID
	Name     string  `gorm:"size:255;not null" json:"name"`                   // 管理员姓名
	Email    string  `gorm:"size:255;not null;uniqueIndex" json:"email"`      // 邮箱
	Password string  `gorm:"size:255;not null" json:"password"`               // 密码
	SchoolID uint    `gorm:"not null;index" json:"schoolId"`                  // 所属学校ID
	Session  Session `gorm:"embedded;embeddedPrefix:session_" json:"session"` // 会话信息
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
