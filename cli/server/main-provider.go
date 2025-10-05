package server

import (
	"fmt"

	"github.com/zobtube/zobtube/internal/controller"
	"github.com/zobtube/zobtube/internal/provider"
)

func providerRegister(c controller.AbstractController, params *Parameters, p provider.Provider) {
	err := c.ProviderRegister(p)
	if err != nil {
		params.Logger.Warn().Str("kind", "provider").Err(err).Msg("unable to register provider")
		c.RegisterError(fmt.Sprintf("Unable to register provider %s with error: %s", p.NiceName(), err.Error()))
	}
}
