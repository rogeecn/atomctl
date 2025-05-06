package conds

import (
	. "github.com/go-jet/jet/v2/postgres"
)

type Cond func(BoolExpression) BoolExpression

func Default() BoolExpression {
	return BoolExp(Bool(true))
}
