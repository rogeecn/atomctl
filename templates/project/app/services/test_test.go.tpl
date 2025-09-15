package services

import (
	"testing"
	"time"

	"{{.ModuleName}}/app/commands/testx"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"

	_ "go.ipao.vip/atom"
	"go.ipao.vip/atom/contracts"
	"go.uber.org/dig"
)

type TestSuiteInjectParams struct {
	dig.In

	Initials []contracts.Initial `group:"initials"` // nolint:structcheck
}

type TestSuite struct {
	suite.Suite

	TestSuiteInjectParams
}

func Test_Test(t *testing.T) {
	providers := testx.Default().With(Provide)

	testx.Serve(providers, t, func(p TestSuiteInjectParams) {
		suite.Run(t, &TestSuite{TestSuiteInjectParams: p})
	})
}

func (t *TestSuite) Test_Test() {
	Convey("test_work", t.T(), func() {
		t.T().Log("start test at", time.Now())
	})
}
