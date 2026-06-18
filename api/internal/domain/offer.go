package domain

import "errors"

const (
	PaymentProviderFondy = "fondy"
)

var (
	ErrPaymentProviderNotUsed = errors.New("payment provider is disabled for current offer")
	ErrUnknownPaymentProvider = errors.New("payment provider is not supported")
)

type Offer struct {
	ID            uint          `gorm:"primaryKey;autoIncrement" json:"id"`                          // 优惠ID
	Name          string        `gorm:"size:255;not null" json:"name"`                               // 优惠名称
	Description   string        `gorm:"type:text" json:"description"`                                // 优惠描述
	Benefits      []string      `gorm:"serializer:json" json:"benefits"`                             // 权益列表
	SchoolID      uint          `gorm:"not null;index" json:"schoolId"`                              // 所属学校ID
	PackageIDs    []uint        `gorm:"serializer:json" json:"packages"`                             // 包含套餐ID列表
	Price         Price         `gorm:"embedded;embeddedPrefix:price_" json:"price"`                 // 价格信息
	PaymentMethod PaymentMethod `gorm:"embedded;embeddedPrefix:paymentMethod_" json:"paymentMethod"` // 支付方式
}

type Price struct {
	Value    uint   `gorm:"not null;default:0" json:"value"`                // 价格数值
	Currency string `gorm:"size:10;not null;default:'USD'" json:"currency"` // 货币类型
}

type PaymentMethod struct {
	UsesProvider bool   `gorm:"not null;default:false" json:"usesProvider"` // 是否使用支付提供商
	Provider     string `gorm:"size:50" json:"provider"`                    // 支付提供商名称
}

func (pm PaymentMethod) Validate() error {
	switch pm.Provider {
	case PaymentProviderFondy:
		return nil
	default:
		return errors.New("unknown payment provider")
	}
}

// ContainsPackage 检查优惠是否包含指定套餐
func (o *Offer) ContainsPackage(packageID uint) bool {
	for _, id := range o.PackageIDs {
		if id == packageID {
			return true
		}
	}
	return false
}

// GetPrice 获取价格
func (o *Offer) GetPrice() uint {
	return o.Price.Value
}

// GetCurrency 获取货币类型
func (o *Offer) GetCurrency() string {
	return o.Price.Currency
}

// UsesPaymentProvider 检查是否使用支付提供商
func (o *Offer) UsesPaymentProvider() bool {
	return o.PaymentMethod.UsesProvider
}

// ValidatePaymentMethod 验证支付方式
func (o *Offer) ValidatePaymentMethod() error {
	if o.PaymentMethod.UsesProvider {
		return o.PaymentMethod.Validate()
	}
	return nil
}

// ValidatePayment 验证更新输入的支付方式
func (i *UpdateOfferInput) ValidatePayment() error {
	if i.PaymentMethod == nil {
		return nil
	}

	if !i.PaymentMethod.UsesProvider {
		return nil
	}

	return i.PaymentMethod.Validate()
}

// CreateOfferInput 优惠创建输入（Service 层使用）
type CreateOfferInput struct {
	Name          string
	Description   string
	Benefits      []string
	SchoolID      uint
	Price         Price
	Packages      []uint
	PaymentMethod PaymentMethod
}

// NewOffer 创建新优惠
func NewOffer(schoolID uint, name, description string, benefits []string, price Price, paymentMethod PaymentMethod, packageIDs []uint) *Offer {
	return &Offer{
		SchoolID:      schoolID,
		Name:          name,
		Description:   description,
		Benefits:      benefits,
		Price:         price,
		PaymentMethod: paymentMethod,
		PackageIDs:    packageIDs,
	}
}
