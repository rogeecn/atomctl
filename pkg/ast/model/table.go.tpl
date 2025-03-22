package models

import (
	"github.com/sirupsen/logrus"
)
// @provider
type {{.CamelTable}}Model struct {
	log *logrus.Entry `inject:"false"`
}

func (m *{{.CamelTable}}Model) Prepare() error {
	m.log = logrus.WithField("model", "{{.CamelTable}}Model")
	return nil
}