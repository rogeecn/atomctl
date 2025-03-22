package model

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

//go:embed table.go.tpl
var tableTpl string

//go:embed models.gen.go.tpl
var modelTpl string

type TableModelParam struct {
	CamelTable  string // user
	PascalTable string // User
}

func Generate(tables []string, transformer Transformer) error {
	baseDir := "app/models"

	tableTpl := template.Must(template.New("model").Parse(string(tableTpl)))
	modelTpl := template.Must(template.New("modelGen").Parse(string(modelTpl)))

	items := []TableModelParam{}
	for _, table := range tables {
		if lo.Contains(transformer.Ignores.Model, table) {
			log.Printf("[WARN] skip model %s\n", table)
			continue
		}

		items = append(items, TableModelParam{
			CamelTable:  lo.CamelCase(table),
			PascalTable: lo.PascalCase(table),
		})

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

		if err := tableTpl.Execute(fd, map[string]string{"CamelTable": lo.CamelCase(table)}); err != nil {
			return fmt.Errorf("failed to render model template: %w", err)
		}
	}

	// 遍历 baseDir 下的所有文件，将不在 tables 中的文件名（不带扩展名）加入
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return fmt.Errorf("遍历目录 %s 失败: %w", baseDir, err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if strings.HasSuffix(name, ".gen.go") {
			continue
		}

		if strings.HasSuffix(name, "_test.go") {
			continue
		}

		baseName := strings.TrimSuffix(name, filepath.Ext(name))
		if lo.Contains(transformer.Ignores.Model, baseName) {
			log.Printf("[WARN] skip model %s\n", baseName)
			continue
		}

		if !lo.Contains(tables, baseName) {
			items = append(items, TableModelParam{
				CamelTable:  lo.CamelCase(baseName),
				PascalTable: lo.PascalCase(baseName),
			})
		}
	}

	// 渲染总的 model 文件

	modelFile := fmt.Sprintf("%s/models.gen.go", baseDir)
	os.Remove(modelFile)
	fd, err := os.Create(modelFile)
	if err != nil {
		return fmt.Errorf("failed to create model file %s: %w", baseDir, err)
	}
	defer fd.Close()

	if err := modelTpl.Execute(fd, items); err != nil {
		return fmt.Errorf("failed to render model template: %w", err)
	}

	return nil
}
