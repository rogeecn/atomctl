package main

import (
	"{{.ModuleName}}/pkg/service/http"
	"{{.ModuleName}}/pkg/service/migrate"
	"{{.ModuleName}}/pkg/service/model"

	"git.ipao.vip/rogeecn/atom"
	log "github.com/sirupsen/logrus"
)

func main() {
	opts := []atom.Option{
		atom.Name("{{ .ProjectName }}"),
		http.Command(),
		migrate.Command(),
		model.Command(),
	}

	if err := atom.Serve(opts...); err != nil {
		log.Fatal(err)
	}
}
