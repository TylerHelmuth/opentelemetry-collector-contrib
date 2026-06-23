// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ambientweatherreceiver

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver/internal/metadata"
)

func newMappingReceiver(t *testing.T, customize func(*Config)) *ambientWeatherReceiver {
	t.Helper()
	cfg := createDefaultConfig().(*Config)
	if customize != nil {
		customize(cfg)
	}
	return newAmbientWeatherReceiver(cfg, receivertest.NewNopSettings(metadata.Type), consumertest.NewNop())
}

func metricsByName(md pmetric.Metrics) map[string]pmetric.Metric {
	out := make(map[string]pmetric.Metric)
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		sms := rms.At(i).ScopeMetrics()
		for j := 0; j < sms.Len(); j++ {
			ms := sms.At(j).Metrics()
			for k := 0; k < ms.Len(); k++ {
				out[ms.At(k).Name()] = ms.At(k)
			}
		}
	}
	return out
}

func doubleGauge(t *testing.T, m pmetric.Metric) float64 {
	t.Helper()
	require.Equal(t, 1, m.Gauge().DataPoints().Len())
	return m.Gauge().DataPoints().At(0).DoubleValue()
}

func intGauge(t *testing.T, m pmetric.Metric) int64 {
	t.Helper()
	require.Equal(t, 1, m.Gauge().DataPoints().Len())
	return m.Gauge().DataPoints().At(0).IntValue()
}

func dataPointByAttr(t *testing.T, m pmetric.Metric, key, value string) pmetric.NumberDataPoint {
	t.Helper()
	dps := m.Gauge().DataPoints()
	for i := 0; i < dps.Len(); i++ {
		if v, ok := dps.At(i).Attributes().Get(key); ok && v.Str() == value {
			return dps.At(i)
		}
	}
	t.Fatalf("no %s data point with %s=%s", m.Name(), key, value)
	return pmetric.NumberDataPoint{}
}

func assertResourceAttr(t *testing.T, res pcommon.Resource, key, want string) {
	t.Helper()
	v, ok := res.Attributes().Get(key)
	require.True(t, ok, "missing resource attribute %s", key)
	assert.Equal(t, want, v.Str())
}

func TestUnitConversions(t *testing.T) {
	assert.InDelta(t, 0.0, fahrenheitToCelsius(32), 1e-9)
	assert.InDelta(t, 100.0, fahrenheitToCelsius(212), 1e-9)
	assert.InDelta(t, 0.44704, mphToMS(1), 1e-9)
	assert.InDelta(t, 25.4, inToMM(1), 1e-9)
	assert.InDelta(t, 33.86389, inHgToHPa(1), 1e-9)
}

func TestParseTimestamp(t *testing.T) {
	logger := zap.NewNop()
	want := time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC)

	got := parseTimestamp(url.Values{"dateutc": {"2024-01-15 10:30:00"}}, logger)
	assert.Equal(t, pcommon.NewTimestampFromTime(want), got)

	for _, raw := range []string{"now", "", "garbage"} {
		vals := url.Values{}
		if raw != "" {
			vals.Set("dateutc", raw)
		}
		ts := parseTimestamp(vals, logger)
		assert.WithinDuration(t, time.Now(), ts.AsTime(), 5*time.Second)
	}
}

func TestStationMAC(t *testing.T) {
	assert.Equal(t, "AA:BB", stationMAC(url.Values{"MAC": {"AA:BB"}, "PASSKEY": {"HEX"}}))
	assert.Equal(t, "HEX", stationMAC(url.Values{"PASSKEY": {"HEX"}}))
	assert.Empty(t, stationMAC(url.Values{}))
}

