package hashids

import (
	"go.ipao.vip/atom/container"
	"go.ipao.vip/atom/opt"
)

const DefaultPrefix = "HashIDs"

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options: []opt.Option{
			opt.Prefix(DefaultPrefix),
		},
	}
}

type Config struct {
	Alphabet  string
	Salt      string
	MinLength uint
}
