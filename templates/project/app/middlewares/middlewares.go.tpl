package middlewares

import (
	log "github.com/sirupsen/logrus"
)

// @provider
type Middlewares struct {
	log *log.Entry `inject:"false"`
}

func (f *Middlewares) Prepare() error {
	f.log = log.WithField("module", "middleware")
	return nil
}
