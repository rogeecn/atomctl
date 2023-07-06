/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func init() {
	genCmd.AddCommand(genActionCmd)
}

// action:Benchmark "关联 Nodes 列表" json get:/v1/assets/kubernetes/clusters/{account_id}/images/{id}/pods "资产中心-K8S集群-镜像"
// action:[Method] "Description" [json|text] [method] [route] [tag]

var genActionCmd = &cobra.Command{
	Use:   "action",
	Short: "generate actions for controller",
	Long: `
	Generate actions for controller by directive in comments blow

	// action:[Method] "Description" [json|text] [method] [route] [tag]

	`,
	Example: "atomctl gen action [filename]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("args invalid, please check")
		}
		filename := args[0]

		tpl, err := template.New("action").Parse(actionTpl)
		if err != nil {
			return err
		}

		renders, err := parseControllerActionTpl(filename)
		if err != nil {
			return err
		}

		c, err := os.ReadFile(filename)
		if err != nil {
			return err
		}
		content := string(c)

		for _, render := range renders {
			var buf bytes.Buffer
			if err := tpl.Execute(&buf, render); err != nil {
				return err
			}

			content = strings.Replace(content, render.ActionTplString, buf.String(), 1)
		}
		return os.WriteFile(filename, []byte(content), os.ModePerm)
	},
}

type RenderAction struct {
	ActionTplString string
	Method          string
	Controller      string
	Description     string
	Tag             string
	ContentType     string
	Params          []RenderActionParams
	Route           string
	RouteMethod     string
}
type RenderActionParams struct {
	Name      string
	NameCamel string
	Type      string
	Desc      string
}

var actionTpl = `
// {{ .Method }} {{.Description }}
//
//	@Summary		{{.Description }}
//	@Description	{{.Description }}
//	@Tags			{{ .Tag }}
//	@Accept			{{ .ContentType }}
//	@Produce		{{ .ContentType }}
{{- range $index, $param := .Params }}
//	@Param			{{ $param.Name }}	path		{{ $param.Type }}	true	"{{ $param.Desc }}"
{{- end }}
//	@Success		200			{string}	TODO:AddData
//	@Router			{{ .Route }} [{{ .RouteMethod }}]
func (c *{{ .Controller }}) {{ .Method }}(ctx *fiber.Ctx{{range $index, $param := .Params }}, {{ $param.NameCamel }} {{ $param.Type }}{{ end }}) error {
	panic("not implemented")
}

`

func parseControllerActionTpl(filename string) ([]RenderAction, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Println("ERR: ", err)
		return nil, err
	}

	var controller string
	actions := []RenderAction{}

	for _, decl := range node.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			for _, spec := range decl.(*ast.GenDecl).Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					controller = spec.(*ast.TypeSpec).Name.Name
				default:
					continue
				}
			}
		default:
			continue
		}
	}

	if controller == "" {
		return nil, errors.New("invalid controller")
	}

	parsePathFields := func(route string) []RenderActionParams {
		pattern := regexp.MustCompile(`\{(.*?)\}`)
		fields := pattern.FindAllStringSubmatch(route, -1)

		params := []RenderActionParams{}
		for _, field := range fields {
			fieldName := field[1]
			typ := "int"
			if !strings.HasSuffix(strings.ToLower(fieldName), "id") {
				typ = "string"
			}
			params = append(params, RenderActionParams{
				Name:      fieldName,
				NameCamel: strcase.ToLowerCamel(fieldName),
				Type:      typ,
				Desc:      strcase.ToCamel(fieldName),
			})
		}
		return params
	}

	parseAction := func(comment string) (RenderAction, error) {
		payload := comment[len("action:"):]
		r := regexp.MustCompile(`[^\s"]+|"([^"]*)"`)
		directives := r.FindAllString(payload, -1)
		if len(directives) != 6 {
			return RenderAction{}, errors.New("invalid action, missing directives")
		}

		return RenderAction{
			Method:      directives[0],
			Description: strings.Trim(directives[1], `"`),
			ContentType: directives[2],
			RouteMethod: directives[3],
			Route:       directives[4],
			Tag:         strings.Trim(directives[5], `"`),
			Params:      parsePathFields(directives[4]),
		}, nil
	}

	for _, decl := range node.Comments {
		for _, comment := range decl.List {
			text := strings.TrimLeft(comment.Text, "/ ")

			if strings.HasPrefix(text, "action:") {
				log.Printf("parsing: %s", comment.Text)
				action, err := parseAction(text)
				if err != nil {
					return nil, err
				}
				action.Controller = controller
				action.ActionTplString = comment.Text
				actions = append(actions, action)
			}
		}
	}
	return actions, nil
}
