package analytic

import (
	"PETShoP/internal/models"
	"context"
)

type AnalyticsMock struct {
	GetUserOrderHistoryFunc func(ctx context.Context, email string) ([]models.OrderDetail, error)
	GetPopularProductsFunc  func(ctx context.Context) ([]models.PopularProduct, error)
}

func (m *AnalyticsMock) GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error) {
	if m.GetUserOrderHistoryFunc != nil {
		return m.GetUserOrderHistoryFunc(ctx, email)
	}
	return nil, nil
}

func (m *AnalyticsMock) GetPopularProducts(ctx context.Context) ([]models.PopularProduct, error) {
	if m.GetPopularProductsFunc != nil {
		return m.GetPopularProductsFunc(ctx)
	}
	return nil, nil
}
