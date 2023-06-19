package cmd

import "testing"

func Test_guessTableName(t *testing.T) {
	type args struct {
		migrationName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1.", args{"create_table"}, "table"},
		{"2.", args{"create_table_sub"}, "table_sub"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := guessTableName(tt.args.migrationName); got != tt.want {
				t.Errorf("guessTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}
