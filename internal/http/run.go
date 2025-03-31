package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func (s *Server) Start(bindAddress string) {
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
	fmt.Println("http server signal received!", mode)

	// shutdown http server
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_ = s.Server.Shutdown(ctx)

	// mode: 1 - shutdown | 2 - restart
	switch mode {
	case 1:
		fmt.Println("server shutdown")
	case 2:
		log.Println("server restart")
		cmd := exec.Command(os.Args[0])
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}
}
