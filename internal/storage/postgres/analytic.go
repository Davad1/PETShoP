package postgres

import (
	"context"
	"fmt"

	"PETShoP/internal/models"
)

func (s *Storage) GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error) {
	const fn = "storage.postgres.GetUserOrderHistory"

	rows, err := s.db.Query(
		ctx,
		`
		SELECT
			o.id,
			COALESCE(oi.product_id, 0),
			COALESCE(p.name, ''),
			COALESCE(oi.quantity, 0),
			o.total_price,
			COALESCE(t.status, ''),
			o.created_at
		FROM users u
		JOIN orders o ON o.user_id = u.id
		LEFT JOIN order_items oi ON oi.order_id = o.id
		LEFT JOIN products p ON p.id = oi.product_id
		LEFT JOIN transactions t ON t.order_id = o.id
		WHERE u.email = $1
		ORDER BY o.created_at DESC, o.id DESC
		`,
		email,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var history []models.OrderDetail

	for rows.Next() {
		var detail models.OrderDetail

		if err := rows.Scan(
			&detail.OrderID,
			&detail.ProductID,
			&detail.ProductName,
			&detail.Quantity,
			&detail.TotalPrice,
			&detail.TransactionStatus,
			&detail.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		history = append(history, detail)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return history, nil
}

func (s *Storage) GetPopularProducts(ctx context.Context) ([]models.PopularProduct, error) {
	const fn = "storage.postgres.GetPopularProducts"

	rows, err := s.db.Query(
		ctx,
		`
		SELECT
			p.id,
			p.name,
			SUM(oi.quantity) AS total_sold
		FROM products p
		JOIN order_items oi ON oi.product_id = p.id
		JOIN transactions t ON t.order_id = oi.order_id
		WHERE t.status = 'completed'
		GROUP BY p.id, p.name
		ORDER BY total_sold DESC, p.id
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var products []models.PopularProduct

	for rows.Next() {
		var product models.PopularProduct

		if err := rows.Scan(
			&product.ProductID,
			&product.Name,
			&product.TotalSold,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return products, nil
}
