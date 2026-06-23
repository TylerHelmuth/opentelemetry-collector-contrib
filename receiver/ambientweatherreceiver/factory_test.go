// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ambientweatherreceiver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver/internal/metadata"
)

func TestNewFactory(t *testing.T) {
	factory := NewFactory()
	require.NotNil(t, factory)
	assert.Equal(t, "ambientweather", factory.Type().String())
}

func TestCreateDefaultConfig(t *testing.T) {
	cfg := NewFactory().CreateDefaultConfig()
	require.NotNil(t, cfg)

	config, ok := cfg.(*Config)
	require.True(t, ok)
	assert.Equal(t, defaultEndpoint, config.NetAddr.Endpoint)
	assert.Equal(t, defaultPath, config.Path)
	assert.Empty(t, config.StationName)
	assert.Equal(t, metadata.NewDefaultMetricsBuilderConfig(), config.MetricsBuilderConfig)
}

func TestCreateMetricsReceiver(t *testing.T) {
	factory := NewFactory()
	r, err := factory.CreateMetrics(
		t.Context(),
		receivertest.NewNopSettings(metadata.Type),
		factory.CreateDefaultConfig(),
		consumertest.NewNop())
	require.NoError(t, err)
	assert.NotNil(t, r)
}

func TestCreateMetricsReceiverInvalidConfig(t *testing.T) {
	_, err := NewFactory().CreateMetrics(
		t.Context(),
		receivertest.NewNopSettings(metadata.Type),
		&struct{}{},
		consumertest.NewNop())
	assert.ErrorIs(t, err, errInvalidConfig)
}

func TestFactoryConfigStruct(t *testing.T) {
	require.NoError(t, componenttest.CheckConfigStruct(NewFactory().CreateDefaultConfig()))
}
