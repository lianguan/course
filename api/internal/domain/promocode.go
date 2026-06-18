package domain

import "time"

type PromoCode struct {
	ID                 uint      `gorm:"primaryKey;autoIncrement" json:"id"`        // 优惠码ID
	SchoolID           uint      `gorm:"not null;index" json:"schoolId"`            // 所属学校ID
	Code               string    `gorm:"size:100;not null;uniqueIndex" json:"code"` // 优惠码
	DiscountPercentage int       `gorm:"not null" json:"discountPercentage"`        // 折扣百分比
	ExpiresAt          time.Time `gorm:"not null;index" json:"expiresAt"`           // 过期时间
	OfferIDs           []uint    `gorm:"serializer:json" json:"offerIds"`           // 适用优惠ID列表
}

type UpdatePromoCodeInput struct {
	ID                 uint      // 优惠码ID
	SchoolID           uint      // 学校ID
	Code               string    // 优惠码
	DiscountPercentage int       // 折扣百分比
	ExpiresAt          time.Time // 过期时间
	OfferIDs           []uint    // 适用优惠ID列表
}

// IsExpired 检查优惠码是否已过期
func (p *PromoCode) IsExpired() bool {
	return p.ExpiresAt.Unix() < time.Now().Unix()
}

// IsValid 检查优惠码是否有效（未过期且折扣有效）
func (p *PromoCode) IsValid() bool {
	return !p.IsExpired() && p.DiscountPercentage > 0 && p.DiscountPercentage < 100
}

// AppliesToOffer 检查优惠码是否适用于指定优惠
func (p *PromoCode) AppliesToOffer(offerID uint) bool {
	for _, id := range p.OfferIDs {
		if id == offerID {
			return true
		}
	}
	return false
}

// CalculateDiscountedPrice 计算折扣后的价格
func (p *PromoCode) CalculateDiscountedPrice(originalPrice uint) uint {
	return CalculateDiscountedPrice(originalPrice, p.DiscountPercentage)
}

// CreatePromoCodeInput 优惠码创建输入（Service 层使用）
type CreatePromoCodeInput struct {
	SchoolID           uint
	Code               string
	DiscountPercentage int
	ExpiresAt          time.Time
	OfferIDs           []uint
}

// NewPromoCode 创建新优惠码
func NewPromoCode(schoolID uint, code string, discountPercentage int, expiresAt time.Time, offerIDs []uint) *PromoCode {
	return &PromoCode{
		SchoolID:           schoolID,
		Code:               code,
		DiscountPercentage: discountPercentage,
		ExpiresAt:          expiresAt,
		OfferIDs:           offerIDs,
	}
}
