package config

import (
	"fmt"
	"net/http"
	"os"

	constants "github.com/devonLoen/leave-request-service/internal/app/rest_api/constant"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Config struct {
	Server     serverConfig
	Database   databaseConfig
	SuperAdmin superAdminConfig
	JWT        jwtConfig
}

type jwtConfig struct {
	Secret string
}

type superAdminConfig struct {
	Email    string
	Password string
}

type serverConfig struct {
	Address string
}

type databaseConfig struct {
	DatabaseDriver string
	DatabaseSource string
}

func NewConfig() *Config {
	err := godotenv.Load("dev.env")

	if err != nil {
		panic("Error loading .env file")
	}

	c := &Config{
		Server: serverConfig{
			Address: GetEnvOrPanic(constants.EnvKeys.ServerAddress),
		},
		Database: databaseConfig{
			DatabaseDriver: GetEnvOrPanic(constants.EnvKeys.DBDriver),
			DatabaseSource: GetEnvOrPanic(constants.EnvKeys.DBSource),
		},
		SuperAdmin: superAdminConfig{
			Email:    GetEnvOrPanic(constants.EnvKeys.SuperAdminEmail),
			Password: GetEnvOrPanic(constants.EnvKeys.SuperAdminPassword),
		},
		JWT: jwtConfig{
			Secret: GetEnvOrPanic(constants.EnvKeys.JwtSecret),
		},
	}

	return c
}

func GetEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("environment variable %s not set", key))
	}

	return value
}

func (conf *Config) CorsNew() gin.HandlerFunc {
	allowedOrigin := GetEnvOrPanic(constants.EnvKeys.CorsAllowedOrigin)

	return cors.New(cors.Config{
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{constants.Headers.Origin},
		ExposeHeaders:    []string{constants.Headers.ContentLength},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == allowedOrigin
		},
		MaxAge: constants.MaxAge,
	})
}
