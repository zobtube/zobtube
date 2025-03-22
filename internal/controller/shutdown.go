package controller

func (c *Controller) Shutdown() {
	c.shutdownChannel <- 1
}

func (c *Controller) Restart() {
	c.shutdownChannel <- 2
}
