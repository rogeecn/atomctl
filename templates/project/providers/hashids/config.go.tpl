package hashids

import (
	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"
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
