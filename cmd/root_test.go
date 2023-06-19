package cmd

import "testing"

func Test_modulePathGenerator(t *testing.T) {
	type args struct {
		modulePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1.", args{"user"}, "modules/user"},
		{"2.", args{"user.profile"}, "modules/user/modules/profile"},
		{"3.", args{"user.profile.info"}, "modules/user/modules/profile/modules/info"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := dotToModule(tt.args.modulePath); got != tt.want {
				t.Errorf("modulePathGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}
