package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load(t *testing.T) {

	tests := []struct {
		name    string
		args    []string
		envVars Config
		want    Config
	}{
		{
			name:    "Defaults",
			args:    nil,
			envVars: Config{RunAddress: ""},
			want: Config{
				RunAddress: defaultRunAddress,
			},
		},
		{
			name:    "Override with flags",
			args:    []string{"-a", "localhost:8888"},
			envVars: Config{RunAddress: ""},
			want: Config{
				RunAddress: "localhost:8888",
			},
		},
		{
			name: "Override with envs",
			args: []string{"-a", "localhost:8888", "-d", "host=local user=u password=p dbname=db", "-r", "localhost:8800", "-i", "2s", "-l", "20"},
			envVars: Config{
				RunAddress: "localhost:8800",
			},
			want: Config{
				RunAddress: "localhost:8800",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("RUN_ADDRESS", tt.envVars.RunAddress)

			c, err := NewConfig("Test", tt.args)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, *c)
		})
	}
}
