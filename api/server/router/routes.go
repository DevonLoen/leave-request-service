package route

import (
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/handler"
	"github.com/gin-gonic/gin"
)

func RegisterPublicEndpoints(router *gin.Engine, userHandlers *handler.User) {
	router.GET("/users", userHandlers.GetAllUsers)
	router.GET("/users/:id", userHandlers.GetUser)
	router.POST("/users", userHandlers.CreateUser)
}
