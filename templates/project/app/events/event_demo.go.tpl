package events

import (
	"encoding/json"

	"git.ipao.vip/rogeecn/atom/contracts"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sirupsen/logrus"
)

var (
	_ contracts.EventHandler   = (*UserRegister)(nil)
	_ contracts.EventPublisher = (*UserRegister)(nil)
)

// @provider(event)
type UserRegister struct {
	log *logrus.Entry `inject:"false" json:"-"`
	ID  int64         `json:"id"`
}

func (e *UserRegister) Prepare() error {
	return nil
}

// Marshal implements contracts.EventPublisher.
func (e *UserRegister) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// PublishToTopic implements contracts.EventHandler.
func (e *UserRegister) PublishToTopic() string {
	return TopicProcessed.String()
}

// Topic implements contracts.EventHandler.
func (e *UserRegister) Topic() string {
	return TopicUserRegister.String()
}

// Handler implements contracts.EventHandler.
func (e *UserRegister) Handler(msg *message.Message) ([]*message.Message, error) {
	var payload UserRegister
	err := json.Unmarshal(msg.Payload, &payload)
	if err != nil {
		return nil, err
	}

	e.log.Infof("received event %+v\n", payload)

	return nil, nil
}
