package service

import (
	"context"
)

type {{.PascalName}}Service struct {
}

func New{{.PascalName}}Service() *{{.PascalName}}Service {
	return &{{.PascalName}}Service{}
}

func (svc *{{.PascalName}}Service) GetName(ctx context.Context) (string, error) {
	return "{{.PascalName}}.GetName", nil
}
