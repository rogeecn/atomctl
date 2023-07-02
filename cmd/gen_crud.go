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
	"os"
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

func init() {
	genCmd.AddCommand(genCrudCmd)
	genCrudCmd.Flags().String("route", "", "manually define route path")
	genCrudCmd.Flags().String("tag", "DEFAULT_TAG_NAME", "define swagger tag")
}

var genCrudCmd = &cobra.Command{
	Use:     "crud",
	Short:   "generate crud for model",
	Example: "atomctl gen crud [model_file_without_gen.go] [module_name_with_dot]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("args invalid, please check")
		}
		pkgName := getPackage()
		modulePath, _ := dotToModule(args[1])
		modelFile := fmt.Sprintf("database/models/%s.gen.go", args[0])

		modelInfo, err := genCrud(pkgName, modelFile, modulePath)
		if err != nil {
			return err
		}
		modelInfo.TagName = cmd.Flag("tag").Value.String()
		modelInfo.RouteName = strcase.ToSnake(args[0])
		if cmd.Flag("route").Value.String() != "" {
			modelInfo.RouteName = strings.Trim(cmd.Flag("route").Value.String(), "/")
		}
		if strings.Contains(modelInfo.RouteName, "{id}") {
			return errors.New("Invalid route, route should not contains {id}")
		}
		modelInfo.GuessIntType()
		modelInfo.parsePathFields()

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

		for _, file := range generateFiles {
			dirname := filepath.Base(filepath.Dir(file))
			if dirname == "dto" {
				continue
			}

			providerFile := filepath.Join(filepath.Dir(file), "provider.go")
			if !utils.IsFile(providerFile) {
				log.Println("[Warn] " + providerFile + " not exists, please add new provider manually")
				continue
			}

			providerContent, err := os.ReadFile(providerFile)
			if err != nil {
				log.Println("[Warn] read " + providerFile + " failed, please add new provider manually")
				continue
			}

			suffix := strcase.ToCamel(dirname)

			providerFunc := fmt.Sprintf("New%s%s", modelInfo.Name, suffix)

			content := string(providerContent)
			if !strings.Contains(content, providerFunc) {
				provider := fmt.Sprintf("if err := container.Container.Provide(%s); err!=nil {\n\treturn err\n\t}\n\treturn nil", providerFunc)
				content = strings.Replace(content, "return nil", provider, 1)
			}

			containerPackage := `"github.com/rogeecn/atom/container"`
			if !strings.Contains(content, containerPackage) {
				content = strings.Replace(content, "import (", "import (\n\t"+containerPackage, 1)
			}
			log.Printf("[Info] add new provider: %s for %s", providerFunc, providerFile)
			_ = os.WriteFile(providerFile, []byte(content), os.ModePerm)
		}

		return nil
	},
}

func genCrud(pkgName, modelFile, moduleName string) (*ModelInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, modelFile, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	modelInfo := parseModelInfo(node)
	return &modelInfo, nil
}

type CrudRenderParams struct {
	PkgName string
	Module  string
	Model   *ModelInfo
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
	Name       string
	CamelName  string
	RouteName  string
	TagName    string
	IntType    string
	Fields     []ModelField
	PathFields []ModelField
}

func (m *ModelInfo) GuessIntType() {
	m.IntType = "int"
	for _, f := range m.Fields {
		if f.Name == "ID" {
			m.IntType = strings.TrimLeft(f.Type, "*")
			return
		}
	}
}

func (m *ModelInfo) parsePathFields() {
	pattern := regexp.MustCompile(`\{(.*?)\}`)
	fields := pattern.FindAllStringSubmatch(m.RouteName, -1)

	for _, field := range fields {
		fieldName := field[1]
		typ := "int"
		if !strings.HasSuffix(strings.ToLower(fieldName), "id") {
			typ = "string"
		}
		m.PathFields = append(m.PathFields, ModelField{
			Name:    strcase.ToLowerCamel(field[1]),
			Type:    typ,
			Comment: strcase.ToCamel(fieldName),
		})
	}
}

type ModelField struct {
	Name    string
	Type    string
	Tag     string
	Comment string
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
				if field.Names[0].Name == "DeletedAt" {
					continue
				}

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

				tag, comment := processModelTag(tag)
				fields = append(fields, ModelField{
					Name:    field.Names[0].Name,
					Type:    "*" + typ,
					Tag:     tag,
					Comment: comment,
				})
			}
			modelInfo.Fields = fields
			return modelInfo
		}
	}
	return modelInfo
}

// get field tag comment
func processModelTag(tag string) (string, string) {
	comment := ""
	patternComment := regexp.MustCompile(`gorm:".*?;comment:(.*?);?"\s+`)
	if patternComment.MatchString(tag) {
		comment = patternComment.FindStringSubmatch(tag)[1]
	}

	patternTag := regexp.MustCompile(`gorm:".*?"\s+`)
	if !patternTag.MatchString(tag) {
		return tag, ""
	}
	tag = patternTag.ReplaceAllString(tag, "")

	patternJson := regexp.MustCompile(`json:"(.*?)"`)
	tag = patternJson.ReplaceAllString(tag, `form:"$1" query:"$1" json:"$1,omitempty"`)

	return tag, comment
}
