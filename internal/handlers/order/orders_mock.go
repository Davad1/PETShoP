package order

import (
	"PETShoP/internal/models"
	"context"
)

type OrderMock struct {
	CreateOrderFunc            func(ctx context.Context, order models.Order) (int, error)
	GetOrderByIDFunc           func(ctx context.Context, id int) (models.Order, error)
	GetOrdersByUserEmailFunc   func(ctx context.Context, email string) ([]models.Order, error)
	AddOrderItemFunc           func(ctx context.Context, orderItem models.OrderItem) error
	GetOrderItemsByOrderIDFunc func(ctx context.Context, orderID int) ([]models.OrderItem, error)
}

func (m *OrderMock) CreateOrder(ctx context.Context, order models.Order) (int, error) {
	if m.CreateOrderFunc != nil {
		return m.CreateOrderFunc(ctx, order)
	}
	return 0, nil
}

func (m *OrderMock) GetOrderByID(ctx context.Context, id int) (models.Order, error) {
	if m.GetOrderByIDFunc != nil {
		return m.GetOrderByIDFunc(ctx, id)
	}
	return models.Order{}, nil
}

func (m *OrderMock) GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error) {
	if m.GetOrdersByUserEmailFunc != nil {
		return m.GetOrdersByUserEmailFunc(ctx, email)
	}
	return []models.Order{}, nil
}

func (m *OrderMock) AddOrderItem(ctx context.Context, orderItem models.OrderItem) error {
	if m.AddOrderItemFunc != nil {
		return m.AddOrderItemFunc(ctx, orderItem)
	}
	return nil
}

func (m *OrderMock) GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	if m.GetOrderItemsByOrderIDFunc != nil {
		return m.GetOrderItemsByOrderIDFunc(ctx, orderID)
	}
	return []models.OrderItem{}, nil
}
