package response

import "basic-trade/internal/entity"

type AdminResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewAdminResponse(admin entity.Admin) AdminResponse {
	return AdminResponse{
		Name:  admin.Name,
		Email: admin.Email,
	}
}

type LoginResponse struct {
	Token string        `json:"token"`
	Admin AdminResponse `json:"admin"`
}

func NewLoginResponse(token string, admin entity.Admin) LoginResponse {
	return LoginResponse{
		Token: token,
		Admin: NewAdminResponse(admin),
	}
}
