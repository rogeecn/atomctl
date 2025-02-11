package event

import (
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"

	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
)

func ProvideKafka(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}

	return container.Container.Provide(func() (*PubSub, error) {
		logger := LogrusAdapter()

		publisher, err := kafka.NewPublisher(kafka.PublisherConfig{
			Brokers:   config.Brokers,
			Marshaler: kafka.DefaultMarshaler{},
		}, logger)
		if err != nil {
			return nil, err
		}

		subscriber, err := kafka.NewSubscriber(kafka.SubscriberConfig{
			Brokers:       config.Brokers,
			Unmarshaler:   kafka.DefaultMarshaler{},
			ConsumerGroup: config.ConsumerGroup,
		}, logger)
		if err != nil {
			return nil, err
		}

		router, err := message.NewRouter(message.RouterConfig{}, logger)
		if err != nil {
			return nil, err
		}

		return &PubSub{
			Publisher:  publisher,
			Subscriber: subscriber,
			Router:     router,
		}, nil
	}, o.DiOptions()...)
}
