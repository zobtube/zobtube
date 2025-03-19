package controller

import (
	"errors"

	"github.com/zobtube/zobtube/internal/provider"
)

func (c *Controller) ProviderRegister(p provider.Provider) {
	c.providers[p.SlugGet()] = p
}

func (c *Controller) ProviderGet(slug string) (p provider.Provider, err error) {
	p, ok := c.providers[slug]
	if ok {
		return p, nil
	}

	return p, errors.New("provider not found")
}
