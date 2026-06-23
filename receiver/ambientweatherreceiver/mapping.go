// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ambientweatherreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver"

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver/internal/metadata"
)

// dateUTCLayout is the timestamp format used by the AMBWeather custom-server protocol.
const dateUTCLayout = "2006-01-02 15:04:05"

// maxSensorChannels is the highest add-on remote-sensor channel the AMBWeather protocol exposes.
const maxSensorChannels = 10

// Unit-conversion factors from the device's imperial units to canonical SI/UCUM units.
const (
	mphToMetersPerSecond = 0.44704
	inHgToHectopascal    = 33.86389
	inchToMillimeter     = 25.4
)

func fahrenheitToCelsius(f float64) float64 { return (f - 32) * 5 / 9 }
func mphToMS(v float64) float64             { return v * mphToMetersPerSecond }
func inHgToHPa(v float64) float64           { return v * inHgToHectopascal }
func inToMM(v float64) float64              { return v * inchToMillimeter }

// rainfallWindows maps each AMBWeather rainfall accumulation field to its rainfall.period value.
var rainfallWindows = []struct {
	field  string
	period string
}{
	{"hourlyrainin", "hourly"},
	{"dailyrainin", "daily"},
	{"24hourrainin", "24hour"},
	{"weeklyrainin", "weekly"},
	{"monthlyrainin", "monthly"},
	{"yearlyrainin", "yearly"},
	{"eventrainin", "event"},
	{"totalrainin", "total"},
}

// buildMetrics converts a set of AMBWeather upload fields into a pmetric.Metrics.
// A fresh MetricsBuilder is created per call so the method is safe for concurrent requests.
func (r *ambientWeatherReceiver) buildMetrics(fields url.Values) pmetric.Metrics {
	mb := metadata.NewMetricsBuilder(r.cfg.MetricsBuilderConfig, r.settings)
	ts := parseTimestamp(fields, r.logger)

	// Temperature (degrees F -> Celsius).
	r.recordFloat(fields, "tempf", fahrenheitToCelsius, func(v float64) {
		mb.RecordAmbientweatherTemperatureDataPoint(ts, v, "outdoor")
	})
	r.recordFloat(fields, "tempinf", fahrenheitToCelsius, func(v float64) {
		mb.RecordAmbientweatherTemperatureDataPoint(ts, v, "indoor")
	})

	// Humidity (percent, passed through).
	r.recordFloat(fields, "humidity", nil, func(v float64) {
		mb.RecordAmbientweatherHumidityDataPoint(ts, v, "outdoor")
	})
	r.recordFloat(fields, "humidityin", nil, func(v float64) {
		mb.RecordAmbientweatherHumidityDataPoint(ts, v, "indoor")
	})

	// PM2.5 (micrograms per cubic meter, passed through).
	r.recordFloat(fields, "pm25", nil, func(v float64) {
		mb.RecordAmbientweatherPm25DataPoint(ts, v, "outdoor")
	})
	r.recordFloat(fields, "pm25_in", nil, func(v float64) {
		mb.RecordAmbientweatherPm25DataPoint(ts, v, "indoor")
	})

	// Battery status (0 = low, 1 = OK).
	r.recordInt(fields, "battout", func(v int64) {
		mb.RecordAmbientweatherBatteryStatusDataPoint(ts, v, "outdoor")
	})
	r.recordInt(fields, "battin", func(v int64) {
		mb.RecordAmbientweatherBatteryStatusDataPoint(ts, v, "indoor")
	})
	r.recordInt(fields, "batt_lightning", func(v int64) {
		mb.RecordAmbientweatherBatteryStatusDataPoint(ts, v, "lightning")
	})

	// Multi-channel add-on sensors.
	for ch := 1; ch <= maxSensorChannels; ch++ {
		sensor := strconv.Itoa(ch)
		r.recordFloat(fields, fmt.Sprintf("temp%df", ch), fahrenheitToCelsius, func(v float64) {
			mb.RecordAmbientweatherTemperatureDataPoint(ts, v, sensor)
		})
		r.recordFloat(fields, fmt.Sprintf("humidity%d", ch), nil, func(v float64) {
			mb.RecordAmbientweatherHumidityDataPoint(ts, v, sensor)
		})
		r.recordFloat(fields, fmt.Sprintf("soilhum%d", ch), nil, func(v float64) {
			mb.RecordAmbientweatherSoilMoistureDataPoint(ts, v, sensor)
		})
		r.recordInt(fields, fmt.Sprintf("batt%d", ch), func(v int64) {
			mb.RecordAmbientweatherBatteryStatusDataPoint(ts, v, sensor)
		})
	}

	// Wind (mph -> m/s; direction in degrees).
	r.recordFloat(fields, "windspeedmph", mphToMS, func(v float64) {
		mb.RecordAmbientweatherWindSpeedDataPoint(ts, v)
	})
	r.recordFloat(fields, "windgustmph", mphToMS, func(v float64) {
		mb.RecordAmbientweatherWindGustDataPoint(ts, v)
	})
	r.recordFloat(fields, "winddir", nil, func(v float64) {
		mb.RecordAmbientweatherWindDirectionDataPoint(ts, v)
	})

	// Barometric pressure (inHg -> hPa).
	r.recordFloat(fields, "baromrelin", inHgToHPa, func(v float64) {
		mb.RecordAmbientweatherPressureRelativeDataPoint(ts, v)
	})
	r.recordFloat(fields, "baromabsin", inHgToHPa, func(v float64) {
		mb.RecordAmbientweatherPressureAbsoluteDataPoint(ts, v)
	})

	// Rainfall (inches -> mm) across each accumulation window, plus the instantaneous rate.
	for _, w := range rainfallWindows {
		period := w.period
		r.recordFloat(fields, w.field, inToMM, func(v float64) {
			mb.RecordAmbientweatherRainfallDataPoint(ts, v, period)
		})
	}
	r.recordFloat(fields, "rainratein", inToMM, func(v float64) {
		mb.RecordAmbientweatherRainRateDataPoint(ts, v)
	})

	// Solar / UV.
	r.recordFloat(fields, "solarradiation", nil, func(v float64) {
		mb.RecordAmbientweatherSolarRadiationDataPoint(ts, v)
	})
	r.recordInt(fields, "uv", func(v int64) {
		mb.RecordAmbientweatherUvIndexDataPoint(ts, v)
	})

	// Lightning. lightning_distance is already in km on consoles configured with metric
	// wind-speed units (km/h or m/s); no conversion is applied.
	r.recordFloat(fields, "lightning_distance", nil, func(v float64) {
		mb.RecordAmbientweatherLightningDistanceDataPoint(ts, v)
	})
	r.recordInt(fields, "lightning_day", func(v int64) {
		mb.RecordAmbientweatherLightningStrikesDataPoint(ts, v, "daily")
	})
	r.recordInt(fields, "lightning_hour", func(v int64) {
		mb.RecordAmbientweatherLightningStrikesDataPoint(ts, v, "hourly")
	})

	rb := mb.NewResourceBuilder()
	if mac := stationMAC(fields); mac != "" {
		rb.SetAmbientweatherStationMac(mac)
	}
	if stationType := fields.Get("stationtype"); stationType != "" {
		rb.SetAmbientweatherStationType(stationType)
	}
	if r.cfg.StationName != "" {
		rb.SetAmbientweatherStationName(r.cfg.StationName)
	}

	return mb.Emit(metadata.WithResource(rb.Emit()))
}

