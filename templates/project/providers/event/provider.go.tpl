package event

import (
	sqlDB "database/sql"

	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/contracts"
	"go.ipao.vip/atom/opt"

	"github.com/ThreeDotsLabs/watermill-kafka/v3/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/redis/go-redis/v9"
)

func ProvideChannel(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}

	return container.Container.Provide(func() (*PubSub, error) {
		logger := LogrusAdapter()

		publishers := make(map[contracts.Channel]message.Publisher)
		subscribers := make(map[contracts.Channel]message.Subscriber)

		// gochannel
		client := gochannel.NewGoChannel(gochannel.Config{}, logger)
		publishers[Go] = client
		subscribers[Go] = client

		// kafka
		if config.Kafka != nil {
			kafkaPublisher, err := kafka.NewPublisher(kafka.PublisherConfig{
				Brokers:   config.Kafka.Brokers,
				Marshaler: kafka.DefaultMarshaler{},
			}, logger)
			if err != nil {
				return nil, err
			}
			publishers[Kafka] = kafkaPublisher

			kafkaSubscriber, err := kafka.NewSubscriber(kafka.SubscriberConfig{
				Brokers:       config.Kafka.Brokers,
				Unmarshaler:   kafka.DefaultMarshaler{},
				ConsumerGroup: config.Kafka.ConsumerGroup,
			}, logger)
			if err != nil {
				return nil, err
			}
			subscribers[Kafka] = kafkaSubscriber
		}

		// redis
		if config.Redis != nil {
			var rdb redis.UniversalClient
			redisSubscriber, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
				Client:        rdb,
				Unmarshaller:  redisstream.DefaultMarshallerUnmarshaller{},
				ConsumerGroup: config.Redis.ConsumerGroup,
			}, logger)
			if err != nil {
				return nil, err
			}
			subscribers[Redis] = redisSubscriber

			redisPublisher, err := redisstream.NewPublisher(redisstream.PublisherConfig{
				Client:     rdb,
				Marshaller: redisstream.DefaultMarshallerUnmarshaller{},
			}, logger)
			if err != nil {
				return nil, err
			}
			publishers[Redis] = redisPublisher
		}

		if config.Sql == nil {
			var db *sqlDB.DB
			sqlPublisher, err := sql.NewPublisher(db, sql.PublisherConfig{
				SchemaAdapter:        sql.DefaultPostgreSQLSchema{},
				AutoInitializeSchema: false,
			}, logger)
			if err != nil {
				return nil, err
			}
			publishers[Sql] = sqlPublisher

			sqlSubscriber, err := sql.NewSubscriber(db, sql.SubscriberConfig{
				SchemaAdapter: sql.DefaultPostgreSQLSchema{},
				ConsumerGroup: config.Sql.ConsumerGroup,
			}, logger)
			if err != nil {
				return nil, err
			}
			subscribers[Sql] = sqlSubscriber
		}

		router, err := message.NewRouter(message.RouterConfig{}, logger)
		if err != nil {
			return nil, err
		}

		return &PubSub{Router: router, publishers: publishers, subscribers: subscribers}, nil
	}, o.DiOptions()...)
}
