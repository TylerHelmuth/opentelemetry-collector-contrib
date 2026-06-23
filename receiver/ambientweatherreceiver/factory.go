// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ambientweatherreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver"

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver/internal/metadata"
)

const (
	defaultEndpoint = "localhost:6255"
	defaultPath     = "/data/report/"
)

var errInvalidConfig = errors.New("config was not an Ambient Weather receiver config")

// NewFactory creates a factory for the Ambient Weather receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, metadata.MetricsStability))
}

func createDefaultConfig() component.Config {
	netAddr := confignet.NewDefaultAddrConfig()
	netAddr.Transport = confignet.TransportTypeTCP
	netAddr.Endpoint = defaultEndpoint
	return &Config{
		ServerConfig:         confighttp.ServerConfig{NetAddr: netAddr},
		MetricsBuilderConfig: metadata.NewDefaultMetricsBuilderConfig(),
		Path:                 defaultPath,
	}
}

func createMetricsReceiver(
	_ context.Context,
	settings receiver.Settings,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (receiver.Metrics, error) {
	rCfg, ok := cfg.(*Config)
	if !ok {
		return nil, errInvalidConfig
	}
	return newAmbientWeatherReceiver(rCfg, settings, nextConsumer), nil
}
