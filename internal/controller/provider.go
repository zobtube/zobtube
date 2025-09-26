package controller

import (
	"errors"

	"github.com/zobtube/zobtube/internal/model"
	"github.com/zobtube/zobtube/internal/provider"
)

func (c *Controller) ProviderRegister(p provider.Provider) error {
	subLog := c.logger.With().Str("kind", "provider").Str("provider", p.SlugGet()).Logger()
	subLog.Debug().Msg("register provider")
	c.providers[p.SlugGet()] = p

	_provider := &model.Provider{
		ID: p.SlugGet(),
	}

	result := c.datastore.First(_provider)

	// check result
	if result.RowsAffected == 1 {
		// provider already registered, updating it if needed
		_provider.NiceName = p.NiceName()
		_provider.AbleToSearchActor = p.CapabilitySearchActor()
		_provider.AbleToScrapePicture = p.CapabilityScrapePicture()
	} else {
		// provider not registered, creating it now
		subLog.Info().Msg("first time seeing provider, creating its configuration")
		_provider = &model.Provider{
			ID:                  p.SlugGet(),
			Enabled:             true,
			NiceName:            p.NiceName(),
			AbleToSearchActor:   p.CapabilitySearchActor(),
			AbleToScrapePicture: p.CapabilityScrapePicture(),
		}
	}

	err := c.datastore.Save(&_provider).Error
	if err != nil {
		subLog.Error().Err(err).Msg("unable to create configuration")
		return err
	}

	return nil
}

func (c *Controller) ProviderGet(slug string) (p provider.Provider, err error) {
	p, ok := c.providers[slug]
	if ok {
		return p, nil
	}

	return p, errors.New("provider not found")
}
