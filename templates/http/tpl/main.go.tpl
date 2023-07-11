//go:generate atomctl gen routes
//go:generate swag fmt
//go:generate swag init -ot json
package main

import (
	"log"

	"{{.Package}}/database/migrations"

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
	}

	if err := atom.Serve(providers, opts...); err != nil {
		log.Fatal(err)
	}
}
