/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"testing"

	"github.com/rogeecn/fabfile"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_processModelTag(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"1.",
			args{`gorm:"column:updated_at;type:timestamp with time zone" json:"updated_at"`},
			"",
		},
		{
			"2.",
			args{`gorm:"column:description;type:text;not null;comment:描述" json:"description"`},
			"描述",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got := processModelTag(tt.args.tag)
			if got != tt.want {
				t.Errorf("processModelTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genCrud(t *testing.T) {
	Convey("gen crud", t, func() {
		path, err := fabfile.Find("tests/crud.go")
		So(err, ShouldBeNil)

		_, err = genCrud("hello", path, "world")
		So(err, ShouldBeNil)
	})
}
