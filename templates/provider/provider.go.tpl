package {{.Name}}

import (
	"{{.ModuleName}}/pkg/atom/container"
	"{{.ModuleName}}/pkg/atom/opt"
)

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	var config Config
	if err := o.UnmarshalConfig(&config); err != nil {
		return err
	}
	return container.Container.Provide(func() (*Config, error) {
		return &Config{}, nil
	}, o.DiOptions()...)
}
