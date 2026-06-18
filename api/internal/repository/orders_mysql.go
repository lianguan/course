package repository

import (
	"context"

	"ultrathreads/internal/domain"
	"ultrathreads/internal/repository/models"
	"gorm.io/gorm"
)

type OrdersRepo struct {
	db *gorm.DB
}

func NewOrdersRepo(db *gorm.DB) *OrdersRepo {
	return &OrdersRepo{db: db}
}

func (r *OrdersRepo) Create(ctx context.Context, order domain.Order) error {
	var model models.OrderModel
	model.FromDomain(order)
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *OrdersRepo) AddTransaction(ctx context.Context, id uint, transaction domain.Transaction) (domain.Order, error) {
	var model models.OrderModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return domain.Order{}, err
	}

	order := model.ToDomain()
	order.Transactions = append(order.Transactions, transaction)
	order.Status = transaction.Status

	model.FromDomain(order)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return domain.Order{}, err
	}

	return model.ToDomain(), nil
}

func (r *OrdersRepo) GetBySchool(ctx context.Context, schoolID uint, query domain.GetOrdersQuery) ([]domain.Order, int64, error) {
	var orderModels []models.OrderModel
	var count int64

	db := r.db.WithContext(ctx).Model(&models.OrderModel{}).Where("school_id = ?", schoolID)

	if query.Search != "" {
		db = db.Where("student_name LIKE ? OR offer_name LIKE ? OR promo_code LIKE ?",
			"%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Search+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.DateFrom != "" {
		db = db.Where("created_at >= ?", query.DateFrom)
	}
	if query.DateTo != "" {
		db = db.Where("created_at <= ?", query.DateTo)
	}

	if err := db.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if query.PaginationQuery.Limit > 0 {
		db = db.Limit(int(query.PaginationQuery.Limit))
	}
	if query.PaginationQuery.Skip > 0 {
		db = db.Offset(int(query.PaginationQuery.Skip))
	}

	err := db.Order("created_at DESC").Find(&orderModels).Error
	if err != nil {
		return nil, 0, err
	}

	orders := make([]domain.Order, len(orderModels))
	for i, m := range orderModels {
		orders[i] = m.ToDomain()
	}
	return orders, count, nil
}

func (r *OrdersRepo) GetByID(ctx context.Context, id uint) (domain.Order, error) {
	var model models.OrderModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		return domain.Order{}, err
	}
	return model.ToDomain(), nil
}

func (r *OrdersRepo) SetStatus(ctx context.Context, id uint, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.OrderModel{}).
		Where("id = ?", id).
		Update("status", status).Error
}
