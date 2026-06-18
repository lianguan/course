package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"ultrathreads/internal/domain"
)

// JSONSlice 通用 JSON 数组类型
type JSONSlice[T any] []T

func (j JSONSlice[T]) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONSlice[T]) Scan(value interface{}) error {
	if value == nil {
		*j = []T{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return errors.New("unsupported type for JSONSlice")
	}
	return json.Unmarshal(bytes, j)
}

// JSONSettings 学校设置 JSON 类型
type JSONSettings struct {
	Color               string          `json:"color"`
	Domains             []string        `json:"domains"`
	ContactInfo         ContactInfoJSON `json:"contactInfo"`
	Pages               PagesJSON       `json:"pages"`
	ShowPaymentImages   bool            `json:"showPaymentImages"`
	Logo                string          `json:"logo"`
	GoogleAnalyticsCode string          `json:"googleAnalyticsCode"`
	Fondy               FondyJSON       `json:"fondy"`
	SendPulse           SendPulseJSON   `json:"sendpulse"`
	DisableRegistration bool            `json:"disableRegistration"`
}

func (j JSONSettings) Value() (driver.Value, error) {
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONSettings) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return errors.New("unsupported type for JSONSettings")
	}
	return json.Unmarshal(bytes, j)
}

// ToDomain 转换为 domain 值对象
func (j *JSONSettings) ToDomain() domain.Settings {
	return domain.Settings{
		Color:               j.Color,
		Domains:             j.Domains,
		ContactInfo:         j.ContactInfo.ToDomain(),
		Pages:               j.Pages.ToDomain(),
		ShowPaymentImages:   j.ShowPaymentImages,
		Logo:                j.Logo,
		GoogleAnalyticsCode: j.GoogleAnalyticsCode,
		Fondy:               j.Fondy.ToDomain(),
		SendPulse:           j.SendPulse.ToDomain(),
		DisableRegistration: j.DisableRegistration,
	}
}

// FromDomain 从 domain 值对象转换
func (j *JSONSettings) FromDomain(s domain.Settings) {
	j.Color = s.Color
	j.Domains = s.Domains
	j.ContactInfo.FromDomain(s.ContactInfo)
	j.Pages.FromDomain(s.Pages)
	j.ShowPaymentImages = s.ShowPaymentImages
	j.Logo = s.Logo
	j.GoogleAnalyticsCode = s.GoogleAnalyticsCode
	j.Fondy.FromDomain(s.Fondy)
	j.SendPulse.FromDomain(s.SendPulse)
	j.DisableRegistration = s.DisableRegistration
}

