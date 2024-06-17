package service

import (
	"basic-trade/common"
	"basic-trade/internal/entity"
	"basic-trade/internal/repository"
	"basic-trade/pkg/password"
	"basic-trade/pkg/token"
	"context"
	"fmt"

	sqlc "basic-trade/internal/repository/sqlc"
)

type IAuthService struct {
	adminRepo repository.AdminRepository
	jwtImpl   token.JWT
}

type AuthService interface {
	Register(ctx context.Context, admin entity.Admin) (entity.Admin, error)
	Login(ctx context.Context, admin entity.Admin) (entity.Admin, string, error)
}

func NewAuthService(adminRepo repository.AdminRepository, jwtImpl token.JWT) AuthService {
	return &IAuthService{adminRepo: adminRepo, jwtImpl: jwtImpl}
}

func (s *IAuthService) Register(ctx context.Context, admin entity.Admin) (entity.Admin, error) {
	hashedPassword, err := password.HashPassword(admin.Password)
	if err != nil {
		return entity.Admin{}, err
	}

	arg := sqlc.CreateAdminParams{
		Name:     admin.Name,
		Email:    admin.Email,
		Password: hashedPassword,
	}

	result, err := s.adminRepo.CreateAdmin(ctx, arg)
	if err != nil {
		return entity.Admin{}, err
	}

	return entity.CreateAdminToViewModel(result), err
}

func (s *IAuthService) Login(ctx context.Context, admin entity.Admin) (entity.Admin, string, error) {
	result, err := s.adminRepo.GetAdmin(ctx, admin.Email)
	if err != nil {
		return entity.Admin{}, "", err
	}

	err = password.CheckPassword(admin.Password, result.Password)
	if err != nil {
		err := fmt.Errorf("%d", common.ErrCredentiials)
		return entity.Admin{}, "", err
	}

	accessToken, err := s.jwtImpl.CreateAccessToken(result.Uuid)
	if err != nil {
		return entity.Admin{}, "", err
	}

	return entity.GetAdminToViewModel(result), accessToken.SignedToken, nil
}
