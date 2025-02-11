package hashids

import (
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"

	"github.com/speps/go-hashids/v2"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}
	return container.Container.Provide(func() (*hashids.HashID, error) {
		data := hashids.NewData()
		data.MinLength = int(config.MinLength)
		if data.MinLength == 0 {
			data.MinLength = 10
		}

		data.Salt = config.Salt
		if data.Salt == "" {
			data.Salt = "default-salt-key"
		}

		data.Alphabet = config.Alphabet
		if config.Alphabet == "" {
			data.Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		}

		return hashids.NewWithData(data)
	}, o.DiOptions()...)
}
