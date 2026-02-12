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
			envVars: Config{RunAddress: "", DatabaseDSN: ""},
			want: Config{
				RunAddress:  defaultRunAddress,
				DatabaseDSN: defaultDatabaseDSN,
			},
		},
		{
			name:    "Override with flags",
			args:    []string{"-a", "localhost:8888"},
			envVars: Config{RunAddress: "", DatabaseDSN: ""},
			want: Config{
				RunAddress:  "localhost:8888",
				DatabaseDSN: defaultDatabaseDSN,
			},
		},
		{
			name: "Override with envs",
			args: []string{"-a", "localhost:8888"},
			envVars: Config{
				RunAddress:  "localhost:8800",
				DatabaseDSN: "",
			},
			want: Config{
				RunAddress:  "localhost:8800",
				DatabaseDSN: defaultDatabaseDSN,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("RUN_ADDRESS", tt.envVars.RunAddress)
			t.Setenv("DATABASE_DSN", tt.envVars.DatabaseDSN)

			c, err := NewConfig("Test", tt.args)

			assert.NoError(t, err)
			assert.Equal(t, tt.want, *c)
		})
	}
}
