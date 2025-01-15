package subscribers

import (
	"encoding/json"

	"{{.ModuleName}}/app/events"
	"{{.ModuleName}}/app/events/publishers"

	"git.ipao.vip/rogeecn/atom/contracts"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sirupsen/logrus"
)

var _ contracts.EventHandler = (*{{.Name}}Subscriber)(nil)

// @provider(event)
type {{.Name}}Subscriber struct {
	log *logrus.Entry `inject:"false"`
}

func (e *{{.Name}}Subscriber) Prepare() error {
	e.log = logrus.WithField("module", "events.subscribers.{{.Name}}Subscriber")
	return nil
}

// PublishToTopic implements contracts.EventHandler.
func (e *{{.Name}}Subscriber) PublishToTopic() string {
	return events.TopicProcessed
}

// Topic implements contracts.EventHandler.
func (e *{{.Name}}Subscriber) Topic() string {
	return events.Topic{{.Name}}
}

// Handler implements contracts.EventHandler.
func (e *{{.Name}}Subscriber) Handler(msg *message.Message) ([]*message.Message, error) {
	var payload publishers.{{.Name}}Event
	err := json.Unmarshal(msg.Payload, &payload)
	if err != nil {
		return nil, err
	}
	e.log.Infof("received event %s", msg.Payload)

	// TODO: handle post deletion

	return nil, nil
}