func TestBuildMetricsCoreFields(t *testing.T) {
	r := newMappingReceiver(t, func(c *Config) { c.StationName = "backyard" })

	fields := url.Values{}
	fields.Set("PASSKEY", "ABC123")
	fields.Set("stationtype", "AMBWeatherV4.3.4")
	fields.Set("dateutc", "2024-01-15 10:30:00")
	fields.Set("tempf", "68.0")
	fields.Set("tempinf", "71.6")
	fields.Set("humidity", "55")
	fields.Set("windspeedmph", "10")
	fields.Set("windgustmph", "20")
	fields.Set("winddir", "180")
	fields.Set("baromrelin", "29.92")
	fields.Set("baromabsin", "29.5")
	fields.Set("dailyrainin", "0.5")
	fields.Set("rainratein", "0.1")
	fields.Set("battout", "1")

	md := r.buildMetrics(fields)

	require.Equal(t, 1, md.ResourceMetrics().Len())
	res := md.ResourceMetrics().At(0).Resource()
	assertResourceAttr(t, res, "ambientweather.station.mac", "ABC123")
	assertResourceAttr(t, res, "ambientweather.station.type", "AMBWeatherV4.3.4")
	assertResourceAttr(t, res, "ambientweather.station.name", "backyard")

	byName := metricsByName(md)

	assert.InDelta(t, 20.0, dataPointByAttr(t, byName["ambientweather.temperature"], "sensor", "outdoor").DoubleValue(), 1e-6)
	assert.InDelta(t, 22.0, dataPointByAttr(t, byName["ambientweather.temperature"], "sensor", "indoor").DoubleValue(), 1e-6)
	assert.InDelta(t, 55.0, dataPointByAttr(t, byName["ambientweather.humidity"], "sensor", "outdoor").DoubleValue(), 1e-6)

	assert.InDelta(t, 4.4704, doubleGauge(t, byName["ambientweather.wind.speed"]), 1e-6)
	assert.InDelta(t, 8.9408, doubleGauge(t, byName["ambientweather.wind.gust"]), 1e-6)
	assert.InDelta(t, 180.0, doubleGauge(t, byName["ambientweather.wind.direction"]), 1e-6)

	assert.InDelta(t, 29.92*inHgToHectopascal, doubleGauge(t, byName["ambientweather.pressure.relative"]), 1e-6)
	assert.InDelta(t, 29.5*inHgToHectopascal, doubleGauge(t, byName["ambientweather.pressure.absolute"]), 1e-6)

	assert.InDelta(t, 12.7, dataPointByAttr(t, byName["ambientweather.rainfall"], "rainfall.period", "daily").DoubleValue(), 1e-6)
	assert.InDelta(t, 2.54, doubleGauge(t, byName["ambientweather.rain.rate"]), 1e-6)

	assert.Equal(t, int64(1), dataPointByAttr(t, byName["ambientweather.battery.status"], "sensor", "outdoor").IntValue())

	want := pcommon.NewTimestampFromTime(time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC))
	assert.Equal(t, want, byName["ambientweather.temperature"].Gauge().DataPoints().At(0).Timestamp())
}

func TestBuildMetricsChannelsAndOptional(t *testing.T) {
	r := newMappingReceiver(t, func(c *Config) {
		c.Metrics.AmbientweatherPm25.Enabled = true
		c.Metrics.AmbientweatherSoilMoisture.Enabled = true
		c.Metrics.AmbientweatherSolarRadiation.Enabled = true
		c.Metrics.AmbientweatherUvIndex.Enabled = true
	})

	fields := url.Values{}
	fields.Set("MAC", "AA:BB:CC")
	fields.Set("temp1f", "50")
	fields.Set("humidity1", "40")
	fields.Set("soilhum1", "30")
	fields.Set("batt1", "0")
	fields.Set("pm25", "12.5")
	fields.Set("solarradiation", "800")
	fields.Set("uv", "5")

	byName := metricsByName(r.buildMetrics(fields))

	assert.InDelta(t, 10.0, dataPointByAttr(t, byName["ambientweather.temperature"], "sensor", "1").DoubleValue(), 1e-6)
	assert.InDelta(t, 40.0, dataPointByAttr(t, byName["ambientweather.humidity"], "sensor", "1").DoubleValue(), 1e-6)
	assert.InDelta(t, 30.0, dataPointByAttr(t, byName["ambientweather.soil.moisture"], "sensor", "1").DoubleValue(), 1e-6)
	assert.Equal(t, int64(0), dataPointByAttr(t, byName["ambientweather.battery.status"], "sensor", "1").IntValue())
	assert.InDelta(t, 12.5, dataPointByAttr(t, byName["ambientweather.pm25"], "sensor", "outdoor").DoubleValue(), 1e-6)
	assert.InDelta(t, 800.0, doubleGauge(t, byName["ambientweather.solar.radiation"]), 1e-6)
	assert.Equal(t, int64(5), intGauge(t, byName["ambientweather.uv.index"]))
}

func TestBuildMetricsLightning(t *testing.T) {
	r := newMappingReceiver(t, nil)

	fields := url.Values{}
	fields.Set("MAC", "AA:BB")
	fields.Set("lightning_distance", "12")
	fields.Set("lightning_day", "5")
	fields.Set("lightning_hour", "2")
	fields.Set("batt_lightning", "1")

	byName := metricsByName(r.buildMetrics(fields))

	// Distance is taken as km (no conversion).
	assert.InDelta(t, 12.0, doubleGauge(t, byName["ambientweather.lightning.distance"]), 1e-6)
	assert.Equal(t, int64(5), dataPointByAttr(t, byName["ambientweather.lightning.strikes"], "lightning.period", "daily").IntValue())
	assert.Equal(t, int64(2), dataPointByAttr(t, byName["ambientweather.lightning.strikes"], "lightning.period", "hourly").IntValue())
	assert.Equal(t, int64(1), dataPointByAttr(t, byName["ambientweather.battery.status"], "sensor", "lightning").IntValue())
}

func TestBuildMetricsSkipsUnparseableFields(t *testing.T) {
	r := newMappingReceiver(t, nil)

	fields := url.Values{}
	fields.Set("MAC", "AA:BB")
	fields.Set("tempf", "not-a-number")
	fields.Set("humidity", "60")

	byName := metricsByName(r.buildMetrics(fields))

	_, hasTemp := byName["ambientweather.temperature"]
	assert.False(t, hasTemp, "unparseable tempf should be skipped")
	assert.InDelta(t, 60.0, dataPointByAttr(t, byName["ambientweather.humidity"], "sensor", "outdoor").DoubleValue(), 1e-6)
}
