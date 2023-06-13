//go:generate atomctl gen routes
//go:generate swag fmt
//go:generate swag init -ot json
package main

import (
	"os"
	"log"

	"{{.Package}}/database/migrations"
	"{{.Package}}/database/seeders"

	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom-addons/services/http"
	"github.com/spf13/cobra"
)

func main() {
	providers := http.Default()
	// providers := atom.DefaultGRPC()

	opts := []atom.Option{
		atom.Name("http"),
		atom.RunE(func(cmd *cobra.Command, args []string) error {
			return http.Serve()
			// return services.ServeGrpc()
		}),
		atom.Seeders(seeders.Seeders...),
		atom.Migrations(migrations.Migrations...),
	}

	if err := atom.Serve(providers, opts...); err != nil {
		log.Fatal(err)
	}
}
