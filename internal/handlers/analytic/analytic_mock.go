package analytic

import (
	"PETShoP/internal/models"
)

type AnalyticsMock struct {
	GetUserOrderHistoryFunc func(email string) ([]models.OrderDetail, error)
	GetPopularProductsFunc  func() ([]models.PopularProduct, error)
}

func (m *AnalyticsMock) GetUserOrderHistory(email string) ([]models.OrderDetail, error) {
	if m.GetUserOrderHistoryFunc != nil {
		return m.GetUserOrderHistoryFunc(email)
	}
	return nil, nil
}

func (m *AnalyticsMock) GetPopularProducts() ([]models.PopularProduct, error) {
	if m.GetPopularProductsFunc != nil {
		return m.GetPopularProductsFunc()
	}
	return nil, nil
}
