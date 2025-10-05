package server

import (
	"fmt"
	"os/exec"

	"github.com/zobtube/zobtube/internal/controller"
)

func dependencyRegister(c controller.AbstractController, params *Parameters, dep string) {
	// external providers
	path, err := exec.LookPath(dep)
	if err != nil {
		params.Logger.Warn().Str("kind", "dependency").Err(err).Msg("unable to check dependency")
		c.RegisterError(fmt.Sprintf("Unable to ensure presence of %s with error: %s", dep, err.Error()))
	}
	params.Logger.Debug().
		Str("kind", "dependency").
		Str("dependency", dep).
		Msg(fmt.Sprintf("available at %s", path))
}
