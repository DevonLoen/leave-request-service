package main

import (
	"time"

	serve "github.com/devonLoen/leave-request-service/api/server"
	routes "github.com/devonLoen/leave-request-service/api/server/router"
	"github.com/devonLoen/leave-request-service/config"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/database"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/handler"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/pkg/util"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/repository"
	"github.com/devonLoen/leave-request-service/internal/app/rest_api/usecase"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	config := config.NewConfig()

	client, err := database.NewSQLClient(database.Config{
		DBDriver:          config.Database.DatabaseDriver,
		DBSource:          config.Database.DatabaseSource,
		MaxOpenConns:      25,
		MaxIdleConns:      25,
		ConnMaxIdleTime:   15 * time.Minute,
		ConnectionTimeout: 5 * time.Second,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database client")
		return
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Error().Msgf("Failed to close database client: %v", err)
		}
	}()

	util.SetupJWT(config.JWT.Secret)

	userRepo := repository.NewUserRepository(client.DB)

	userUsecase := usecase.NewUserUsecase(userRepo)

	userHandler := handler.NewUserHandler(userUsecase)

	authUsecase := usecase.NewAuthUsecase(userRepo)

	authHandler := handler.NewAuthHandler(authUsecase)

	leaveRequestRepo := repository.NewLeaveRequestRepository(client.DB)

	leaveRequestUsecase := usecase.NewLeaveRequestUsecase(leaveRequestRepo)

	leaveRequestHandler := handler.NewLeaveRequestHandler(leaveRequestUsecase)

	cors := config.CorsNew()

	router := gin.Default()
	router.Use(cors)

	routes.RegisterPublicEndpoints(router, userHandler, authHandler, leaveRequestHandler)

	server := serve.NewServer(log.Logger, router, config)
	server.Serve()
}
