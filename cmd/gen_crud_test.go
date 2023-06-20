/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import "testing"

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
