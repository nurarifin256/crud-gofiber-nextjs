package user

import (
	"context"
	model "simple-crud/database/models"
	"simple-crud/helpers"
)

type Service interface {
	FindUserByNik(ctx context.Context, nik string) (*model.User, error)
	SubmitUser(ctx context.Context, req SubmitUserRequest, picture string) (*model.User, error)
}

type ServiceImpl struct {
	repo Repository
}

func NewUserService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) FindUserByNik(ctx context.Context, nik string) (*model.User, error) {
	return s.repo.FindUserByNik(ctx, nik)
}

func (s *ServiceImpl) SubmitUser(ctx context.Context, req SubmitUserRequest, picture string) (*model.User, error) {
	password := helpers.HashPassword(req.Password)
	return s.repo.SubmitUser(ctx, req, picture, password)
}
