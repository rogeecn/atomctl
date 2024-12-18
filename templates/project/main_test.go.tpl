package main

import (
	"testing"

	"{{.ModuleName}}/pkg/service/model"

	"git.ipao.vip/rogeecn/atom"
)

func Test_GenModel(t *testing.T) {
	err := atom.Serve(model.Options()...)
	if err != nil {
		t.Fatal(err)
	}
}
