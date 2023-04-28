package main

import (
	"os"

	"{{.Package}}/database/migrations"
	"{{.Package}}/database/seeders"

	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/providers/database/sqlite"
	"github.com/rogeecn/atom/providers/http"
	"github.com/rogeecn/atom/providers/http/gin"
	"github.com/rogeecn/atom/providers/log"
	"github.com/rogeecn/atom/services"
	"github.com/rogeecn/atom/utils/opt"
	"github.com/spf13/cobra"
)

var providers = container.Providers{
	{
		Provider: log.Provide,
		Options: []opt.Option{
			opt.Prefix(log.DefaultPrefix),
		},
	},
	{
		Provider: gin.Provide,
		Options: []opt.Option{
			opt.Prefix(http.DefaultPrefix),
		},
	},
	{
		Provider: sqlite.Provide,
		Options: []opt.Option{
			opt.Prefix(sqlite.DefaultPrefix),
		},
	}
}

func main() {
	// add module providers
	// providers = append(providers, system.Providers()...)

	opts := []atom.Option{
		atom.Name("http"),
		atom.RunE(func(cmd *cobra.Command, args []string) error {
			return services.ServeHttp()
		}),
		atom.Seeders(seeders.Seeders...),
		atom.Migrations(migrations.Migrations...),
	}

	if err := atom.Serve(providers, opts...); err != nil {
		log.Fatal(err)
	}
}
