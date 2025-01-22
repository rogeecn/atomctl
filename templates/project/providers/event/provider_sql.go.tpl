package event

import (
	sqlDB "database/sql"

	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/opt"

	"github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
)

func ProvideSQL(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}

	return container.Container.Provide(func(db *sqlDB.DB) (*PubSub, error) {
		logger := LogrusAdapter()

		publisher, err := sql.NewPublisher(db, sql.PublisherConfig{
			SchemaAdapter:        sql.DefaultPostgreSQLSchema{},
			AutoInitializeSchema: false,
		}, logger)
		if err != nil {
			return nil, err
		}

		subscriber, err := sql.NewSubscriber(db, sql.SubscriberConfig{
			SchemaAdapter: sql.DefaultPostgreSQLSchema{},
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
