package dto

import entities "github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=254"`
	Password string `json:"password" binding:"required,min=3,max=50"`
}

type LoginResponse struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Token    string `json:"token"`
	Message  string `json:"message"`
}

func (lr *LoginResponse) FromLogin(user *entities.User, token string) *LoginResponse {

	return &LoginResponse{
		ID:       user.ID,
		FullName: user.FullName,
		Email:    user.Email,
		Role:     string(user.Role),
		Token:    token,
		Message:  "Login successful.",
	}
}
