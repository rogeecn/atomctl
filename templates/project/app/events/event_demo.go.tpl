package events

import (
	"encoding/json"
	"fmt"
	"time"

	"git.ipao.vip/rogeecn/atom/contracts"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/sirupsen/logrus"
)

var _ contracts.EventHandler = (*UserRegister)(nil)

type Event struct {
	ID int `json:"id"`
}

type ProcessedEvent struct {
	ProcessedID int       `json:"processed_id"`
	Time        time.Time `json:"time"`
}

// @provider(event)
type UserRegister struct {
	log *logrus.Entry `inject:"false"`
}

func (u *UserRegister) Prepare() error {
	return nil
}

// Handler implements contracts.EventHandler.
func (u *UserRegister) Handler(msg *message.Message) ([]*message.Message, error) {
	consumedPayload := Event{}
	err := json.Unmarshal(msg.Payload, &consumedPayload)
	if err != nil {
		// When a handler returns an error, the default behavior is to send a Nack (negative-acknowledgement).
		// The message will be processed again.
		//
		// You can change the default behaviour by using middlewares, like Retry or PoisonQueue.
		// You can also implement your own middleware.
		return nil, err
	}

	fmt.Printf("received event %+v\n", consumedPayload)

	newPayload, err := json.Marshal(ProcessedEvent{
		ProcessedID: consumedPayload.ID,
		Time:        time.Now(),
	})
	if err != nil {
		return nil, err
	}

	newMessage := message.NewMessage(watermill.NewUUID(), newPayload)

	return nil, nil
	return []*message.Message{newMessage}, nil
}

// PublishToTopic implements contracts.EventHandler.
func (u *UserRegister) PublishToTopic() string {
	return "event:processed"
}

// Topic implements contracts.EventHandler.
func (u *UserRegister) Topic() string {
	return "event:user-register"
}