// ContactInfoJSON 联系信息 JSON 类型
type ContactInfoJSON struct {
	BusinessName       string `json:"businessName"`
	RegistrationNumber string `json:"registrationNumber"`
	Address            string `json:"address"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
}

func (c *ContactInfoJSON) ToDomain() domain.ContactInfo {
	return domain.ContactInfo{
		BusinessName:       c.BusinessName,
		RegistrationNumber: c.RegistrationNumber,
		Address:            c.Address,
		Email:              c.Email,
		Phone:              c.Phone,
	}
}

func (c *ContactInfoJSON) FromDomain(ci domain.ContactInfo) {
	c.BusinessName = ci.BusinessName
	c.RegistrationNumber = ci.RegistrationNumber
	c.Address = ci.Address
	c.Email = ci.Email
	c.Phone = ci.Phone
}

// PagesJSON 页面内容 JSON 类型
type PagesJSON struct {
	Confidential      string `json:"confidential"`
	ServiceAgreement  string `json:"serviceAgreement"`
	NewsletterConsent string `json:"newsletterConsent"`
}

func (p *PagesJSON) ToDomain() domain.Pages {
	return domain.Pages{
		Confidential:      p.Confidential,
		ServiceAgreement:  p.ServiceAgreement,
		NewsletterConsent: p.NewsletterConsent,
	}
}

func (p *PagesJSON) FromDomain(pg domain.Pages) {
	p.Confidential = pg.Confidential
	p.ServiceAgreement = pg.ServiceAgreement
	p.NewsletterConsent = pg.NewsletterConsent
}

// FondyJSON Fondy 支付配置 JSON 类型
type FondyJSON struct {
	MerchantID       string `json:"merchantId"`
	MerchantPassword string `json:"merchantPassword"`
	Connected        bool   `json:"connected"`
}

func (f *FondyJSON) ToDomain() domain.Fondy {
	return domain.Fondy{
		MerchantID:       f.MerchantID,
		MerchantPassword: f.MerchantPassword,
		Connected:        f.Connected,
	}
}

func (f *FondyJSON) FromDomain(fd domain.Fondy) {
	f.MerchantID = fd.MerchantID
	f.MerchantPassword = fd.MerchantPassword
	f.Connected = fd.Connected
}

// SendPulseJSON SendPulse 邮件配置 JSON 类型
type SendPulseJSON struct {
	ID        string `json:"id"`
	Secret    string `json:"secret"`
	ListID    string `json:"listId"`
	Connected bool   `json:"connected"`
}

func (s *SendPulseJSON) ToDomain() domain.SendPulse {
	return domain.SendPulse{
		ID:        s.ID,
		Secret:    s.Secret,
		ListID:    s.ListID,
		Connected: s.Connected,
	}
}

func (s *SendPulseJSON) FromDomain(sp domain.SendPulse) {
	s.ID = sp.ID
	s.Secret = sp.Secret
	s.ListID = sp.ListID
	s.Connected = sp.Connected
}

// OfferModel 优惠 GORM 模型
type OfferModel struct {
	ID                     uint            `gorm:"primaryKey;autoIncrement"`
	Name                   string          `gorm:"size:255;not null"`
	Description            string          `gorm:"type:text"`
	Benefits               JSONSlice[string] `gorm:"type:json"`
	SchoolID               uint            `gorm:"not null;index"`
	PackageIDs             JSONSlice[uint] `gorm:"type:json;column:packages"`
	PriceValue             uint            `gorm:"not null;default:0"`
	PriceCurrency          string          `gorm:"size:10;not null;default:'USD'"`
	PaymentMethodUsesProvider bool         `gorm:"not null;default:false"`
	PaymentMethodProvider  string          `gorm:"size:50"`
}

func (m *OfferModel) ToDomain() domain.Offer {
	return domain.Offer{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Benefits:    m.Benefits,
		SchoolID:    m.SchoolID,
		PackageIDs:  m.PackageIDs,
		Price: domain.Price{
			Value:    m.PriceValue,
			Currency: m.PriceCurrency,
		},
		PaymentMethod: domain.PaymentMethod{
			UsesProvider: m.PaymentMethodUsesProvider,
			Provider:     m.PaymentMethodProvider,
		},
	}
}

func (m *OfferModel) FromDomain(o domain.Offer) {
	m.ID = o.ID
	m.Name = o.Name
	m.Description = o.Description
	m.Benefits = o.Benefits
	m.SchoolID = o.SchoolID
	m.PackageIDs = o.PackageIDs
	m.PriceValue = o.Price.Value
	m.PriceCurrency = o.Price.Currency
	m.PaymentMethodUsesProvider = o.PaymentMethod.UsesProvider
	m.PaymentMethodProvider = o.PaymentMethod.Provider
}

// PromoCodeModel 优惠码 GORM 模型
type PromoCodeModel struct {
	ID                 uint            `gorm:"primaryKey;autoIncrement"`
	SchoolID           uint            `gorm:"not null;index"`
	Code               string          `gorm:"size:100;not null;uniqueIndex"`
	DiscountPercentage int             `gorm:"not null"`
	ExpiresAt          time.Time       `gorm:"not null;index"`
	OfferIDs           JSONSlice[uint] `gorm:"type:json"`
}

func (m *PromoCodeModel) ToDomain() domain.PromoCode {
	return domain.PromoCode{
		ID:                 m.ID,
		SchoolID:           m.SchoolID,
		Code:               m.Code,
		DiscountPercentage: m.DiscountPercentage,
		ExpiresAt:          m.ExpiresAt.Unix(),
		OfferIDs:           m.OfferIDs,
	}
}

func (m *PromoCodeModel) FromDomain(p domain.PromoCode) {
	m.ID = p.ID
	m.SchoolID = p.SchoolID
	m.Code = p.Code
	m.DiscountPercentage = p.DiscountPercentage
	m.ExpiresAt = time.Unix(p.ExpiresAt, 0)
	m.OfferIDs = p.OfferIDs
}

// OrderModel 订单 GORM 模型
type OrderModel struct {
	ID           uint              `gorm:"primaryKey;autoIncrement"`
	SchoolID     uint              `gorm:"not null;index"`
	StudentID    uint              `gorm:"not null"`
	StudentName  string            `gorm:"size:255"`
	StudentEmail string            `gorm:"size:255"`
	OfferID      uint              `gorm:"not null"`
	OfferName    string            `gorm:"size:255"`
	PromoID      uint
	PromoCode    string            `gorm:"size:100"`
	CreatedAt    time.Time         `gorm:"not null"`
	Amount       uint              `gorm:"not null"`
	Currency     string            `gorm:"size:10;not null"`
	Status       string            `gorm:"size:50;not null;index"`
	Transactions JSONTransactions  `gorm:"type:json"`
}

// JSONTransactions 交易记录 JSON 类型
type JSONTransactions []domain.Transaction

func (j JSONTransactions) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONTransactions) Scan(value interface{}) error {
	if value == nil {
		*j = []domain.Transaction{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return errors.New("unsupported type for JSONTransactions")
	}
	return json.Unmarshal(bytes, j)
}

func (m *OrderModel) ToDomain() domain.Order {
	return domain.Order{
		ID:       m.ID,
		SchoolID: m.SchoolID,
		Student: domain.StudentInfoShort{
			ID:    m.StudentID,
			Name:  m.StudentName,
			Email: m.StudentEmail,
		},
		Offer: domain.OrderOfferInfo{
			ID:   m.OfferID,
			Name: m.OfferName,
		},
		Promo: domain.OrderPromoInfo{
			ID:   m.PromoID,
			Code: m.PromoCode,
		},
		CreatedAt:    m.CreatedAt.Unix(),
		Amount:       m.Amount,
		Currency:     m.Currency,
		Status:       m.Status,
		Transactions: m.Transactions,
	}
}

func (m *OrderModel) FromDomain(o domain.Order) {
	m.ID = o.ID
	m.SchoolID = o.SchoolID
	m.StudentID = o.Student.ID
	m.StudentName = o.Student.Name
	m.StudentEmail = o.Student.Email
	m.OfferID = o.Offer.ID
	m.OfferName = o.Offer.Name
	m.PromoID = o.Promo.ID
	m.PromoCode = o.Promo.Code
	m.CreatedAt = time.Unix(o.CreatedAt, 0)
	m.Amount = o.Amount
	m.Currency = o.Currency
	m.Status = o.Status
	m.Transactions = o.Transactions
}

// FileModel 文件 GORM 模型
type FileModel struct {
	ID              uint              `gorm:"primaryKey;autoIncrement"`
	SchoolID        uint              `gorm:"not null;index"`
	Type            domain.FileType   `gorm:"size:50;not null;index"`
	ContentType     string            `gorm:"size:100"`
	Name            string            `gorm:"size:255;not null"`
	Size            int64             `gorm:"not null"`
	Status          domain.FileStatus `gorm:"not null;default:0;index"`
	UploadStartedAt time.Time         `gorm:"not null"`
	URL             string            `gorm:"size:500"`
}

func (m *FileModel) ToDomain() domain.File {
	return domain.File{
		ID:              m.ID,
		SchoolID:        m.SchoolID,
		Type:            m.Type,
		ContentType:     m.ContentType,
		Name:            m.Name,
		Size:            m.Size,
		Status:          m.Status,
		UploadStartedAt: m.UploadStartedAt.Unix(),
		URL:             m.URL,
	}
}

func (m *FileModel) FromDomain(f domain.File) {
	m.ID = f.ID
	m.SchoolID = f.SchoolID
	m.Type = f.Type
	m.ContentType = f.ContentType
	m.Name = f.Name
	m.Size = f.Size
	m.Status = f.Status
	m.UploadStartedAt = time.Unix(f.UploadStartedAt, 0)
	m.URL = f.URL
}

// SurveyResultModel 问卷结果 GORM 模型
type SurveyResultModel struct {
	ID          uint              `gorm:"primaryKey;autoIncrement"`
	StudentID   uint              `gorm:"not null"`
	StudentName string            `gorm:"size:255"`
	StudentEmail string           `gorm:"size:255"`
	ModuleID    uint              `gorm:"not null;index"`
	SubmittedAt time.Time         `gorm:"not null"`
	Answers     JSONSurveyAnswers `gorm:"type:json"`
}

// JSONSurveyAnswers 问卷答案 JSON 类型
type JSONSurveyAnswers []domain.SurveyAnswer

func (j JSONSurveyAnswers) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONSurveyAnswers) Scan(value interface{}) error {
	if value == nil {
		*j = []domain.SurveyAnswer{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return errors.New("unsupported type for JSONSurveyAnswers")
	}
	return json.Unmarshal(bytes, j)
}

func (m *SurveyResultModel) ToDomain() domain.SurveyResult {
	return domain.SurveyResult{
		ID: m.ID,
		Student: domain.StudentInfoShort{
			ID:    m.StudentID,
			Name:  m.StudentName,
			Email: m.StudentEmail,
		},
		ModuleID:    m.ModuleID,
		SubmittedAt: m.SubmittedAt.Unix(),
		Answers:     m.Answers,
	}
}

func (m *SurveyResultModel) FromDomain(sr domain.SurveyResult) {
	m.ID = sr.ID
	m.StudentID = sr.Student.ID
	m.StudentName = sr.Student.Name
	m.StudentEmail = sr.Student.Email
	m.ModuleID = sr.ModuleID
	m.SubmittedAt = time.Unix(sr.SubmittedAt, 0)
	m.Answers = sr.Answers
}

// SurveyModel 问卷 GORM 模型（存储在 Module 中）
type SurveyModel struct {
	Title     string              `json:"title"`
	Questions JSONSurveyQuestions `json:"questions" gorm:"type:json"`
	Required  bool                `json:"required"`
}

// JSONSurveyQuestions 问卷问题 JSON 类型
type JSONSurveyQuestions []domain.SurveyQuestion

func (j JSONSurveyQuestions) Value() (driver.Value, error) {
	if j == nil {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONSurveyQuestions) Scan(value interface{}) error {
	if value == nil {
		*j = []domain.SurveyQuestion{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return errors.New("unsupported type for JSONSurveyQuestions")
	}
	return json.Unmarshal(bytes, j)
}

func (m *SurveyModel) ToDomain() domain.Survey {
	return domain.Survey{
		Title:     m.Title,
		Questions: m.Questions,
		Required:  m.Required,
	}
}

func (m *SurveyModel) FromDomain(s domain.Survey) {
	m.Title = s.Title
	m.Questions = s.Questions
	m.Required = s.Required
}
