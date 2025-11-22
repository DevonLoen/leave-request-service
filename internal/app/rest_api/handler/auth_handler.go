package handler

import (
	"errors"
	"net/http"

	dto "github.com/devonLoen/leave-request-service/internal/app/rest_api/model/dto"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/pkg/util"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Auth struct {
	AuthService *usecase.Auth
}

func NewAuthHandler(AuthService *usecase.Auth) *Auth {
	return &Auth{AuthService: AuthService}
}

func (h *Auth) Login(ctx *gin.Context) {
	var loginRequest dto.LoginRequest

	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
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

	loginResponse, loginError := h.AuthService.Login(&loginRequest)
	if loginError != nil {
		ctx.AbortWithStatusJSON(loginError.Code, loginError)
		return
	}

	ctx.JSON(http.StatusOK, loginResponse)
}
