package main

import (
	"os"

	"{{.Package}}/database/migrations"
	"{{.Package}}/database/seeders"

	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom/services"
	"github.com/spf13/cobra"
)

func main() {
	providers := atom.DefaultHTTP()
	// providers := atom.DefaultGRPC()

	opts := []atom.Option{
		atom.Name("http"),
		atom.RunE(func(cmd *cobra.Command, args []string) error {
			return services.ServeHttp()
			// return services.ServeGrpc()
		}),
		atom.Seeders(seeders.Seeders...),
		atom.Migrations(migrations.Migrations...),
	}

	if err := atom.Serve(providers, opts...); err != nil {
		os.Exit(1)
	}
}
