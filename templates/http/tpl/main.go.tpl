//go:generate atomctl gen routes
//go:generate swag fmt
//go:generate swag init -ot json
package main

import (
	"log"

	"github.com/rogeecn/atom"
	serviceHttp "github.com/atom-providers/service-http"
	// serviceGrpc "github.com/atom-providers/service-grpc"
	"github.com/spf13/cobra"
)

func main() {
	providers := serviceHttp.Default()
	// providers := serviceGrpc.Default()

	opts := []atom.Option{
		atom.Name("{{ .Name }}"),
		atom.RunE(func(cmd *cobra.Command, args []string) error {
			return serviceHttp.Serve()
			// return serviceGrpc.Serve()
		}),
	}

	if err := atom.Serve(providers, opts...); err != nil {
		log.Fatal(err)
	}
}
