package http

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func (s *Server) Start(bindAddress string) {
	s.Logger.Info().Str("kind", "system").Str("bind", bindAddress).Msg("http server binding")
	// #nosec G112
	s.Server = &http.Server{
		Addr:    bindAddress,
		Handler: s.Router.Handler(),
	}
	err := s.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Panic(err.Error())
	}
}

func (s *Server) WaitForStopSignal(c <-chan int) {
	mode := <-c
	s.Logger.Warn().Str("kind", "system").Int("signal", mode).Msg("http server signal received")

	// shutdown http server
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_ = s.Server.Shutdown(ctx)

	// mode: 1 - shutdown | 2 - restart
	switch mode {
	case 1:
		s.Logger.Warn().Str("kind", "system").Msg("server shutdown requested")
	case 2:
		s.Logger.Warn().Str("kind", "system").Msg("server restart requested")
		// #nosec G204
		cmd := exec.Command(os.Args[0])
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}
}
