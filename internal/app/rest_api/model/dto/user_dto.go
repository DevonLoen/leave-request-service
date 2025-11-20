package dto

import entities "github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"

type UserResponse struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type GetAllUsersResponse struct {
	Users []*UserResponse `json:"users"`
}

type CreateUserRequest struct {
	FullName string `json:"fullName" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=254"`
	Role     string `json:"role" binding:"required,oneof=superadmin admin employee"`
}

type CreateUserResponse struct {
	FullName string `json:"fullName" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email,max=254"`
	Role     string `json:"role" binding:"required,oneof=superadmin admin employee"`
	Message  string `json:"message" binding:"required"`
}

func (r *GetAllUsersResponse) MapUsersResponse(users []*entities.User) {
	for _, users := range users {
		user := &UserResponse{
			FullName: users.FullName,
			Email:    users.Email,
			Role:     string(users.Role),
		}
		r.Users = append(r.Users, user)
	}
}

func (r *UserResponse) MapUserResponse(user *entities.User) {
	r.FullName = user.FullName
	r.Email = user.Email
	r.Role = string(user.Role)
}

func (ur *CreateUserRequest) ToUser() *entities.User {
	return &entities.User{
		FullName: ur.FullName,
		Email:    ur.Email,
		Role:     entities.UserRole(ur.Role),
	}
}

func (ur *CreateUserResponse) FromUser(user *entities.User) *CreateUserResponse {
	return &CreateUserResponse{
		FullName: user.FullName,
		Email:    user.Email,
		Role:     string(user.Role),
		Message:  "User created successfully.",
	}
}
