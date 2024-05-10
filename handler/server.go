package handler

import (
	"github.com/hartono-wen/sawitpro-technical-interview-software-architect/config"
	"github.com/hartono-wen/sawitpro-technical-interview-software-architect/repository"
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
