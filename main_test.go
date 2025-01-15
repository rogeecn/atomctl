package main

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_router(t *testing.T) {
	Convey("Test routerPattern", t, func() {
		jsonReg := regexp.MustCompile(`Json\[\[?\]?(\w+)\]`)
		items := []string{
			"Json[abc]",
			"Json[[]abc]",
		}

		types := []string{
			"string",
			"int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64",
			"bool",
		}
		for _, item := range items {
			match := jsonReg.FindStringSubmatch(item)
			if len(match) ==2 && !lo.Contains(types, match[1]) { {

			}
		}
	})
}
