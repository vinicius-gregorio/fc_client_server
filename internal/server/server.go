package server

import "github.com/vinicius-gregorio/fc_client_server/internal/db"

type Server struct {
	DB db.DB
}

func (s *Server) NewServer(database db.DB) *Server {
	return &Server{DB: database}
}
