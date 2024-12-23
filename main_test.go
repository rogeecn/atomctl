package main

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_router(t *testing.T) {
	routerPattern := regexp.MustCompile(`^(/[\w./\-{}\(\)+:$]*)[[:blank:]]+\[(\w+)]`)

	Convey("Test routerPattern", t, func() {
		Convey("Pattern 1", func() {
			commentLine := "/api/v1/health [GET] # Check health status"
			matches := routerPattern.FindStringSubmatch(commentLine)
			t.Logf("matches: %v", matches)
		})

		Convey("Pattern 2", func() {
			commentLine := "/api/v1/:health [get] "
			matches := routerPattern.FindStringSubmatch(commentLine)
			t.Logf("matches: %v", matches)
		})

		Convey("Pattern 3", func() {
			commentLine := "/api/v1/get_users-:id [get] "
			pattern := regexp.MustCompile(`<.*?>`)
			commentLine = pattern.ReplaceAllString(commentLine, "")

			matches := routerPattern.FindStringSubmatch(commentLine)
			t.Logf("matches: %v", matches)
		})

		Convey("Pattern 4", func() {
			commentLine := "/api/v1/get_users-:id<int>/name/:name<string> [get] "
			pattern := regexp.MustCompile(`:(\w+)(<.*?>)?`)
			commentLine = pattern.ReplaceAllString(commentLine, "{$1}")

			matches := routerPattern.FindStringSubmatch(commentLine)
			t.Logf("matches: %v", matches)
		})
	})
}
