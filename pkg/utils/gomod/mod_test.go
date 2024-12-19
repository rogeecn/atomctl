package gomod

import (
	"testing"

	"github.com/rogeecn/fabfile"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_ParseGoMod(t *testing.T) {
	Convey("Test ParseGoMod", t, func() {
		Convey("parse go.mod", func() {
			modFile := fabfile.MustFind("go.mod")
			err := Parse(modFile)
			So(err, ShouldBeNil)

			t.Logf("%+v", goMod)
		})
	})
}

func Test_getPackageName(t *testing.T) {
	Convey("Test getPackageName", t, func() {
		Convey("", func() {
			Convey("github.com/redis/go-redis/v9@v9.7.0", func() {
				name, err := getPackageName("github.com/redis/go-redis/v9", "v9.7.0")
				So(err, ShouldBeNil)
				So(name, ShouldEqual, "redis")
			})
		})
	})
}
