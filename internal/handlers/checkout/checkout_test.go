package checkout

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"PETShoP/internal/models"
	"PETShoP/internal/storage"
)

// =======================
// Place Order
// =======================

func TestPlaceOrder_Success(t *testing.T) {
	mock := &CheckoutMock{
		PlaceOrderFunc: func(ctx context.Context, userEmail string, items []models.OrderItem) (int, error) {
			return 1, nil
		},
	}

	body := `{
		"user_email": "alice@example.com",
		"items": [
			{
				"ProductID": 1,
				"Quantity": 2
			},
			{
				"ProductID": 2,
				"Quantity": 1
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestPlaceOrder_BadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(`{"user_email":`))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &CheckoutMock{})
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestPlaceOrder_EmptyEmail(t *testing.T) {
	body := `{
		"user_email": "",
		"items": [
			{
				"ProductID": 1,
				"Quantity": 2
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &CheckoutMock{})
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestPlaceOrder_InvalidEmail(t *testing.T) {
	body := `{
		"user_email": "invalid-email",
		"items": [
			{
				"ProductID": 1,
				"Quantity": 2
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &CheckoutMock{})
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestPlaceOrder_EmptyItems(t *testing.T) {
	body := `{
		"user_email": "alice@example.com",
		"items": []
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &CheckoutMock{})
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestPlaceOrder_InvalidProductID(t *testing.T) {
	body := `{
		"user_email": "alice@example.com",
		"items": [
			{
				"ProductID": 0,
				"Quantity": 2
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &CheckoutMock{})
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestPlaceOrder_InvalidQuantity(t *testing.T) {
	body := `{
		"user_email": "alice@example.com",
		"items": [
			{
				"ProductID": 1,
				"Quantity": 0
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &CheckoutMock{})
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestPlaceOrder_NotFound(t *testing.T) {
	mock := &CheckoutMock{
		PlaceOrderFunc: func(ctx context.Context, userEmail string, items []models.OrderItem) (int, error) {
			return 0, storage.ErrNotFound
		},
	}

	body := `{
		"user_email": "alice@example.com",
		"items": [
			{
				"ProductID": 999,
				"Quantity": 1
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestPlaceOrder_InsufficientStock(t *testing.T) {
	mock := &CheckoutMock{
		PlaceOrderFunc: func(ctx context.Context, userEmail string, items []models.OrderItem) (int, error) {
			return 0, storage.ErrInsufficientStock
		},
	}

	body := `{
		"user_email": "alice@example.com",
		"items": [
			{
				"ProductID": 1,
				"Quantity": 999
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestPlaceOrder_Fail(t *testing.T) {
	mock := &CheckoutMock{
		PlaceOrderFunc: func(ctx context.Context, userEmail string, items []models.OrderItem) (int, error) {
			return 0, errors.New("db error")
		},
	}

	body := `{
		"user_email": "alice@example.com",
		"items": [
			{
				"ProductID": 1,
				"Quantity": 2
			}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
