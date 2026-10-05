package user

import (
	"context"
	"PETShoP/internal/models"
	"log/slog"
	"net/http"
	"net/mail"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Users interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, u models.User) error
}

type Handler struct {
	log *slog.Logger
	storage Users
}

func New(logger *slog.Logger, storage Users) *Handler {
	return &Handler{
		log:     logger,
		storage: storage,
	}
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.user.GetAllUsers"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	users, err := h.storage.GetAllUsers(r.Context())
	if err != nil {
		log.Error("Failed to get all users", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve users",
		})
		return
	}

	render.JSON(w, r, users)
}


func (h* Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.user.GetUserByEmail"

	 log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	email := chi.URLParam(r, "email")
	if email == "" {
		log.Error("Email parameter is missing")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Email parameter is required",
		})
		return
	}

	addr, err := mail.ParseAddress(email)
	if err != nil || addr.Address != email {
		log.Error("Invalid email format", slog.Any("error", err), slog.String("email", email))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid email format",
		})
		return
	}
	 

	user, err := h.storage.GetUserByEmail(r.Context(), email)
	if err != nil{
		log.Error("Failed to get user by email", slog.Any("error", err))
		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, map[string]string{
			"error":   "Not found",
			"message": "User not found",
		})
		return
	}

	log.Info("Retrieved user successfully",
		slog.String("url", r.URL.String()),
		slog.Int("id", user.ID),
	)

	render.JSON(w, r, user)

}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.user.CreateUser"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Creating new user", slog.String("url", r.URL.String()))

	var user models.User
	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.Error("Failed to decode user data", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid user data",
		})
		return
	}

	if user.Name == "" {
		log.Error("user name is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "User name is required",
		})
		return
	}

	if user.Email == "" {
		log.Error("user email is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "User email is required",
		})
		return
	}

	err := h.storage.CreateUser(r.Context(), user)
	if err != nil {
		log.Error("Failed to create user", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to create user",
		})
		return
	}

	log.Info("User created successfully",
		slog.String("name", user.Name),
		slog.String("email", user.Email),
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, map[string]interface{}{
		"status":  "User created successfully",
		"user": user,
	})
}