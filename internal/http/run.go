package http

import "log"

func (s *Server) Start(bindAddress string) {
	err := s.Server.Run(bindAddress)
	if err != nil {
		log.Panic(err.Error())
	}
}
