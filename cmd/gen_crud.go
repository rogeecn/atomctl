/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/rogeecn/atomctl/templates/crud"
	"github.com/rogeecn/atomctl/utils"
	"github.com/spf13/cobra"
)

var genCrudCmd = &cobra.Command{
	Use:     "crud",
	Short:   "generate crud for model",
	Example: "atomctl gen crud [model_file_without_gen.go] [module_name_with_dot]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("args invalid, please check")
		}
		pkgName := getPackage()
		modulePath := strings.Join(append([]string{""}, strings.Split(args[1], ".")...), "modules/")
		modelFile := fmt.Sprintf("database/models/%s.gen.go", args[0])

		modelInfo, err := genCrud(pkgName, modelFile, modulePath)
		if err != nil {
			return err
		}
		modelInfo.RouteName = args[0]

		render := CrudRenderParams{
			PkgName: pkgName,
			Module:  modulePath,
			Model:   modelInfo,
		}

		generateFiles, err := render.prepareFiles(crud.Files, args[0])
		if err != nil {
			return err
		}

		if err := utils.Generate(generateFiles, crud.Templates, render); err != nil {
			return err
		}

		log.Println("generate crud success")
		log.Println("REMEMBER TO ADD NEW PROVIDERS")
		return nil
	},
}

func init() {
	genCmd.AddCommand(genCrudCmd)
}

func genCrud(pkgName, modelFile, moduleName string) (ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, modelFile, nil, parser.ParseComments)
	if err != nil {
		return ModelInfo{}, err
	}

	modelInfo := parseModelInfo(node)
	return modelInfo, nil
}

type CrudRenderParams struct {
	PkgName string
	Module  string
	Model   ModelInfo
}

func (m *CrudRenderParams) prepareFiles(files map[string]string, filename string) (map[string]string, error) {
	result := make(map[string]string)
	for tpl, target := range files {
		// get target file name
		if len(target) == 0 {
			if utils.IsTplFile(tpl) {
				target = utils.TplToGo(tpl)
			} else {
				target = tpl
			}
		} else {
			newName := bytes.NewBuffer(nil)
			t, err := template.New("name").Parse(target)
			if err != nil {
				return nil, errors.Wrap(err, "init template failed")
			}
			if err := t.Execute(newName, m); err != nil {
				return nil, errors.Wrapf(err, "generate target file failed, tpl: %s", tpl)
			}
			target = newName.String()
		}
		tplFilePath := "tpl/" + tpl

		target = filepath.Join(m.Module, strings.Replace(target, "{filename}", filename, -1))
		result[tplFilePath] = target
		if utils.IsFile(target) {
			return nil, errors.New(target + " file exists")
		}
	}

	return result, nil
}

type ModelInfo struct {
	Name      string
	CamelName string
	RouteName string
	Fields    []ModelField
}

type ModelField struct {
	Name string
	Type string
	Tag  string
}

func parseModelInfo(file *ast.File) ModelInfo {
	modelInfo := ModelInfo{}
	fields := []ModelField{}
	for _, decl := range file.Decls {
		switch decl.(type) {
		case *ast.GenDecl:
			if decl.(*ast.GenDecl).Tok != token.TYPE {
				continue
			}
			spec := decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec)
			modelInfo.Name = spec.Name.Name
			modelInfo.CamelName = strcase.ToLowerCamel(modelInfo.Name)

			for _, field := range spec.Type.(*ast.StructType).Fields.List {
				tag := ""
				if field.Tag != nil {
					tag = field.Tag.Value
				}

				var typ string
				switch field.Type.(type) {
				case *ast.Ident:
					typ = field.Type.(*ast.Ident).Name
				case *ast.StarExpr:
					paramsType := field.Type.(*ast.StarExpr)
					X := paramsType.X.(*ast.SelectorExpr)
					typ = fmt.Sprintf("%s.%s", X.X.(*ast.Ident).Name, X.Sel.Name)
				case *ast.SelectorExpr:
					typ = fmt.Sprintf("%s.%s", field.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name, field.Type.(*ast.SelectorExpr).Sel.Name)
				}
				fields = append(fields, ModelField{
					Name: field.Names[0].Name,
					Type: "*" + typ,
					Tag:  processModelTag(tag),
				})
			}
			modelInfo.Fields = fields
			return modelInfo
		}
	}
	return modelInfo
}

func processModelTag(tag string) string {
	pattern := regexp.MustCompile(`gorm:".*?"\s+`)
	if pattern.MatchString(tag) {
		return pattern.ReplaceAllString(tag, "")
	}
	return tag
}
