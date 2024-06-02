package service

import (
	"basic-trade/internal/entity"
	"basic-trade/internal/repository"
	"basic-trade/pkg/config"
	"basic-trade/pkg/password"
	"basic-trade/pkg/token"

	sqlc "basic-trade/internal/repository/sqlc"
)

type AuthService struct {
	adminRepo  repository.IAdminRepository
	tokenMaker token.Maker
	config     config.Token
}

type IAuthService interface {
	Register(admin entity.Admin) (entity.Admin, error)
	Login(admin entity.Admin) (entity.Admin, string, error)
}

func NewAuthService(adminRepo repository.IAdminRepository, tokenMaker token.Maker, cfg config.Token) *AuthService {
	return &AuthService{adminRepo: adminRepo, tokenMaker: tokenMaker, config: cfg}
}

func (s *AuthService) Register(admin entity.Admin) (entity.Admin, error) {
	hashedPassword, err := password.HashPassword(admin.Password)
	if err != nil {
		return entity.Admin{}, err
	}

	arg := sqlc.CreateAdminParams{
		Name:     admin.Name,
		Email:    admin.Email,
		Password: hashedPassword,
	}

	result, err := s.adminRepo.CreateAdmin(arg)
	if err != nil {
		return entity.Admin{}, err
	}

	return entity.CreateAdminToViewModel(&result), err
}

func (s *AuthService) Login(admin entity.Admin) (entity.Admin, string, error) {
	result, err := s.adminRepo.GetAdmin(admin.Email)
	if err != nil {
		return entity.Admin{}, "", err
	}

	err = password.CheckPassword(admin.Password, result.Password)
	if err != nil {
		return entity.Admin{}, "", err
	}

	accessToken, _, err := s.tokenMaker.CreateToken(
		result.Uuid,
		s.config.Duration,
	)
	if err != nil {
		return entity.Admin{}, "", err
	}

	return entity.GetAdminToViewModel(&result), accessToken, nil
}
