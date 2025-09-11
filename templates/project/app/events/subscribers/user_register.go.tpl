package subscribers

import (
	"encoding/json"

	"{{.ModuleName}}/app/events"
	"{{.ModuleName}}/app/events/publishers"
	"{{.ModuleName}}/providers/event"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sirupsen/logrus"
	"go.ipao.vip/atom/contracts"
)

var _ contracts.EventHandler = (*UserRegister)(nil)

// @provider(event)
type UserRegister struct {
	event.DefaultChannel
	event.DefaultPublishTo

	log *logrus.Entry `inject:"false"`
}

func (e *UserRegister) Prepare() error {
	e.log = logrus.WithField("module", "events.subscribers.user_register")
	return nil
}

// Topic implements contracts.EventHandler.
func (e *UserRegister) Topic() string {
	return events.TopicUserRegister
}

// Handler implements contracts.EventHandler.
func (e *UserRegister) Handler(msg *message.Message) ([]*message.Message, error) {
	var payload publishers.UserRegister
	err := json.Unmarshal(msg.Payload, &payload)
	if err != nil {
		return nil, err
	}
	e.log.Infof("received event %s", msg.Payload)

	return nil, nil
}
