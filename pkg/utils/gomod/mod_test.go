package gomod

import (
	"bufio"
	"os"
	"regexp"
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
		Convey("github.com/redis/go-redis/v9@v9.7.0", func() {
			name, err := getPackageName("github.com/redis/go-redis/v9", "v9.7.0")
			So(err, ShouldBeNil)
			So(name, ShouldEqual, "redis")
		})

		Convey("github.com/pkg/errors@v0.9.1", func() {
			name, err := getPackageName("github.com/pkg/errors", "v0.9.1")
			So(err, ShouldBeNil)
			So(name, ShouldEqual, "errors")
		})
	})
}

func Test_file(t *testing.T) {
	Convey("Test file", t, func() {
		Convey("Test file", func() {
			packagePattern := regexp.MustCompile(`^package\s+(\w+)$`)
			file := "/root/go/pkg/mod/github.com/redis/go-redis/v9@v9.7.0/acl_commands.go"

			// read file line by line
			f, err := os.Open(file)
			So(err, ShouldBeNil)
			defer f.Close()

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if matches := packagePattern.FindStringSubmatch(line); matches != nil {
					t.Logf("Matched package name: %s", matches[1])
				}
			}
			So(scanner.Err(), ShouldBeNil)
		})
	})
}
