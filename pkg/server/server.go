package server

import (
	"github.com/statictask/newsletter/internal/log"
)

type Server struct {
	*log.Logger
}

func New() *Server {
	return &Server{log.NewLogger()}
}
