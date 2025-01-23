package opt

import (
	"github.com/spf13/viper"
	"go.uber.org/dig"
)

type Options struct {
	Config *viper.Viper
	Prefix string
	Name   string
	Group  string
}

type Option func(o *Options)

func New(opts ...Option) *Options {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *Options) UnmarshalConfig(config interface{}) error {
	return o.Config.UnmarshalKey(o.Prefix, &config)
}

func (o *Options) DiOptions() []dig.ProvideOption {
	options := []dig.ProvideOption{}
	if o.Name != "" {
		options = append(options, dig.Name(o.Name))
	}
	if o.Group != "" {
		options = append(options, dig.Group(o.Group))
	}
	return options
}

func Config(config *viper.Viper) Option {
	return func(o *Options) {
		o.Config = config
	}
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func Group(group string) Option {
	return func(o *Options) {
		o.Group = group
	}
}

func Prefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}
