package controller

import (
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
)

// --- setup helper ---

func setupShutdownController() (*Controller, chan int) {
	logger := zerolog.Nop()

	// Full-duplex channel for testing (we can both send and receive)
	shutdown := make(chan int, 1)

	ctrl := &Controller{
		logger:          &logger,
		shutdownChannel: shutdown, // assign send-only view
	}

	return ctrl, shutdown
}

// --- tests ---

func TestController_Shutdown_BroadcastsAndWaits(t *testing.T) {
	ctrl, shutdown := setupShutdownController()

	var wg sync.WaitGroup
	wg.Add(1)

	// Listen for shutdown signal using the real (bidirectional) channel
	go func() {
		defer wg.Done()
		select {
		case <-shutdown:
			// received signal
		case <-time.After(2 * time.Second):
			t.Error("timeout waiting for shutdown signal")
		}
	}()

	ctrl.Shutdown()
}

func TestController_Shutdown_SignalReceived(t *testing.T) {
	ctrl, shutdown := setupShutdownController()

	done := make(chan struct{}, 1)
	go func() {
		select {
		case <-shutdown:
			done <- struct{}{}
		case <-time.After(2 * time.Second):
			t.Error("timeout waiting for shutdown signal")
		}
	}()

	ctrl.Shutdown()

	select {
	case <-done:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("did not receive shutdown broadcast")
	}
}
