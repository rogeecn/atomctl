package {{.ModuleName}}

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
	log "github.com/sirupsen/logrus"
)

// @provider:except
type Service struct {
	db  *sql.DB
	log *log.Entry `inject:"false"`
}

func (svc *Service) Prepare() error {
	svc.log = log.WithField("module", "{{.ModuleName}}.service")
	_ = Int(1)
	return nil
}
