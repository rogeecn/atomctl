package events

import (
	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}

	return container.Container.Provide(func() (*PubSub, error) {
		logger := LogrusAdapter()

		client := gochannel.NewGoChannel(gochannel.Config{}, logger)
		router, err := message.NewRouter(message.RouterConfig{}, logger)
		if err != nil {
			return nil, err
		}

		return &PubSub{
			Publisher:  client,
			Subscriber: client,
			Router:     router,
		}, nil
	}, o.DiOptions()...)
}
