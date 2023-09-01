package crud

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"dto/model.go.tpl":        "dto/{filename}.go",
	"dao/model.go.tpl":        "dao/{filename}.go",
	"controller/model.go.tpl": "controller/{filename}.go",
	"service/model.go.tpl":    "service/{filename}.go",
}

var BackendFiles = map[string]string{
	"backend/api.ts.tpl":               "src/api/{module}/{filename}.ts",
	"backend/views/list.vue.tpl":       "src/views/{module}/{filename}/list.vue",
	"backend/views/view.vue.tpl":       "src/views/{module}/{filename}/view.vue",
	"backend/views/create.vue.tpl":     "src/views/{module}/{filename}/create.vue",
	"backend/views/edit.vue.tpl":       "src/views/{module}/{filename}/edit.vue",
	"backend/views/form-items.vue.tpl": "src/views/{module}/{filename}/form-items.vue",
}
