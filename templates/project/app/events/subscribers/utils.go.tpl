package subscribers

import (
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

func toMessage(event any) (*message.Message, error) {
	b, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	return message.NewMessage(watermill.NewUUID(), b), nil
}

func toMessageList(event any) ([]*message.Message, error) {
	m, err := toMessage(event)
	if err != nil {
		return nil, err
	}
	return []*message.Message{m}, nil
}
