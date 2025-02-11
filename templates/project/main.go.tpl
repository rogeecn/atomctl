package main

import (
	"{{.ModuleName}}/app/service/http"

	log "github.com/sirupsen/logrus"
	"go.ipao.vip/atom"
)

// @title           ApiDoc
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   UserName
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @BasePath  /api/v1
// @securityDefinitions.basic  BasicAuth
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	opts := []atom.Option{
		atom.Name("{{ .ProjectName }}"),
		http.Command(),
	}

	if err := atom.Serve(opts...); err != nil {
		log.Fatal(err)
	}
}
