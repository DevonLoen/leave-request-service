package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/devonLoen/leave-request-service/internal/app/rest_api/entity"
	dto "github.com/devonLoen/leave-request-service/internal/app/rest_api/model/dto"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/pkg/util"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type User struct {
	userUsecase *usecase.User
}

func NewUserHandler(userUsecase *usecase.User) *User {
	return &User{userUsecase: userUsecase}
}

func (h *User) GetAllUsers(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")
	sortByStr := ctx.DefaultQuery("sortBy", "id")
	orderByStr := ctx.DefaultQuery("orderBy", "asc")

	filter := entity.UserFilter{
		Role: ctx.Query("role"),
	}

	search := ctx.Query("search")

	page, errConv := strconv.Atoi(pageStr)
	if errConv != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Page not valid "})

		return
	}

	limit, errConv := strconv.Atoi(limitStr)
	if errConv != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "limit not valid "})

		return
	}

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	allUsers, err := h.userUsecase.GetAllUsers(limit, offset, sortByStr, orderByStr, search, filter)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Code, err)

		return
	}

	ctx.JSON(http.StatusOK, allUsers)
}

func (h *User) GetUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ID not valid"})

		return
	}

	user, userErr := h.userUsecase.GetUser(userID)
	if userErr != nil {
		ctx.AbortWithStatusJSON(userErr.Code, userErr)

		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *User) CreateUser(ctx *gin.Context) {
	var createUserRequest dto.CreateUserRequest

	if err := util.StrictBindJSON(ctx, &createUserRequest); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.New().Struct(createUserRequest); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make(map[string]string)
			for _, fe := range ve {
				out[fe.Field()] = util.MsgForTag(fe)
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createUserResponse, signupError := h.userUsecase.CreateUser(&createUserRequest)
	if signupError != nil {
		ctx.AbortWithStatusJSON(signupError.Code, signupError)

		return
	}

	ctx.JSON(http.StatusCreated, createUserResponse)
}
