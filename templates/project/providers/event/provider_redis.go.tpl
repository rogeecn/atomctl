package event

import (
	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func ProvideRedis(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}

	return container.Container.Provide(func(rdb redis.UniversalClient) (*PubSub, error) {
		logger := LogrusAdapter()

		subscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
			Client:        rdb,
			Unmarshaller:  redisstream.DefaultMarshallerUnmarshaller{},
			ConsumerGroup: config.ConsumerGroup,
		}, logger)
		if err != nil {
			return nil, err
		}

		publisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
			Client:     rdb,
			Marshaller: redisstream.DefaultMarshallerUnmarshaller{},
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
