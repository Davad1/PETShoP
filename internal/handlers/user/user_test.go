package user

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
// Get All Users
// =======================

func TestGetAllUsers_Success(t *testing.T) {
	mock := &UsersMock{
		GetAllUsersFunc: func(ctx context.Context) ([]models.User, error) {
			return []models.User{
				{ID: 1, Name: "Alice", Email: "alice@example.com"},
				{ID: 2, Name: "Bob", Email: "bob@example.com"},
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetAllUsers_Error(t *testing.T) {
	mock := &UsersMock{
		GetAllUsersFunc: func(ctx context.Context) ([]models.User, error) {
			return nil, errors.New("db error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllUsers(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

// =======================
// Get User By Email
// =======================

func TestGetUserByEmail_Success(t *testing.T) {
	mock := &UsersMock{
		GetUserByEmailFunc: func(ctx context.Context, email string) (models.User, error) {
			return models.User{ID: 1, Name: "Alice", Email: email}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users/alice@example.com", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("email", "alice@example.com")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)
	handler.GetUserByEmail(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestGetUserByEmail_MissingEmail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("email", "")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), &UsersMock{})
	handler.GetUserByEmail(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetUserByEmail_InvalidEmail(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/invalid-email", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("email", "invalid-email")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), &UsersMock{})
	handler.GetUserByEmail(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	mock := &UsersMock{
		GetUserByEmailFunc: func(ctx context.Context, email string) (models.User, error) {
			return models.User{}, errors.New("user not found")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users/notfound@example.com", nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("email", "notfound@example.com")

	req = req.WithContext(
		context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx),
	)

	w := httptest.NewRecorder()
	handler := New(slog.Default(), mock)
	handler.GetUserByEmail(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

// =======================
// Create User
// =======================

func TestCreateUser_Success(t *testing.T) {
	mock := &UsersMock{
		CreateUserFunc: func(ctx context.Context, u models.User) error {
			return nil
		},
	}

	body := `{
		"name": "Alice",
		"email": "alice@example.com"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateUser(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestCreateUser_BadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":`))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &UsersMock{})
	handler.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateUser_EmptyName(t *testing.T) {
	body := `{
		"name": "",
		"email": "alice@example.com"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &UsersMock{})
	handler.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateUser_EmptyEmail(t *testing.T) {
	body := `{
		"name": "Alice",
		"email": ""
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), &UsersMock{})
	handler.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestCreateUser_Fail(t *testing.T) {
	mock := &UsersMock{
		CreateUserFunc: func(ctx context.Context, u models.User) error {
			return errors.New("db error")
		},
	}

	body := `{
		"name": "Alice",
		"email": "alice@example.com"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateUser(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}
