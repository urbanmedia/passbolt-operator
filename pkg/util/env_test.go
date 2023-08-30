package util

import (
	"os"
	"testing"
)

func TestStringEnvOrDefault(t *testing.T) {
	type args struct {
		key          string
		defaultValue string

		envVars map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1 env empty",
			args: args{
				key:          "test1",
				defaultValue: "test1",
				envVars:      map[string]string{},
			},
			want: "test1",
		},
		{
			name: "test1 env set",
			args: args{
				key:          "test1",
				defaultValue: "test1",
				envVars: map[string]string{
					"test1": "abc",
				},
			},
			want: "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.args.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.args.envVars {
					os.Unsetenv(k)
				}
			}()

			if got := StringEnvOrDefault(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("StringEnvOrDefault() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestIntEnvOrDefault(t *testing.T) {
	type args struct {
		key          string
		defaultValue int

		envVars map[string]string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test1 env not set",
			args: args{
				key:          "number1",
				defaultValue: 1,
				envVars:      map[string]string{},
			},
			want: 1,
		},
		{
			name: "test1 env is set with number",
			args: args{
				key:          "number1",
				defaultValue: 1,
				envVars: map[string]string{
					"number1": "100",
				},
			},
			want: 100,
		},
		{
			name: "test1 env is set with non-number",
			args: args{
				key:          "number1",
				defaultValue: 10,
				envVars: map[string]string{
					"number1": "asdasd",
				},
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.args.envVars {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.args.envVars {
					os.Unsetenv(k)
				}
			}()

			if got := IntEnvOrDefault(tt.args.key, tt.args.defaultValue); got != tt.want {
				t.Errorf("IntEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}
