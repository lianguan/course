package domain

import "errors"

const (
	PaymentProviderFondy = "fondy"
)

var (
	ErrPaymentProviderNotUsed = errors.New("payment provider is disabled for current offer")
	ErrUnknownPaymentProvider = errors.New("payment provider is not supported")
)

// Offer 优惠实体
type Offer struct {
	ID            uint          // 优惠ID
	Name          string        // 优惠名称
	Description   string        // 优惠描述
	Benefits      []string      // 权益列表
	SchoolID      uint          // 所属学校ID
	PackageIDs    []uint        // 包含套餐ID列表
	Price         Price         // 价格信息
	PaymentMethod PaymentMethod // 支付方式
}

// Price 价格值对象
type Price struct {
	Value    uint   // 价格数值
	Currency string // 货币类型
}

// PaymentMethod 支付方式值对象
type PaymentMethod struct {
	UsesProvider bool   // 是否使用支付提供商
	Provider     string // 支付提供商名称
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
