package checkout

import (
	"errors"
	"log/slog"
	"net/http"
	"net/mail"

	"PETShoP/internal/models"
	"PETShoP/internal/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Checkout interface {
	PlaceOrder(userEmail string, items []models.OrderItem) (int, error)
}

type Handler struct {
	log     *slog.Logger
	storage Checkout
}

func New(logger *slog.Logger, storage Checkout) *Handler {
	return &Handler{
		log:     logger,
		storage: storage,
	}
}

type CheckoutRequest struct {
	UserEmail string             `json:"user_email"`
	Items     []models.OrderItem `json:"items"`
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.checkout.PlaceOrder"
	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req CheckoutRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode checkout request", slog.Any("error", err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	if req.UserEmail == "" {
		log.Error("user email is empty")

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "User email is required",
		})
		return
	}

	addr, err := mail.ParseAddress(req.UserEmail)
	if err != nil || addr.Address != req.UserEmail {
		log.Error("invalid email format",
			slog.Any("error", err),
			slog.String("email", req.UserEmail),
		)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid email format",
		})
		return
	}

	if len(req.Items) == 0 {
		log.Error("checkout items are empty")

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order items are required",
		})
		return
	}
	for _, item := range req.Items {
		if item.ProductID <= 0 || item.Quantity <= 0 {
			log.Error("invalid order item",
				slog.Int("product_id", item.ProductID),
				slog.Int("quantity", item.Quantity),
			)

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error":   "Bad request",
				"message": "Product ID and quantity must be greater than zero",
			})
			return
		}
	}

	orderID, err := h.storage.PlaceOrder(req.UserEmail, req.Items)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Error("resource not found", slog.Any("error", err))

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{
				"error":   "Not found",
				"message": "User or product not found",
			})
			return
		}

		if errors.Is(err, storage.ErrInsufficientStock) {
			log.Error("insufficient stock", slog.Any("error", err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error":   "Bad request",
				"message": "Insufficient product stock",
			})
			return
		}

		log.Error("failed to place order", slog.Any("error", err))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to place order",
		})
		return
	}

	log.Info("order placed successfully",
		slog.Int("order_id", orderID),
		slog.String("email", req.UserEmail),
	)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"status":   "Order placed successfully",
		"order_id": orderID,
	})
}
