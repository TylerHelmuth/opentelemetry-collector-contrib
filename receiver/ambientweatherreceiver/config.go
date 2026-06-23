// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ambientweatherreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver"

import (
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/config/confighttp"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver/internal/metadata"
)

// Config defines configuration for the Ambient Weather receiver.
type Config struct {
	confighttp.ServerConfig       `mapstructure:",squash"`
	metadata.MetricsBuilderConfig `mapstructure:",squash"`

	// Path is the URL path the weather console is configured to upload to.
	Path string `mapstructure:"path"`

	// StationName is an optional friendly name applied as the
	// ambientweather.station.name resource attribute.
	StationName string `mapstructure:"station_name"`

	// prevent unkeyed literal initialization
	_ struct{}
}

func (cfg *Config) Validate() error {
	if cfg.NetAddr.Endpoint == "" {
		return errors.New("endpoint cannot be empty")
	}
	if !strings.HasPrefix(cfg.Path, "/") {
		return fmt.Errorf("path %q must begin with '/'", cfg.Path)
	}
	return nil
}
