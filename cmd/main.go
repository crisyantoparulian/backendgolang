package main

import (
	"fmt"

	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.LoadConfig()

	e := echo.New()

	var server generated.ServerInterface = newServer(cfg)

	generated.RegisterHandlers(e, server)
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.App.Port)))
}

func newServer(cfg *config.Config) *handler.Server {
	validator := validator.New()

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: cfg.Database.PostgreDSN,
	})

	opts := handler.NewServerOptions{
		Repository: repo,
		Validator:  validator,
		Config:     cfg,
	}

	return handler.NewServer(opts)
}
