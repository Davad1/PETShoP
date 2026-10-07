package postgres

import (
	"PETShoP/internal/models"
	"PETShoP/internal/storage"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreateUser(ctx context.Context, u models.User) error {
	const fn = "storage.postgres.user.CreateUser"

	_, err := s.db.Exec(ctx,
		`INSERT INTO users (Name, email) VALUES ($1, $2)`,
		u.Name, u.Email)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const fn = "storage.postgres.user.GetUserByEmail"

	row := s.db.QueryRow(ctx, `SELECT id, name, email FROM users WHERE email = $1`, email)

	var u models.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, storage.ErrNotFound
		}
		return models.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	return u, nil
}

func (s *Storage) GetAllUsers(ctx context.Context) ([]models.User, error) {
	const fn = "storage.postgres.user.GetAllUsers"
	rows, err := s.db.Query(ctx, `SELECT id, name, email FROM users ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		users = append(users, u)
	}

	// Проверяем ошибки итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return users, nil
}
