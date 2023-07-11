package run_migrate

import (
	_ "embed"
)

//go:embed tpl/up.tpl
var Up string

//go:embed tpl/down.tpl
var Down string
