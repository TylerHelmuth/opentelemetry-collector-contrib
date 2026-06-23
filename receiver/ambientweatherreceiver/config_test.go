// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ambientweatherreceiver

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/confmap/xconfmap"
)

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	sub, err := cm.Sub("ambientweather")
	require.NoError(t, err)

	cfg := createDefaultConfig()
	require.NoError(t, sub.Unmarshal(cfg))
	require.NoError(t, xconfmap.Validate(cfg))

	got := cfg.(*Config)
	assert.Equal(t, "0.0.0.0:6255", got.NetAddr.Endpoint)
	assert.Equal(t, "/custom/report/", got.Path)
	assert.Equal(t, "backyard", got.StationName)
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*Config)
		wantErr string
	}{
		{
			name:   "valid default",
			modify: func(*Config) {},
		},
		{
			name:    "empty endpoint",
			modify:  func(c *Config) { c.NetAddr.Endpoint = "" },
			wantErr: "endpoint cannot be empty",
		},
		{
			name:    "path missing leading slash",
			modify:  func(c *Config) { c.Path = "data/report" },
			wantErr: "must begin with '/'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := createDefaultConfig().(*Config)
			tt.modify(cfg)
			err := cfg.Validate()
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}
