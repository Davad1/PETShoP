package checkout

import (
	"PETShoP/internal/models"
)

type CheckoutMock struct {
	PlaceOrderFunc func(userEmail string, items []models.OrderItem) (int, error)
}

func (m *CheckoutMock) PlaceOrder(userEmail string, items []models.OrderItem) (int, error) {
	if m.PlaceOrderFunc != nil {
		return m.PlaceOrderFunc(userEmail, items)
	}

	return 0, nil
}
