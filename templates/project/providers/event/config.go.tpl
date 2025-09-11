package event

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/contracts"
	"go.ipao.vip/atom/opt"
)

const DefaultPrefix = "Events"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: ProvideChannel,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

type Config struct {
	Sql   *ConfigSql
	Kafka *ConfigKafka
	Redis *ConfigRedis
}

type ConfigSql struct {
	ConsumerGroup string
}

type ConfigRedis struct {
	ConsumerGroup string
	Streams       []string
}

type ConfigKafka struct {
	ConsumerGroup string
	Brokers       []string
}

type PubSub struct {
	Router *message.Router

	publishers  map[contracts.Channel]message.Publisher
	subscribers map[contracts.Channel]message.Subscriber
}

func (ps *PubSub) Serve(ctx context.Context) error {
	if err := ps.Router.Run(ctx); err != nil {
		return err
	}
	return nil
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
	return ps.getPublisher(e.Channel()).Publish(e.Topic(), msg)
}

// getPublisher returns the publisher for the specified channel.
func (ps *PubSub) getPublisher(channel contracts.Channel) message.Publisher {
	if pub, ok := ps.publishers[channel]; ok {
		return pub
	}
	return ps.publishers[Go]
}

func (ps *PubSub) getSubscriber(channel contracts.Channel) message.Subscriber {
	if sub, ok := ps.subscribers[channel]; ok {
		return sub
	}
	return ps.subscribers[Go]
}

func (ps *PubSub) Handle(handlerName string, sub contracts.EventHandler) {
	publishToCh, publishToTopic := sub.PublishTo()

	ps.Router.AddHandler(
		handlerName,
		sub.Topic(),
		ps.getSubscriber(sub.Channel()),
		publishToTopic,
		ps.getPublisher(publishToCh),
		sub.Handler,
	)
}
