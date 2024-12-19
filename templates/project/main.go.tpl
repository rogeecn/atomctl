package main

import (
	"{{.ModuleName}}/pkg/service/http"

	"git.ipao.vip/rogeecn/atom"
	log "github.com/sirupsen/logrus"
)

func main() {
	opts := []atom.Option{
		atom.Name("{{ .ProjectName }}"),
		http.Command(),
	}

	if err := atom.Serve(opts...); err != nil {
		log.Fatal(err)
	}
}
