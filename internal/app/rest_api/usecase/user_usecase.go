package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
	entities "github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
	models "github.com/devonLoen/leave-request-service/internal/app/rest_api/model"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/model/dto"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/pkg/util"
	repositories "github.com/devonLoen/leave-request-service/internal/app/rest_api/repository"
)

type User struct {
	userRepo *repositories.User
}

func NewUserService(userRepo *repositories.User) *User {
	return &User{userRepo: userRepo}
}

func (us *User) GetAllUsers(limit, offset int, sortBy, orderBy, search string, filter entities.UserFilter) (*dto.GetAllUsersResponse, *models.ErrorResponse) {
	response := &dto.GetAllUsersResponse{}

	allowedSorts := map[string]bool{
		"id":       true,
		"fullName": true,
		"email":    true,
		"role":     true,
	}

	safeSortBy := "id"
	if allowedSorts[sortBy] {
		safeSortBy = sortBy
	} else {
		return nil, &models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid sort parameter",
		}
	}

	safeOrderBy := "ASC"
	if strings.ToUpper(orderBy) == "DESC" {
		safeOrderBy = "DESC"
	} else if strings.ToUpper(orderBy) != "ASC" {
		return nil, &models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid order parameter",
		}
	}

	if filter.Role != "" {
		roleEnum := entity.UserRole(filter.Role)
		if !roleEnum.IsValid() {
			return nil, &models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid Role Filter parameter",
			}
		}
	}

	queriedUsers, err := us.userRepo.GetAllUsers(limit, offset, safeSortBy, safeOrderBy, search, filter)
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	response.MapUsersResponse(queriedUsers)

	return response, nil
}

func (us *User) GetUser(userID int) (*dto.UserResponse, *models.ErrorResponse) {
	response := &dto.UserResponse{}

	user, err := us.userRepo.FindById(userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &models.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "User Not Found",
			}
		}
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	response.MapUserResponse(user)

	return response, nil
}

func (us *User) CreateUser(createUserRequest *dto.CreateUserRequest) (*dto.CreateUserResponse, *models.ErrorResponse) {
	userResponse := &dto.CreateUserResponse{}

	errEmail := us.checkIfEmailExists(createUserRequest.Email)
	if errEmail != nil {
		return nil, errEmail
	}

	plainPassword, errPass := util.GenerateSecurePassword(12)
	if errPass != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	hashedPassword, errPass := util.HashPassword(plainPassword)
	if errPass != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}

	user := createUserRequest.ToUser()
	user.Password = hashedPassword

	err := us.userRepo.Create(user)
	if err != nil {
		return nil, &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
		}
	}
	fmt.Println(plainPassword)

	//send email implementation.....

	return userResponse.FromUser(user), nil
}

func (us *User) checkIfEmailExists(email string) *models.ErrorResponse {
	userWithEmail, err := us.userRepo.FindByEmail(email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &models.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
	}
	if userWithEmail != nil {
		return &models.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Email already in use",
		}
	}
	return nil
}
