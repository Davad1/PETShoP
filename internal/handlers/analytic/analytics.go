package analytic

import (
	"log/slog"
	"net/http"
	"net/mail"

	"PETShoP/internal/models"
	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/render"
)

type Analytics interface {
	GetUserOrderHistory(email string) ([]models.OrderDetail, error)
	GetPopularProducts() ([]models.PopularProduct, error)
}

type Handler struct {
	log     *slog.Logger
	storage Analytics
}

func New(log *slog.Logger, storage Analytics) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) GetUserOrderHistory(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.analytics.GetUserOrderHistory"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("getting user order history", slog.String("url", r.URL.String()))

	email := r.URL.Query().Get("email")
	if email == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Email is required",
		})
		return
	}

	if _, err := mail.ParseAddress(email); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid email",
		})
		return
	}

	history, err := h.storage.GetUserOrderHistory(email)
	if err != nil {
		log.Error("failed to get user order history",
			slog.String("email", email),
			slog.Any("error", err),
		)

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to get user order history",
		})
		return
	}

	log.Info("user order history retrieved successfully",
		slog.String("email", email),
		slog.Int("orders_count", len(history)),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"status":  "User order history retrieved successfully",
		"history": history,
	})
}

func (h *Handler) GetPopularProducts(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.analytics.GetPopularProducts"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("getting popular products", slog.String("url", r.URL.String()))

	products, err := h.storage.GetPopularProducts()
	if err != nil {
		log.Error("failed to get popular products",
			slog.Any("error", err),
		)

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to get popular products",
		})
		return
	}

	log.Info("popular products retrieved successfully",
		slog.Int("products_count", len(products)),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"status":   "Popular products retrieved successfully",
		"products": products,
	})
}
