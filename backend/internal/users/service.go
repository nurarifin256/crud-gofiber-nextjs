package user

import "context"

type Service interface {
	FindUserByNik(ctx context.Context, nik string) (*User, error)
	SubmitUser(ctx context.Context, req SubmitUserRequest, picture string) (*User, error)
}

type ServiceImpl struct {
	repo Repository
}

func NewUserService(repo Repository) Service {
	return &ServiceImpl{
		repo: repo,
	}
}

func (s *ServiceImpl) FindUserByNik(ctx context.Context, nik string) (*User, error) {
	return s.repo.FindUserByNik(ctx, nik)
}

func (s *ServiceImpl) SubmitUser(ctx context.Context, req SubmitUserRequest, picture string) (*User, error) {
	return s.repo.SubmitUser(ctx, req, picture)
}
