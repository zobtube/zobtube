package controller

func (c *Controller) Shutdown() {
	c.logger.Warn().Str("kind", "system").Msg("shutdown requested")
	c.shutdownChannel <- 1
}

func (c *Controller) Restart() {
	c.logger.Warn().Str("kind", "system").Msg("restart requested")
	c.shutdownChannel <- 2
}
