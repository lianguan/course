package domain

import "time"

const (
	OrderStatusCreated  = "created"  // 已创建
	OrderStatusPaid     = "paid"     // 已支付
	OrderStatusFailed   = "failed"   // 支付失败
	OrderStatusCanceled = "canceled" // 已取消
	OrderStatusOther    = "other"    // 其他
)

type Order struct {
	ID           uint             `gorm:"primaryKey;autoIncrement" json:"id"`       // 订单ID
	SchoolID     uint             `gorm:"not null;index" json:"schoolId"`           // 所属学校ID
	Student      StudentInfoShort `gorm:"serializer:json" json:"student"`           // 学生信息
	Offer        OrderOfferInfo   `gorm:"serializer:json" json:"offer"`             // 优惠信息
	Promo        OrderPromoInfo   `gorm:"serializer:json" json:"promo"`             // 优惠码信息
	CreatedAt    time.Time        `gorm:"not null" json:"createdAt"`                // 创建时间
	Amount       uint             `gorm:"not null" json:"amount"`                   // 订单金额
	Currency     string           `gorm:"size:10;not null" json:"currency"`         // 货币类型
	Status       string           `gorm:"size:50;not null;index" json:"status"`     // 订单状态
	Transactions []Transaction    `gorm:"serializer:json" json:"transactions"`      // 交易记录
}

type OrderOfferInfo struct {
	ID   uint   `json:"id"`   // 优惠ID
	Name string `json:"name"` // 优惠名称
}

type OrderPromoInfo struct {
	ID   uint   `json:"id"`   // 优惠码ID
	Code string `json:"code"` // 优惠码
}

type Transaction struct {
	Status         string    `json:"status"`         // 交易状态
	CreatedAt      time.Time `json:"createdAt"`      // 交易时间
	AdditionalInfo string    `json:"additionalInfo"` // 附加信息
}

// IsCreated 检查订单是否处于创建状态
func (o *Order) IsCreated() bool {
	return o.Status == OrderStatusCreated
}

// IsPaid 检查订单是否已支付
func (o *Order) IsPaid() bool {
	return o.Status == OrderStatusPaid
}

// IsFailed 检查订单是否支付失败
func (o *Order) IsFailed() bool {
	return o.Status == OrderStatusFailed
}

// IsCanceled 检查订单是否已取消
func (o *Order) IsCanceled() bool {
	return o.Status == OrderStatusCanceled
}

// MarkAsPaid 标记订单为已支付
func (o *Order) MarkAsPaid() {
	o.Status = OrderStatusPaid
}

// MarkAsFailed 标记订单为支付失败
func (o *Order) MarkAsFailed() {
	o.Status = OrderStatusFailed
}

// MarkAsCanceled 标记订单为已取消
func (o *Order) MarkAsCanceled() {
	o.Status = OrderStatusCanceled
}

// AddTransaction 添加交易记录
func (o *Order) AddTransaction(transaction Transaction) {
	o.Transactions = append(o.Transactions, transaction)
}

// HasPromo 检查订单是否使用了优惠码
func (o *Order) HasPromo() bool {
	return o.Promo.ID != 0
}

// NewOrder 创建新订单
func NewOrder(schoolID uint, student StudentInfoShort, offer OrderOfferInfo, amount uint, currency string, promo OrderPromoInfo) *Order {
	return &Order{
		SchoolID:     schoolID,
		Student:      student,
		Offer:        offer,
		Promo:        promo,
		Amount:       amount,
		Currency:     currency,
		CreatedAt:    time.Now(),
		Status:       OrderStatusCreated,
		Transactions: make([]Transaction, 0),
	}
}

// CalculateDiscountedPrice 计算折扣后的价格
func CalculateDiscountedPrice(originalPrice uint, discountPercentage int) uint {
	if discountPercentage <= 0 || discountPercentage >= 100 {
		return originalPrice
	}
	return (originalPrice * uint(100-discountPercentage)) / 100
}
