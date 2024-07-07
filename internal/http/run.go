package http

import ()

func (s *Server) Start(bindAddress string) {
	s.Server.Run(bindAddress)
}
