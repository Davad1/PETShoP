package product

import (
	"context"
	"errors"
	"PETShoP/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"github.com/go-chi/chi"
)

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage — он вернёт один продукт.
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return []models.Product{
				{ID: 1, Name: "Dog Food"},
			}, nil
		},
	}

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), mock)

	// Вызываем метод GetAllProducts, который является http.HandlerFunc
	handler.GetAllProducts(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
func TestGetAllProducts_Error(t *testing.T) {
	// Мокаем storage — он будет возвращать ошибку
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return nil, errors.New("DB error")
		},
	}

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllProducts(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Create Product
// =======================

func TestCreateProduct_Success(t *testing.T) {
	mock := &ProductsMock{
		CreateProductFunc: func(
			ctx context.Context,
			product models.Product,
		) (int, error) {
			return 1, nil
		},
	}

	body := `{
		"name": "Dog Food",
		"price": 10,
		"stock": 5
	}`

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	handler.CreateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestCreateProduct_BadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"name":`))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &ProductsMock{})

	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateProduct_Fail(t *testing.T) {
	mock := &ProductsMock{
		CreateProductFunc: func(
			ctx context.Context,
			product models.Product,
		) (int, error) {
			return 0, errors.New("db error")
		},
	}

	body := `{
		"name": "Dog Food",
		"price": 10,
		"stock": 5
	}`

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	handler.CreateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) {
	mock := &ProductsMock{
		UpdateProductFunc: func(
			ctx context.Context,
			product models.Product,
		) error {
			return nil
		},
	}
	body := `{
		"name": "Updated Dog Food",
		"price": 15,
		"stock": 10
	}`

	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(body))
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestUpdateProduct_BadRequest(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPut,
		"/products/1",
		strings.NewReader(`{"name":`),
	)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), &ProductsMock{})
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestUpdateProduct_Fail(t *testing.T) {
	mock := &ProductsMock{
		UpdateProductFunc: func(
			ctx context.Context,
			product models.Product,
		) error {
			return errors.New("db error")
		},
	}
	body := `{
		"name": "Updated Dog Food",
		"price": 15,
		"stock": 10
	}`

	req := httptest.NewRequest(http.MethodPut,"/products/1", strings.NewReader(body))

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Delete Product
// =======================

func TestDeleteProduct_Success(t *testing.T) {
	mock := &ProductsMock{
		DeleteProductFunc: func(
			ctx context.Context,
			id int,
		) error {
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	handler.DeleteProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/products/", nil)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()

	handler := New(slog.Default(), &ProductsMock{})

	handler.DeleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestDeleteProduct_Fail(t *testing.T) {
	mock := &ProductsMock{
		DeleteProductFunc: func(
			ctx context.Context,
			id int,
		) error {
			return errors.New("db error")
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", "1")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)

	handler.DeleteProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}