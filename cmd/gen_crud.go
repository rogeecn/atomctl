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
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/rogeecn/atomctl/templates/crud"
	"github.com/rogeecn/atomctl/utils"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var (
	flagForce   bool
	backendDest string
	title       string
	moduleTitle string
)

func init() {
	genCmd.AddCommand(genCrudCmd)
	genCrudCmd.Flags().String("route", "", "manually define route path")
	genCrudCmd.Flags().String("tag", "DEFAULT_TAG_NAME", "define swagger tag")
	genCrudCmd.Flags().BoolVar(&flagForce, "force", false, "overwrite file if exists")
	genCrudCmd.Flags().StringVar(&backendDest, "backend", "", "generate backend")
	genCrudCmd.Flags().StringVar(&title, "title", "", "generate title")
	genCrudCmd.Flags().StringVar(&moduleTitle, "module-title", "", "generate module title")
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
		modelInfo.Filename = args[0]

		// render go files
		render := CrudRenderParams{
			PkgName: pkgName,
			Module:  modulePath,
			Model:   modelInfo,
		}
		generateFiles, err := render.prepareFiles(crud.Files, args[0], flagForce)
		if err != nil {
			return err
		}

		if err := utils.Generate(generateFiles, crud.Templates, render); err != nil {
			return err
		}

		log.Println("generate crud success")

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
				provider := fmt.Sprintf("\n\tif err := container.Container.Provide(%s); err!=nil {\n\t\treturn err\n\t}\n\treturn nil", providerFunc)
				content = strings.Replace(content, "return nil", provider, 1)
			}

			containerPackage := `"github.com/rogeecn/atom/container"`
			if !strings.Contains(content, containerPackage) {
				content = strings.Replace(content, "import (", "import (\n\t"+containerPackage, 1)
			}
			log.Printf("[Info] add new provider: %s for %s", providerFunc, providerFile)

			_ = os.WriteFile(providerFile, []byte(content), os.ModePerm)
		}

		if backendDest == "" {
			return nil
		}

		// render vue template backend files
		backendRender := &CrudRenderParams{
			PkgName: pkgName,
			Module:  backendDest,
			Model:   modelInfo,
			Vars: map[string]string{
				"module": strings.ReplaceAll(args[1], ".", "/"),
			},
			Title:       title,
			ModuleTitle: moduleTitle,
		}
		generateFiles, err = backendRender.prepareFiles(crud.BackendFiles, args[0], flagForce)
		if err != nil {
			return err
		}

		if err := utils.Generate(generateFiles, crud.Templates, backendRender); err != nil {
			return err
		}
		log.Println("generate to ", backendDest, " success")

		outputRoutes(backendRender)
		return nil
	},
}

func outputRoutes(render *CrudRenderParams) {
	t := `
	// {{ .Vars.moduleTitle }}
    { path: '{{ .Model.Filename }}', name: '{{ .Model.Name }}List', component: () => import('@/views/{{ .Vars.module }}/{{ .Model.Filename }}/list.vue'), meta: { title: '{{ .Title }}列表', requiresAuth: true }, },
    { path: '{{ .Model.Filename }}/view/:id', name: '{{ .Model.Name }}View', component: () => import('@/views/{{ .Vars.module }}/{{ .Model.Filename }}/view.vue'), meta: { title: '{{ .Title }}详情', requiresAuth: true, hideInMenu: true } },
    { path: '{{ .Model.Filename }}/edit/:id', name: '{{ .Model.Name }}Edit', component: () => import('@/views/{{ .Vars.module }}/{{ .Model.Filename }}/edit.vue'), meta: { title: '{{ .Title }}编辑', requiresAuth: true, hideInMenu: true } },
    { path: '{{ .Model.Filename }}/create', name: '{{ .Model.Name }}Create', component: () => import('@/views/{{ .Vars.module }}/{{ .Model.Filename }}/create.vue'), meta: { title: '{{ .Title }}创建', requiresAuth: true, hideInMenu: true } },
	`

	buf := bytes.NewBuffer(nil)
	tpl, err := template.New("Route").Parse(t)
	if err != nil {
		fmt.Println("ERR: ", err)
		return
	}

	if err := tpl.Execute(buf, render); err != nil {
		fmt.Println("ERR: ", err)
		return
	}

	fmt.Printf("\n\n%s\n", buf.Bytes())
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
	PkgName     string
	Module      string
	Model       *ModelInfo
	Title       string
	ModuleTitle string
	Vars        map[string]string
}

