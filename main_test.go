package main

import (
	"strings"
	"testing"
)

type ParamDefinition struct {
	Name     string
	Type     string
	Key      string
	Table    string
	Model    string
	Position string
}

func parseBind(bind string) ParamDefinition {
	var param ParamDefinition
	parts := strings.FieldsFunc(bind, func(r rune) bool {
		return r == ' ' || r == '(' || r == ')'
	})

	// 过滤掉空的元素
	var newParts []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			newParts = append(newParts, part)
		}
	}

	for i, part := range parts {
		switch part {
		case "@Bind":
			param.Name = parts[i+1]
			param.Position = parts[i+2]
		case "key":
			param.Key = parts[i+1]
		case "table":
			param.Table = parts[i+1]
		case "model":
			param.Model = parts[i+1]
		}
	}
	return param
}

func Test_T(t *testing.T) {
	// @Bind [Name] [Type] [Key] [Table] [Model]
	suites := []string{
		`@Bind name query key("a") table(b) model("c")`,
		`@Bind id uri key(a)`,
		`@Bind id uri table(b)`,
		`@Bind id uri key(b) model(c)`,
		`@Bind id uri key(b) table(c)`,
		`@Bind id uri table(b) key(c)`,
		`@Bind id uri`,
	}

	for _, suite := range suites {
		param := parseBind(suite)
		t.Logf("Parsed Param: %+v", param)
	}
}