// recordFloat parses a floating-point field, optionally converts it, and records it.
// Missing or unparseable fields are skipped (logged at debug).
func (r *ambientWeatherReceiver) recordFloat(fields url.Values, key string, convert func(float64) float64, record func(float64)) {
	raw := fields.Get(key)
	if raw == "" {
		return
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		r.logger.Debug("could not parse numeric field",
			zap.String("field", key), zap.String("value", raw), zap.Error(err))
		return
	}
	if convert != nil {
		v = convert(v)
	}
	record(v)
}

// recordInt parses an integer-valued field (accepting "1" or "1.0") and records it.
func (r *ambientWeatherReceiver) recordInt(fields url.Values, key string, record func(int64)) {
	raw := fields.Get(key)
	if raw == "" {
		return
	}
	v, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		r.logger.Debug("could not parse integer field",
			zap.String("field", key), zap.String("value", raw), zap.Error(err))
		return
	}
	record(int64(v))
}

// parseTimestamp resolves the point timestamp from the console's dateutc field,
// falling back to the receive time when it is absent, "now", or unparseable.
func parseTimestamp(fields url.Values, logger *zap.Logger) pcommon.Timestamp {
	raw := fields.Get("dateutc")
	if raw != "" && !strings.EqualFold(raw, "now") {
		if t, err := time.Parse(dateUTCLayout, raw); err == nil {
			return pcommon.NewTimestampFromTime(t.UTC())
		}
		logger.Debug("could not parse dateutc; using receive time", zap.String("dateutc", raw))
	}
	return pcommon.NewTimestampFromTime(time.Now())
}

// stationMAC returns the station identity, preferring the AMBWeather MAC over the Ecowitt PASSKEY.
func stationMAC(fields url.Values) string {
	if mac := fields.Get("MAC"); mac != "" {
		return mac
	}
	return fields.Get("PASSKEY")
}
