package usecase

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/devonLoen/leave-request-service/internal/app/rest_api/model"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/model/dto"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/pkg/util"
	repository "github.com/devonLoen/leave-request-service/internal/app/rest_api/repository"
)

type Auth struct {
	userRepo *repository.User
}

func NewAuthUsecase(userRepo *repository.User) *Auth {
	return &Auth{userRepo: userRepo}
}

func (a *Auth) Login(req *dto.LoginRequest) (*dto.LoginResponse, *model.ErrorResponse) {
	loginResponse := &dto.LoginResponse{}

	user, err := a.userRepo.FindByEmailWithPassword(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &model.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Invalid email or password",
			}
		}
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, &model.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Invalid email or password",
		}
	}

	token, err := util.GenerateJWT(user.ID, string(user.Role))
	if err != nil {
		return nil, &model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
		}
	}

	return loginResponse.FromLogin(user, token), nil
}
