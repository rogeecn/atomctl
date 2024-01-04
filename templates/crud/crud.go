package crud

import (
	"embed"
)

//go:embed tpl/*
var Templates embed.FS

var Files = map[string]string{
	"model_dto.go.tpl":        "{filename}_dto.go",
	"model_controller.go.tpl": "{filename}_controller.go",
	"model_service.go.tpl":    "{filename}_service.go",
}

var BackendFiles = map[string]string{
	"backend/api.ts.tpl":               "src/api/{module}/{filename}.ts",
	"backend/views/list.vue.tpl":       "src/views/{module}/{filename}/list.vue",
	"backend/views/view.vue.tpl":       "src/views/{module}/{filename}/view.vue",
	"backend/views/create.vue.tpl":     "src/views/{module}/{filename}/create.vue",
	"backend/views/edit.vue.tpl":       "src/views/{module}/{filename}/edit.vue",
	"backend/views/form-items.vue.tpl": "src/views/{module}/{filename}/form-items.vue",
}
