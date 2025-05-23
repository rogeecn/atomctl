package model

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	log "github.com/sirupsen/logrus"
	"go.ipao.vip/atomctl/pkg/utils/gomod"
)

//go:embed table.go.tpl
var tableTpl string

//go:embed table_test.go.tpl
var tableTestTpl string

//go:embed provider.gen.go.tpl
var providerTplStr string

type TableModelParam struct {
	PkgName     string
	CamelTable  string // user
	PascalTable string // User
}

func Generate(tables []string, transformer Transformer) error {
	baseDir := "app/model"
	modelDir := "database/schemas/public/model"
	tableDir := "database/schemas/public/table"
	defer func() {
		os.RemoveAll("database/schemas")
	}()

	os.RemoveAll("database/table")
	// move tableDir to database/table
	if err := os.Rename(tableDir, "database/table"); err != nil {
		return err
	}

	// remove all files in app/model with ext .gen.go
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".gen.go") {
			if err := os.RemoveAll(filepath.Join(baseDir, file.Name())); err != nil {
				return err
			}
		}
	}

	// move files remove ext .go to .gen.go
	files, err = os.ReadDir(modelDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		// get filename without ext
		name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))

		from := filepath.Join(modelDir, file.Name())
		to := filepath.Join(baseDir, name+".gen.go")
		log.Infof("Move %s to %s", from, to)
		if err := os.Rename(from, to); err != nil {
			return err
		}
	}

	// remove database/schemas/public/model
	if err := os.RemoveAll(modelDir); err != nil {
		return err
	}

	tableTpl := template.Must(template.New("model").Parse(string(tableTpl)))
	tableTestTpl := template.Must(template.New("model").Parse(string(tableTestTpl)))
	providerTpl := template.Must(template.New("modelGen").Parse(string(providerTplStr)))

	items := []TableModelParam{}
	for _, table := range tables {
		if lo.Contains(transformer.Ignores.Model, table) {
			log.Printf("[WARN] skip model %s\n", table)
			continue
		}

		tableInfo := TableModelParam{
			CamelTable:  lo.CamelCase(table),
			PascalTable: lo.PascalCase(table),
			PkgName:     gomod.GetModuleName(),
		}
		items = append(items, tableInfo)

		modelFile := fmt.Sprintf("%s/%s.go", baseDir, table)
		// 如果 modelFile 已存在，则跳过
		if _, err := os.Stat(modelFile); err == nil {
			fmt.Printf("Model file %s already exists. Skipping...\n", modelFile)
			continue
		}

		// 如果 modelFile 不存在，则创建
		fd, err := os.Create(modelFile)
		if err != nil {
			return fmt.Errorf("failed to create model file %s: %w", modelFile, err)
		}
		defer fd.Close()

		if err := tableTpl.Execute(fd, tableInfo); err != nil {
			return fmt.Errorf("failed to render model template: %w", err)
		}

		modelTestFile := fmt.Sprintf("%s/%s_test.go", baseDir, table)
		// 如果 modelTestFile 已存在，则跳过
		if _, err := os.Stat(modelTestFile); err == nil {
			fmt.Printf("Model test file %s already exists. Skipping...\n", modelTestFile)
			continue
		}

		// 如果 modelTestFile 不存在，则创建
		fd, err = os.Create(modelTestFile)
		if err != nil {
			return fmt.Errorf("failed to create model test file %s: %w", modelTestFile, err)
		}
		defer fd.Close()

		if err := tableTestTpl.Execute(fd, tableInfo); err != nil {
			return fmt.Errorf("failed to render model test template: %w", err)
		}
	}

	// 渲染总的 provider 文件
	providerFile := fmt.Sprintf("%s/provider.gen.go", baseDir)
	os.Remove(providerFile)
	fd, err := os.Create(providerFile)
	if err != nil {
		return fmt.Errorf("failed to create provider file %s: %w", providerFile, err)
	}
	defer fd.Close()

	if err := providerTpl.Execute(fd, items); err != nil {
		return fmt.Errorf("failed to render model template: %w", err)
	}

	return nil
}

func addProviderComment(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if strings.Contains(string(content), "// @provider") {
		return nil
	}

	// Write this comment to the up line of the type xxx struct
	newLines := []string{}
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.Contains(line, "type ") && strings.Contains(line, "struct") {
			newLines = append(newLines, "// @provider")
			// append rest lines
			newLines = append(newLines, lines[i:]...)
			break
		}
		newLines = append(newLines, line)
	}
	newContent := strings.Join(newLines, "\n")
	if _, err := file.WriteAt([]byte(newContent), 0); err != nil {
		return err
	}
	return nil
}
