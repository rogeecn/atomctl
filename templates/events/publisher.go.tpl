package publishers

import (
	"encoding/json"

	"{{.ModuleName}}/app/events"

	"{{.ModuleName}}/pkg/atom/contracts"
)

var _ contracts.EventPublisher = (*{{.Name}}Event)(nil)

type {{.Name}}Event struct {
	ID int64 `json:"id"`
}

func (e *{{.Name}}Event) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *{{.Name}}Event) Topic() string {
	return events.Topic{{.Name}}
}
