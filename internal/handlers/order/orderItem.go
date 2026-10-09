package order

import (
	"PETShoP/internal/models"
	"PETShoP/internal/storage"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type OrdersItem interface {
	AddOrderItem(ctx context.Context, orderItem models.OrderItem) error
	GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error)
}

type HandlerItem struct {
	log     *slog.Logger
	storage OrdersItem
}

func NewItem(logger *slog.Logger, storage OrdersItem) *HandlerItem {
	return &HandlerItem{
		log:     logger,
		storage: storage,
	}
}

func (h *HandlerItem) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.order.AddOrderItem"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Adding a new order item", slog.String("url", r.URL.String()))

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID is required",
		})
		return
	}

	orderID, err := strconv.Atoi(idStr)
	if err != nil || orderID <= 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID must be a positive number",
		})
		return
	}

	var orderItem models.OrderItem
	if err := render.DecodeJSON(r.Body, &orderItem); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid JSON payload",
		})
		return
	}

	orderItem.OrderID = orderID

	if orderItem.ProductID <= 0 || orderItem.Quantity <= 0 {
		log.Error("invalid order item",
			slog.Int("order_id", orderItem.OrderID),
			slog.Int("product_id", orderItem.ProductID),
			slog.Int("quantity", orderItem.Quantity),
		)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Product ID and quantity must be greater than zero",
		})
		return
	}

	if err := h.storage.AddOrderItem(r.Context(), orderItem); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{
				"error":   "Not found",
				"message": "Order or product not found",
			})
			return
		}

		log.Error("failed to add order item", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to add order item",
		})
		return
	}

	log.Info("Order item added successfully",
		slog.String("url", r.URL.String()),
		slog.Int("order_id", orderItem.OrderID),
		slog.Int("product_id", orderItem.ProductID),
	)

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]string{
		"status":  "success",
		"message": "Order item added successfully",
	})
}

func (h *HandlerItem) GetOrderItemsByOrderID(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.order.GetOrderItemsByOrderID"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	orderIDStr := chi.URLParam(r, "id")
	if orderIDStr == "" {
		log.Error("order ID is missing in the request")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID is missing",
		})
		return
	}

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil || orderID <= 0 {
		log.Error("invalid order ID", slog.String("id", orderIDStr))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID must be a positive integer",
		})
		return
	}

	orderItems, err := h.storage.GetOrderItemsByOrderID(r.Context(), orderID)
	if err != nil {
		log.Error("failed to get order items", slog.Any("error", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to get order items",
		})
		return
	}

	if orderItems == nil {
		orderItems = make([]models.OrderItem, 0)
	}

	log.Info("Order items retrieved successfully", slog.Int("order_id", orderID))
	render.JSON(w, r, orderItems)
}
