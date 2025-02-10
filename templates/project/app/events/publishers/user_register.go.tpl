package publishers

import (
	"encoding/json"

	"{{.ModuleName}}/app/events"
	"{{.ModuleName}}/pkg/atom/contracts"
)

var _ contracts.EventPublisher = (*UserRegister)(nil)

type UserRegister struct {
	ID int64 `json:"id"`
}

func (e *UserRegister) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

func (e *UserRegister) Topic() string {
	return events.TopicUserRegister
}
