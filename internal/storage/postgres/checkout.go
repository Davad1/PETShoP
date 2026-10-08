package postgres

import (
	"PETShoP/internal/models"
	"PETShoP/internal/storage"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) PlaceOrder(userEmail string, items []models.OrderItem) (int, error) {
	const fn = "storage.postgres.PlaceOrder"

	ctx := context.Background()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var userID int
	err = tx.QueryRow(
		ctx,
		`SELECT id FROM users WHERE email = $1`,
		userEmail,
	).Scan(&userID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, storage.ErrNotFound
		}

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	var orderID int
	err = tx.QueryRow(
		ctx,
		`INSERT INTO orders (user_id, total_price)
		 VALUES ($1, 0)
		 RETURNING id`,
		userID,
	).Scan(&orderID)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	var totalPrice float64

	for _, item := range items {
		var price float64

		err = tx.QueryRow(
			ctx,
			`UPDATE products
			 SET stock = stock - $1
			 WHERE id = $2 AND stock >= $1
			 RETURNING price`,
			item.Quantity,
			item.ProductID,
		).Scan(&price)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				var exists bool

				checkErr := tx.QueryRow(
					ctx,
					`SELECT EXISTS(
						SELECT 1
						FROM products
						WHERE id = $1
					)`,
					item.ProductID,
				).Scan(&exists)

				if checkErr != nil {
					return 0, fmt.Errorf("%s: %w", fn, checkErr)
				}

				if !exists {
					return 0, storage.ErrNotFound
				}

				return 0, storage.ErrInsufficientStock
			}

			return 0, fmt.Errorf("%s: %w", fn, err)
		}

		_, err = tx.Exec(
			ctx,
			`INSERT INTO order_items (order_id, product_id, quantity)
			 VALUES ($1, $2, $3)`,
			orderID,
			item.ProductID,
			item.Quantity,
		)

		if err != nil {
			return 0, fmt.Errorf("%s: %w", fn, err)
		}

		totalPrice += price * float64(item.Quantity)
	}
	_, err = tx.Exec(
		ctx,
		`UPDATE orders
		 SET total_price = $1
		 WHERE id = $2`,
		totalPrice,
		orderID,
	)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO transactions (order_id, amount, status)
		 VALUES ($1, $2, $3)`,
		orderID,
		totalPrice,
		"completed",
	)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return orderID, nil
}
