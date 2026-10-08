package order

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"PETShoP/internal/models"

	"github.com/go-chi/chi"
)

// =======================
// Add Order Item
// =======================

func TestAddOrderItem_Success(t *testing.T) {
	mock := &OrderMock{
		AddOrderItemFunc: func(orderItem models.OrderItem) error {
			return nil
		},
	}

	body := `{
		"ProductID": 1,
		"Quantity": 2
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders/1/items", strings.NewReader(body))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := NewItem(slog.Default(), mock)
	handler.AddOrderItem(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestAddOrderItem_BadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/orders/1/items", strings.NewReader(`{"ProductID":`))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := NewItem(slog.Default(), &OrderMock{})
	handler.AddOrderItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestAddOrderItem_Fail(t *testing.T) {
	mock := &OrderMock{
		AddOrderItemFunc: func(orderItem models.OrderItem) error {
			return errors.New("db error")
		},
	}

	body := `{
		"ProductID": 1,
		"Quantity": 2
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders/1/items", strings.NewReader(body))

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()

	handler := NewItem(slog.Default(), mock)
	handler.AddOrderItem(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Get Order Items By Order ID
// =======================

func TestGetOrderItemsByOrderID_Success(t *testing.T) {
	mock := &OrderMock{
		GetOrderItemsByOrderIDFunc: func(orderID int) ([]models.OrderItem, error) {
			return []models.OrderItem{
				{ID: 1, OrderID: orderID, ProductID: 1, Quantity: 2},
				{ID: 2, OrderID: orderID, ProductID: 2, Quantity: 1},
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/1/items", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := NewItem(slog.Default(), mock)
	handler.GetOrderItemsByOrderID(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetOrderItemsByOrderID_BadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/orders/invalid/items", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "invalid")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := NewItem(slog.Default(), &OrderMock{})
	handler.GetOrderItemsByOrderID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetOrderItemsByOrderID_Fail(t *testing.T) {
	mock := &OrderMock{
		GetOrderItemsByOrderIDFunc: func(orderID int) ([]models.OrderItem, error) {
			return nil, errors.New("db error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/1/items", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := NewItem(slog.Default(), mock)
	handler.GetOrderItemsByOrderID(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
