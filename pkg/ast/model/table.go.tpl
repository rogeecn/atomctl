package model

import (
	log "github.com/sirupsen/logrus"
)

func (m *{{.PascalTable}}) log() *log.Entry {
	return log.WithField("model", "{{.PascalTable}}Model")
}