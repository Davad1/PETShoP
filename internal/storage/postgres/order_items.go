package postgres

import (
	"PETShoP/internal/models"
	"PETShoP/internal/storage"
	"context"
	"fmt"
)

func (s *Storage) AddOrderItem(orderItem models.OrderItem) error {
	const fn = "storage.postgres.AddOrderItem"

	if _, err := s.GetOrderByID(orderItem.OrderID); err != nil {
		return err
	}

	if _, err := s.GetProductByID(context.Background(), orderItem.ProductID); err != nil {
		return err
	}

	_, err := s.db.Exec(context.Background(),
		`INSERT INTO order_items (order_id, product_id, quantity) VALUES ($1, $2, $3)`,
		orderItem.OrderID, orderItem.ProductID, orderItem.Quantity)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return s.UpdateOrderTotalPrice(orderItem.OrderID)
}

func (s *Storage) GetOrderItemsByOrderID(orderID int) ([]models.OrderItem, error) {
	const fn = "storage.postgres.GetOrderItemsByOrderID"

	rows, err := s.db.Query(context.Background(),
		`SELECT id, order_id, product_id, quantity FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var orderItems []models.OrderItem
	for rows.Next() {
		var orderItem models.OrderItem
		if err := rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		orderItems = append(orderItems, orderItem)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return orderItems, nil
}

func (s *Storage) UpdateOrderTotalPrice(orderID int) error {
	const fn = "storage.postgres.UpdateOrderTotalPrice"

	result, err := s.db.Exec(context.Background(),
		`UPDATE orders
		 SET total_price = (
		     SELECT COALESCE(SUM(p.price * oi.quantity), 0)
		     FROM order_items oi
		     JOIN products p ON p.id = oi.product_id
		     WHERE oi.order_id = $1
		 )
		 WHERE id = $1`, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if result.RowsAffected() == 0 {
		return storage.ErrNotFound
	}
	return nil
}
