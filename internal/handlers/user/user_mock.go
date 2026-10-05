package user

import (
	"context"
	"PETShoP/internal/models"
)

type UsersMock struct {
	GetAllUsersFunc    func(ctx context.Context) ([]models.User, error)
	GetUserByEmailFunc func(ctx context.Context, email string) (models.User, error)
	CreateUserFunc     func(ctx context.Context, u models.User) error
}


func (m *UsersMock) GetAllUsers(ctx context.Context) ([]models.User, error) {
	if m.GetAllUsersFunc != nil {
		return m.GetAllUsersFunc(ctx)
	}
	return []models.User{}, nil
}

func (m *UsersMock) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return models.User{}, nil
}

func (m *UsersMock) CreateUser(ctx context.Context, u models.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, u)
	}	
	return nil
}