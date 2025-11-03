package server

import (
	"github.com/zobtube/zobtube/internal/controller"
	"github.com/zobtube/zobtube/internal/http"
)

func startFailsafeWebServer(httpServer *http.Server, err error, c controller.AbstractController) {
	httpServer.Logger.Warn().
		Str("mode", "failsafe").
		Str("reason", "error during boot").
		Err(err).
		Msg("start zobtube in failsafe mode")
	httpServer.ControllerSetupFailsafeError(c, err)

	// handle shutdown
	go httpServer.WaitForStopSignal(shutdownChannel)

	_err := httpServer.Start("0.0.0.0:8069")
	if _err != nil {
		httpServer.Logger.Error().Err(err).Msg("unable to start failsafe server")
		return
	}

	// Wait for all HTTP fetches to complete.
	wg.Wait()

	httpServer.Logger.Warn().Msg("zobtube exiting failsafe webserver")
}
