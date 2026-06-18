package domain

const (
	OrderStatusCreated  = "created"  // 已创建
	OrderStatusPaid     = "paid"     // 已支付
	OrderStatusFailed   = "failed"   // 支付失败
	OrderStatusCanceled = "canceled" // 已取消
	OrderStatusOther    = "other"    // 其他
)

// Order 订单实体
type Order struct {
	ID           uint             // 订单ID
	SchoolID     uint             // 所属学校ID
	Student      StudentInfoShort // 学生信息
	Offer        OrderOfferInfo   // 优惠信息
	Promo        OrderPromoInfo   // 优惠码信息
	CreatedAt    int64            // 创建时间（Unix 时间戳）
	Amount       uint             // 订单金额
	Currency     string           // 货币类型
	Status       string           // 订单状态
	Transactions []Transaction    // 交易记录
}

// OrderOfferInfo 订单优惠信息值对象
type OrderOfferInfo struct {
	ID   uint   // 优惠ID
	Name string // 优惠名称
}

// OrderPromoInfo 订单优惠码信息值对象
type OrderPromoInfo struct {
	ID   uint   // 优惠码ID
	Code string // 优惠码
}

// Transaction 交易记录值对象
type Transaction struct {
	Status         string // 交易状态
	CreatedAt      int64  // 交易时间（Unix 时间戳）
	AdditionalInfo string // 附加信息
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
func NewOrder(schoolID uint, student StudentInfoShort, offer OrderOfferInfo, amount uint, currency string, promo OrderPromoInfo, now int64) *Order {
	return &Order{
		SchoolID:     schoolID,
		Student:      student,
		Offer:        offer,
		Promo:        promo,
		Amount:       amount,
		Currency:     currency,
		CreatedAt:    now,
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
