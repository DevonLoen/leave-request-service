package route

import (
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/handler"
	"github.com/gin-gonic/gin"
)

func RegisterPublicEndpoints(router *gin.Engine, userHandlers *handler.User, authHandlers *handler.Auth) {
	router.GET("/users", userHandlers.GetAllUsers)
	router.GET("/users/:id", userHandlers.GetUser)
	router.POST("/users", userHandlers.CreateUser)

	router.POST("/auth/login", authHandlers.Login)
}
