package order

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"PETShoP/internal/models"
	"PETShoP/internal/storage"

	"github.com/go-chi/chi"
)

// =======================
// Create Order
// =======================

func TestCreateOrder_Success(t *testing.T) {
	mock := &OrderMock{
		CreateOrderFunc: func(order models.Order) (int, error) {
			return 1, nil
		},
	}

	body := `{
		"CustomerID": 1,
		"TotalPrice": 99.99
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateOrder(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestCreateOrder_BadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{"customer_id":`))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &OrderMock{})
	handler.CreateOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateOrder_InvalidCustomerID(t *testing.T) {
	body := `{
		"CustomerID": 0,
		"TotalPrice": 50.0
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &OrderMock{})
	handler.CreateOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateOrder_Fail(t *testing.T) {
	mock := &OrderMock{
		CreateOrderFunc: func(order models.Order) (int, error) {
			return 0, errors.New("db error")
		},
	}

	body := `{
		"CustomerID": 1,
		"TotalPrice": 99.99
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateOrder(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Get Order By ID
// =======================

func TestGetOrderByID_Success(t *testing.T) {
	now := time.Now()
	mock := &OrderMock{
		GetOrderByIDFunc: func(id int) (models.Order, error) {
			return models.Order{
				ID:         id,
				CustomerID: 1,
				TotalPrice: 150.50,
				CreatedAt:  now,
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetOrderByID_EmptyID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders/", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), &OrderMock{})
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetOrderByID_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders/abc", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "abc")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), &OrderMock{})
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetOrderByID_ZeroOrNegativeID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders/0", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "0")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), &OrderMock{})
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetOrderByID_NotFound(t *testing.T) {
	mock := &OrderMock{
		GetOrderByIDFunc: func(id int) (models.Order, error) {
			return models.Order{}, storage.ErrNotFound
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/999", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "999")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestGetOrderByID_Fail(t *testing.T) {
	mock := &OrderMock{
		GetOrderByIDFunc: func(id int) (models.Order, error) {
			return models.Order{}, errors.New("database connection failed")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Get Orders By User Email
// =======================

func TestGetOrdersByUserEmail_Success(t *testing.T) {
	now := time.Now()
	mock := &OrderMock{
		GetOrdersByUserEmailFunc: func(email string) ([]models.Order, error) {
			return []models.Order{
				{ID: 1, CustomerID: 1, TotalPrice: 50.0, CreatedAt: now},
				{ID: 2, CustomerID: 1, TotalPrice: 100.0, CreatedAt: now},
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/user/alice@example.com", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("email", "alice@example.com")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)
	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_EmptyEmail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders/user/", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("email", "")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), &OrderMock{})
	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_InvalidEmail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders/user/invalid-email", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("email", "invalid-email")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), &OrderMock{})
	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_Fail(t *testing.T) {
	mock := &OrderMock{
		GetOrdersByUserEmailFunc: func(email string) ([]models.Order, error) {
			return nil, errors.New("db error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/user/alice@example.com", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("email", "alice@example.com")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)
	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
