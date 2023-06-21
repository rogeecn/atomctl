package cmd

import (
	"testing"

	"github.com/rogeecn/fabfile"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_astParseRoutes(t *testing.T) {
	Convey("ast parse routes", t, func() {
		path, err := fabfile.Find("tests/routes.go")
		So(err, ShouldBeNil)

		_ = astParseRoutes(path)
	})
}
