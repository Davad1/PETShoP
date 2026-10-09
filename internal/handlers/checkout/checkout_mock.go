package checkout

import (
	"PETShoP/internal/models"
	"context"
)

type CheckoutMock struct {
	PlaceOrderFunc func(ctx context.Context, userEmail string, items []models.OrderItem) (int, error)
}

func (m *CheckoutMock) PlaceOrder(ctx context.Context, userEmail string, items []models.OrderItem) (int, error) {
	if m.PlaceOrderFunc != nil {
		return m.PlaceOrderFunc(ctx, userEmail, items)
	}

	return 0, nil
}
