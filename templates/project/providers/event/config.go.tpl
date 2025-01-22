package event

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/contracts"
	"{{.ModuleName}}/pkg/atom/opt"
)

const DefaultPrefix = "Events"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

type Config struct {
	ConsumerGroup string

	Brokers []string
}

type PubSub struct {
	Publisher  message.Publisher
	Subscriber message.Subscriber
	Router     *message.Router
}

func (ps *PubSub) Serve(ctx context.Context) error {
	if err := ps.Router.Run(ctx); err != nil {
		return err
	}
	return nil
}

func (ps *PubSub) Handle(
	handlerName string,
	consumerTopic string,
	publisherTopic string,
	handler message.HandlerFunc,
) {
	ps.Router.AddHandler(handlerName, consumerTopic, ps.Subscriber, publisherTopic, ps.Publisher, handler)
}

// publish
func (ps *PubSub) Publish(e contracts.EventPublisher) error {
	if e == nil {
		return nil
	}

	payload, err := e.Marshal()
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	return ps.Publisher.Publish(e.Topic(), msg)
}