func (m *CrudRenderParams) prepareFiles(files map[string]string, filename string, force bool) (map[string]string, error) {
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
			t, err := template.New("tpl_" + tpl).Parse(target)
			if err != nil {
				return nil, errors.Wrap(err, "init template failed")
			}
			fmt.Println("-----", m)
			if err := t.Execute(newName, m); err != nil {
				return nil, errors.Wrapf(err, "generate target file failed, tpl: %s", tpl)
			}
			target = newName.String()
		}
		tplFilePath := "tpl/" + tpl

		target = strings.Replace(target, "{filename}", filename, -1)
		if m.Vars != nil {
			for k, v := range m.Vars {
				target = strings.Replace(target, "{"+k+"}", v, -1)
			}
		}
		target = filepath.Join(m.Module, target)
		result[tplFilePath] = target
		if utils.IsFile(target) && !force {
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
	Imports    []string
	Filename   string
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
	Name         string
	Type         string
	Tag          string
	Comment      string
	Package      string
	PackageAlias string
}

func parseModelInfo(file *ast.File) ModelInfo {
	modelInfo := ModelInfo{}
	fields := []ModelField{}
	vPattern := regexp.MustCompile(`^v\d+$`)
	imports := make(map[string]string)
	for _, imp := range file.Imports {

		name := ""
		if imp.Name == nil {
			paths := strings.Split(strings.Trim(imp.Path.Value, "\""), "/")
			name = paths[len(paths)-1]
			if vPattern.MatchString(name) {
				name = paths[len(paths)-2]
			}
		} else {
			name = imp.Name.Name
		}

		if name == "_" {
			paths := strings.Split(strings.Trim(imp.Path.Value, "\""), "/")
			name = paths[len(paths)-1]
			if vPattern.MatchString(name) {
				name = paths[len(paths)-2]
			}
			if _, ok := imports[name]; ok {
				name = fmt.Sprintf("%s%d", name, rand.Intn(100))
			}
		}
		pkg := strings.Trim(imp.Path.Value, `"`)
		imports[name] = pkg
	}

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

				var pkg string
				var pkgAlias string
				var typ string
				switch field.Type.(type) {
				case *ast.Ident:
					typ = field.Type.(*ast.Ident).Name
				case *ast.IndexExpr:
					pkgAlias = field.Type.(*ast.IndexExpr).Index.(*ast.SelectorExpr).X.(*ast.Ident).Name
					p, ok := imports[pkgAlias]
					if !ok {
						continue
					}
					pkg = p
					typ = field.Type.(*ast.IndexExpr).Index.(*ast.SelectorExpr).Sel.Name
				case *ast.StarExpr:
					paramsType := field.Type.(*ast.StarExpr)
					switch paramsType.X.(type) {
					case *ast.SelectorExpr:
						X := paramsType.X.(*ast.SelectorExpr)

						pkgAlias = X.X.(*ast.Ident).Name
						p, ok := imports[pkgAlias]
						if !ok {
							continue
						}
						pkg = p

						typ = X.Sel.Name
					default:
						typ = paramsType.X.(*ast.Ident).Name
					}
				case *ast.SelectorExpr:
					pkgAlias = field.Type.(*ast.SelectorExpr).X.(*ast.Ident).Name
					p, ok := imports[pkgAlias]
					if !ok {
						continue
					}
					pkg = p
					typ = field.Type.(*ast.SelectorExpr).Sel.Name
				}

				tag, comment := processModelTag(tag)
				fields = append(fields, ModelField{
					Name:         field.Names[0].Name,
					Tag:          tag,
					Comment:      comment,
					Type:         typ,
					Package:      pkg,
					PackageAlias: pkgAlias,
				})
			}
			modelInfo.Fields = fields
			modelInfo.Imports = lo.Uniq(lo.FilterMap(fields, func(field ModelField, _ int) (string, bool) {
				if field.Package == "" {
					return "", false
				}
				if field.PackageAlias == "" {
					return fmt.Sprintf("%q", field.Package), true
				}

				if field.PackageAlias == field.Package {
					return fmt.Sprintf("%q", field.Package), true
				}

				if strings.HasSuffix(field.Package, "/"+field.PackageAlias) {
					return fmt.Sprintf("%q", field.Package), true
				}

				return fmt.Sprintf("%s %q", field.PackageAlias, field.Package), true
			}))
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
	tag = patternJson.FindStringSubmatch(tag)[1]

	return tag, comment
}
