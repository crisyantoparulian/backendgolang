package handler

import (
	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	Validator  *validator.Validate
	Repository repository.RepositoryInterface
	Config     *config.Config
}

type NewServerOptions struct {
	Validator  *validator.Validate
	Repository repository.RepositoryInterface
	Config     *config.Config
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Validator:  opts.Validator,
		Repository: opts.Repository,
		Config:     opts.Config,
	}
}
