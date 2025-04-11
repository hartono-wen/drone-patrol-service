package handler

import (
	"github.com/hartono-wen/drone-patrol-service/config"
	"github.com/hartono-wen/drone-patrol-service/repository"
)

type Server struct {
	Repository repository.RepositoryInterface
	Config     *config.Config
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Config     *config.Config
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{opts.Repository, opts.Config}
}
