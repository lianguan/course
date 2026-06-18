package dto

import "ultrathreads/internal/domain"

// ===== Offer DTOs =====

type CreateOfferRequest struct {
	Name          string            `json:"name" binding:"required"`
	Description   string            `json:"description"`
	Benefits      []string          `json:"benefits"`
	Price         PriceRequest      `json:"price" binding:"required"`
	Packages      []uint            `json:"packages"`
	PaymentMethod PaymentMethodReq  `json:"paymentMethod" binding:"required"`
}

type UpdateOfferRequest struct {
	Name          *string           `json:"name"`
	Description   *string           `json:"description"`
	Benefits      []string          `json:"benefits"`
	Price         *PriceRequest     `json:"price"`
	Packages      []uint            `json:"packages"`
	PaymentMethod *PaymentMethodReq `json:"paymentMethod"`
}

type PriceRequest struct {
	Value    uint   `json:"value" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

type PaymentMethodReq struct {
	UsesProvider bool   `json:"usesProvider"`
	Provider     string `json:"provider"`
}

type OfferResponse struct {
	ID            uint              `json:"id"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Benefits      []string          `json:"benefits"`
	Price         PriceResponse     `json:"price"`
	PaymentMethod PaymentMethodResp `json:"paymentMethod"`
	Packages      []uint            `json:"packages"`
}

type PriceResponse struct {
	Value    uint   `json:"value"`
	Currency string `json:"currency"`
}

type PaymentMethodResp struct {
	UsesProvider bool   `json:"usesProvider"`
	Provider     string `json:"provider"`
}

func OfferToResponse(o domain.Offer) OfferResponse {
	return OfferResponse{
		ID:          o.ID,
		Name:        o.Name,
		Description: o.Description,
		Benefits:    o.Benefits,
		Price: PriceResponse{
			Value:    o.Price.Value,
			Currency: o.Price.Currency,
		},
		PaymentMethod: PaymentMethodResp{
			UsesProvider: o.PaymentMethod.UsesProvider,
			Provider:     o.PaymentMethod.Provider,
		},
		Packages: o.PackageIDs,
	}
}

func OffersToResponse(offers []domain.Offer) []OfferResponse {
	res := make([]OfferResponse, len(offers))
	for i, o := range offers {
		res[i] = OfferToResponse(o)
	}
	return res
}

// ===== PromoCode DTOs =====

type CreatePromoCodeRequest struct {
	Code               string `json:"code" binding:"required"`
	DiscountPercentage int    `json:"discountPercentage" binding:"required,min=1,max=100"`
	ExpiresAt          int64  `json:"expiresAt" binding:"required"`
	OfferIDs           []uint `json:"offerIds"`
}

type UpdatePromoCodeRequest struct {
	Code               *string `json:"code"`
	DiscountPercentage *int    `json:"discountPercentage"`
	ExpiresAt          *int64  `json:"expiresAt"`
	OfferIDs           []uint  `json:"offerIds"`
}

type PromoCodeResponse struct {
	ID                 uint   `json:"id"`
	Code               string `json:"code"`
	DiscountPercentage int    `json:"discountPercentage"`
	ExpiresAt          int64  `json:"expiresAt"`
	OfferIDs           []uint `json:"offerIds"`
}

func PromoCodeToResponse(p domain.PromoCode) PromoCodeResponse {
	return PromoCodeResponse{
		ID:                 p.ID,
		Code:               p.Code,
		DiscountPercentage: p.DiscountPercentage,
		ExpiresAt:          p.ExpiresAt,
		OfferIDs:           p.OfferIDs,
	}
}

func PromoCodesToResponse(promos []domain.PromoCode) []PromoCodeResponse {
	res := make([]PromoCodeResponse, len(promos))
	for i, p := range promos {
		res[i] = PromoCodeToResponse(p)
	}
	return res
}

// ===== Order DTOs =====

type OrderResponse struct {
	ID           uint              `json:"id"`
	Student      StudentInfoResp   `json:"student"`
	Offer        OfferInfoResp     `json:"offer"`
	Promo        PromoInfoResp     `json:"promo"`
	CreatedAt    int64             `json:"createdAt"`
	Amount       uint              `json:"amount"`
	Currency     string            `json:"currency"`
	Status       string            `json:"status"`
	Transactions []TransactionResp `json:"transactions"`
}

type StudentInfoResp struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type OfferInfoResp struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PromoInfoResp struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
}

type TransactionResp struct {
	Status         string `json:"status"`
	CreatedAt      int64  `json:"createdAt"`
	AdditionalInfo string `json:"additionalInfo"`
}

func OrderToResponse(o domain.Order) OrderResponse {
	transactions := make([]TransactionResp, len(o.Transactions))
	for i, t := range o.Transactions {
		transactions[i] = TransactionResp{
			Status:         t.Status,
			CreatedAt:      t.CreatedAt,
			AdditionalInfo: t.AdditionalInfo,
		}
	}

	return OrderResponse{
		ID: o.ID,
		Student: StudentInfoResp{
			ID:    o.Student.ID,
			Name:  o.Student.Name,
			Email: o.Student.Email,
		},
		Offer: OfferInfoResp{
			ID:   o.Offer.ID,
			Name: o.Offer.Name,
		},
		Promo: PromoInfoResp{
			ID:   o.Promo.ID,
			Code: o.Promo.Code,
		},
		CreatedAt:    o.CreatedAt,
		Amount:       o.Amount,
		Currency:     o.Currency,
		Status:       o.Status,
		Transactions: transactions,
	}
}

func OrdersToResponse(orders []domain.Order) []OrderResponse {
	res := make([]OrderResponse, len(orders))
	for i, o := range orders {
		res[i] = OrderToResponse(o)
	}
	return res
}
