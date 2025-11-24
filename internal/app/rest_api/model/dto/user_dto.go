package dto

import entity "github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"

type UserResponse struct {
	ID       int    `json:"id"`
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

func (r *GetAllUsersResponse) MapUsersResponse(users []*entity.User) {
	for _, users := range users {
		user := &UserResponse{
			ID:       users.ID,
			FullName: users.FullName,
			Email:    users.Email,
			Role:     string(users.Role),
		}
		r.Users = append(r.Users, user)
	}
}

func (r *UserResponse) MapUserResponse(user *entity.User) {
	r.ID = user.ID
	r.FullName = user.FullName
	r.Email = user.Email
	r.Role = string(user.Role)
}

func (ur *CreateUserRequest) ToUser() *entity.User {
	return &entity.User{
		FullName: ur.FullName,
		Email:    ur.Email,
		Role:     entity.UserRole(ur.Role),
	}
}

func (ur *CreateUserResponse) FromUser(user *entity.User) *CreateUserResponse {
	return &CreateUserResponse{
		FullName: user.FullName,
		Email:    user.Email,
		Role:     string(user.Role),
		Message:  "User created successfully.",
	}
}
