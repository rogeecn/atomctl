package cmd

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.ipao.vip/atomctl/pkg/ast/route"
	"go.ipao.vip/atomctl/pkg/utils/gomod"
)

func CommandGenRoute(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:      "route",
		Short:    "generate routes",
		RunE:     commandGenRouteE,
		PostRunE: commandGenProviderE,
	}

	root.AddCommand(cmd)
}

// https://go.ipao.vip/atomctl/pkg/swag?tab=readme-ov-file#api-operation
func commandGenRouteE(cmd *cobra.Command, args []string) error {
	var err error
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	path, _ = filepath.Abs(path)

	err = gomod.Parse(filepath.Join(path, "go.mod"))
	if err != nil {
		return err
	}

	routes := []route.RouteDefinition{}

	modulePath := filepath.Join(path, "app/http")
	if _, err := os.Stat(modulePath); os.IsNotExist(err) {
		log.Fatal("modules dir not exist, ", modulePath)
	}

	// controllerPattern := regexp.MustCompile(`controller(_?\w+)?\.go`)
	err = filepath.WalkDir(modulePath, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		// if !controllerPattern.MatchString(d.Name()) {
		// 	return nil
		// }
		if strings.HasSuffix(path, ".gen.go") {
			return nil
		}

		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		routes = append(routes, route.ParseFile(path)...)
		return nil
	})
	if err != nil {
		return err
	}

	routeGroups := lo.GroupBy(routes, func(item route.RouteDefinition) string {
		return filepath.Dir(item.Path)
	})

	for path, routes := range routeGroups {
		if err := route.Render(path, routes); err != nil {
			log.WithError(err).WithField("path", path).Error("render route failed")
		}
	}

	return nil
}
