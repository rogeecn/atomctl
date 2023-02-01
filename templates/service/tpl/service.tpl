package service

import (
	"context"
)

type {{.PascalName}}Service interface {
	GetName(ctx context.Context) (string, error)
}

type {{.CamelName}}Service struct {
}

func New{{.PascalName}}Service() {{.PascalName}}Service {
	return &{{.CamelName}}Service{}
}

func (svc *{{.CamelName}}Service) GetName(ctx context.Context) (string, error) {
	return "{{.CamelName}}.GetName", nil
}
