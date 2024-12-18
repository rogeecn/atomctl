package uuid

import (
	"git.ipao.vip/rogeecn/atom/container"
	"git.ipao.vip/rogeecn/atom/utils/opt"

	"github.com/gofrs/uuid"
)

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options:  []opt.Option{},
	}
}

type Generator struct {
	generator uuid.Generator
}

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	return container.Container.Provide(func() (*Generator, error) {
		return &Generator{
			generator: uuid.DefaultGenerator,
		}, nil
	}, o.DiOptions()...)
}

func (u *Generator) MustGenerate() string {
	uuid, _ := u.Generate()
	return uuid
}

func (u *Generator) Generate() (string, error) {
	uuid, err := u.generator.NewV4()
	if err != nil {
		return "", err
	}
	return uuid.String(), err
}
