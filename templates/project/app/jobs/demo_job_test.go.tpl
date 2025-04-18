package jobs

import (
	"context"
	"testing"

	"{{.ModuleName}}/app/service/testx"
	"{{.ModuleName}}/providers/app"
	"{{.ModuleName}}/providers/job"

	. "github.com/riverqueue/river"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	_ "go.ipao.vip/atom"
	"go.ipao.vip/atom/contracts"
	"go.uber.org/dig"
)

type DemoJobSuiteInjectParams struct {
	dig.In

	Initials []contracts.Initial `group:"initials"` // nolint:structcheck
	Job      *job.Job
	App      *app.Config
}

type DemoJobSuite struct {
	suite.Suite

	DemoJobSuiteInjectParams
}

func Test_DemoJob(t *testing.T) {
	providers := testx.Default().With(Provide, models.Provide)

	testx.Serve(providers, t, func(p DemoJobSuiteInjectParams) {
		suite.Run(t, &DemoJobSuite{DemoJobSuiteInjectParams: p})
	})
}

func (t *DemoJobSuite) Test_Work() {
	Convey("test_work", t.T(), func() {
		Convey("step 1", func() {
			job := &Job[DemoJob]{
				Args: DemoJob{
					Strings: []string{"a", "b", "c"},
				},
			}

			worker := &DemoJobWorker{
				job: t.Job,
				app: t.App,
			}

			err := worker.Work(context.Background(), job)
			So(err, ShouldBeNil)
		})
	})
}
