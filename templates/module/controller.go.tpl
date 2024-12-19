package {{.ModuleName}}

import (
	log "github.com/sirupsen/logrus"
)

// @provider
type Controller struct {
	svc *Service
	log *log.Entry `inject:"false"`
}

func (c *Controller) Prepare() error {
	c.log = log.WithField("module", "{{.ModuleName}}.Controller")
	return nil
}
