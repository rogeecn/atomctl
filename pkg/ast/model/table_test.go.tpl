package model

import (
	"context"
	"testing"

	"{{ .PkgName }}/app/service/testx"
	"{{ .PkgName }}/database"
	"{{ .PkgName }}/database/schemas/public/table"

	. "github.com/smartystreets/goconvey/convey"
	"go.ipao.vip/atom/contracts"

	// . "github.com/go-jet/jet/v2/postgres"
	"github.com/stretchr/testify/suite"
	"go.uber.org/dig"
)

type {{ .PascalTable }}InjectParams struct {
	dig.In
	Initials []contracts.Initial `group:"initials"`
}

type {{ .PascalTable }}TestSuite struct {
	suite.Suite

	{{ .PascalTable }}InjectParams
}

func Test_{{ .PascalTable }}(t *testing.T) {
	providers := testx.Default().With(Provide)
    testx.Serve(providers, t, func(params {{ .PascalTable }}InjectParams) {
		suite.Run(t, &{{ .PascalTable }}TestSuite{
			{{ .PascalTable }}InjectParams: params,
		})
	})
}

func (s *{{ .PascalTable }}TestSuite) Test_Demo() {
	Convey("Test_Demo", s.T(), func() {
		database.Truncate(context.Background(), db, table.{{ .PascalTable }}.TableName())
	})
}
