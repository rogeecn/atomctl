package {{.ModuleName}}

import (
	"testing"

	"backend/pkg/service/testx"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"go.uber.org/dig"
)

type ServiceInjectParams struct {
	dig.In
	Svc *Service
}

type ServiceTestSuite struct {
	suite.Suite
	ServiceInjectParams
}

func Test_DiscoverMedias(t *testing.T) {
	providers := testx.Default().With(
		Provide,
	)

	testx.Serve(providers, t, func(params ServiceInjectParams) {
		suite.Run(t, &ServiceTestSuite{ServiceInjectParams: params})
	})
}

func (s *ServiceTestSuite) Test_Service() {
	Convey("Test Service", s.T(), func() {
		So(s.Svc, ShouldNotBeNil)
	})
}
