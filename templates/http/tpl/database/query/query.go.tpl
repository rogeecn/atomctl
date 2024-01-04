package query

import (
	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/utils/opt"
	"gorm.io/gorm"
)

func Provide(...opt.Option) error {
	return container.Container.Provide(func(db *gorm.DB) *Query {
		SetDefault(db)
		return Q
	}, atom.GroupInitial)
}
func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options:  []opt.Option{},
	}
}
