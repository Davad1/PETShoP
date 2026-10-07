package order

import (
	"PETShoP/internal/models"
)

type OrderMock struct {
	CreateOrderFunc            func(order models.Order) (int, error)
	GetOrderByIDFunc           func(id int) (models.Order, error)
	GetOrdersByUserEmailFunc   func(email string) ([]models.Order, error)
	AddOrderItemFunc           func(orderItem models.OrderItem) error
	GetOrderItemsByOrderIDFunc func(orderID int) ([]models.OrderItem, error)
}

func (m *OrderMock) CreateOrder(order models.Order) (int, error) {
	if m.CreateOrderFunc != nil {
		return m.CreateOrderFunc(order)
	}
	return 0, nil
}

func (m *OrderMock) GetOrderByID(id int) (models.Order, error) {
	if m.GetOrderByIDFunc != nil {
		return m.GetOrderByIDFunc(id)
	}
	return models.Order{}, nil
}

func (m *OrderMock) GetOrdersByUserEmail(email string) ([]models.Order, error) {
	if m.GetOrdersByUserEmailFunc != nil {
		return m.GetOrdersByUserEmailFunc(email)
	}
	return []models.Order{}, nil
}

func (m *OrderMock) AddOrderItem(orderItem models.OrderItem) error {
	if m.AddOrderItemFunc != nil {
		return m.AddOrderItemFunc(orderItem)
	}
	return nil
}

func (m *OrderMock) GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error) {
	if m.GetOrderItemsByOrderIDFunc != nil {
		return m.GetOrderItemsByOrderIDFunc(orderID)
	}
	return []models.OrderItem{}, nil
}
