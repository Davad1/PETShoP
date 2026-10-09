package order

import (
	"PETShoP/internal/models"
	"PETShoP/internal/storage"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/mail"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Orders interface {
	CreateOrder(ctx context.Context, order models.Order) (int, error)
	GetOrderByID(ctx context.Context, id int) (models.Order, error)
	GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error)
}

type Handler struct {
	log     *slog.Logger
	storage Orders
}

func New(logger *slog.Logger, storage Orders) *Handler {
	return &Handler{
		log:     logger,
		storage: storage,
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.order.CreateOrder"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Creating a new order", slog.String("url", r.URL.String()))

	var order models.Order
	if err := render.DecodeJSON(r.Body, &order); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	if order.CustomerID <= 0 {
		log.Error("invalid customer ID", slog.Int("customer_id", order.CustomerID))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid customer ID",
		})
		return
	}

	orderID, err := h.storage.CreateOrder(r.Context(), order)
	if err != nil {
		log.Error("failed to create order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to create order",
		})
		return
	}

	log.Info("Order created successfully",
		slog.Int("order_id", orderID),
		slog.Int("customer_id", order.CustomerID),
		slog.Float64("total_price", order.TotalPrice),
	)

	order.ID = orderID
	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"status":      "Order created successfully",
		"id":          orderID,
		"order":       order,
		"total_price": order.TotalPrice,
	})
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.order.GetOrderByID"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	IDStr := chi.URLParam(r, "id")
	if IDStr == "" {
		log.Error("empty id")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID is required",
		})
		return
	}

	id, err := strconv.Atoi(IDStr)
	if err != nil {
		log.Error("invalid id", slog.Any("error", err), slog.String("id", IDStr))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid product ID",
		})
		return
	}

	if id <= 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID must be bigger than zero",
		})
		return
	}

	order, err := h.storage.GetOrderByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{
				"error":   "Not found",
				"message": "Order not found",
			})
			return
		}

		log.Error("failed to get order by id", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to get order by ID",
		})
		return
	}

	log.Info("Retrieved order successfully",
		slog.String("url", r.URL.String()),
		slog.Int("id", order.ID),
	)

	render.JSON(w, r, order)
}

func (h *Handler) GetOrdersByUserEmail(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.order.GetOrdersByUserEmail"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	email := r.URL.Query().Get("email")
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
		log.Error("Invalid email format",
			slog.Any("error", err),
			slog.String("email", email),
		)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid email format",
		})
		return
	}

	orders, err := h.storage.GetOrdersByUserEmail(r.Context(), email)
	if err != nil {
		log.Error("failed to get orders by user email", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to get orders by user email",
		})
		return
	}

	if orders == nil {
		orders = make([]models.Order, 0)
	}

	log.Info("Retrieved orders successfully",
		slog.String("url", r.URL.String()),
		slog.String("email", email),
	)

	render.JSON(w, r, orders)
}
