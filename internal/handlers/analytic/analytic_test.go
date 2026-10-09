package analytic

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"PETShoP/internal/models"
)

func TestGetUserOrderHistory_Success(t *testing.T) {
	mock := &AnalyticsMock{
		GetUserOrderHistoryFunc: func(ctx context.Context, email string) ([]models.OrderDetail, error) {
			return []models.OrderDetail{
				{
					OrderID:           1,
					ProductID:         1,
					ProductName:       "Dog Food",
					Quantity:          2,
					TotalPrice:        21,
					TransactionStatus: "completed",
					CreatedAt:         time.Now(),
				},
			}, nil
		},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/history?email=alice@example.com",
		nil,
	)

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetUserOrderHistory(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetUserOrderHistory_EmptyEmail(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/users/history",
		nil,
	)

	w := httptest.NewRecorder()

	handler := New(slog.Default(), &AnalyticsMock{})
	handler.GetUserOrderHistory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetUserOrderHistory_InvalidEmail(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodGet,
		"/users/history?email=invalid-email",
		nil,
	)

	w := httptest.NewRecorder()

	handler := New(slog.Default(), &AnalyticsMock{})
	handler.GetUserOrderHistory(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetUserOrderHistory_Fail(t *testing.T) {
	mock := &AnalyticsMock{
		GetUserOrderHistoryFunc: func(ctx context.Context, email string) ([]models.OrderDetail, error) {
			return nil, errors.New("db error")
		},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/history?email=alice@example.com",
		nil,
	)

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetUserOrderHistory(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestGetPopularProducts_Success(t *testing.T) {
	mock := &AnalyticsMock{
		GetPopularProductsFunc: func(ctx context.Context) ([]models.PopularProduct, error) {
			return []models.PopularProduct{
				{
					ProductID: 1,
					Name:      "Dog Food",
					TotalSold: 5,
				},
			}, nil
		},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/products/popular",
		nil,
	)

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetPopularProducts(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetPopularProducts_Fail(t *testing.T) {
	mock := &AnalyticsMock{
		GetPopularProductsFunc: func(ctx context.Context) ([]models.PopularProduct, error) {
			return nil, errors.New("db error")
		},
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/products/popular",
		nil,
	)

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetPopularProducts(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
