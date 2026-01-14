// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logging

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestMetricWithoutDataPoints_Sum(t *testing.T) {
	metric := pmetric.NewMetric()
	metric.SetName("test.sum")
	metric.SetDescription("test sum metric")
	metric.SetUnit("ms")
	sum := metric.SetEmptySum()
	sum.SetIsMonotonic(true)
	sum.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)

	// Add datapoints - these should NOT be marshaled
	dp1 := sum.DataPoints().AppendEmpty()
	dp1.SetIntValue(100)
	dp2 := sum.DataPoints().AppendEmpty()
	dp2.SetIntValue(200)

	// Create observer to capture log output
	core, logs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)

	// Marshal the metric without datapoints
	logger.Info("test", zap.Object("metric", MetricWithoutDataPoints(metric)))

	// Verify the log was created
	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]

	// Verify metric fields are present
	metricObj := entry.ContextMap()["metric"].(map[string]interface{})
	require.Equal(t, "test.sum", metricObj["name"])
	require.Equal(t, "test sum metric", metricObj["description"])
	require.Equal(t, "ms", metricObj["unit"])
	require.Equal(t, "Sum", metricObj["type"])
	require.Equal(t, "Cumulative", metricObj["aggregation_temporality"])
	require.Equal(t, true, metricObj["is_monotonic"])

	// Verify datapoints array is NOT present
	_, hasDatapoints := metricObj["datapoints"]
	require.False(t, hasDatapoints, "MetricWithoutDataPoints should not include datapoints array")
}

func TestMetricWithoutDataPoints_Gauge(t *testing.T) {
	metric := pmetric.NewMetric()
	metric.SetName("test.gauge")
	gauge := metric.SetEmptyGauge()

	// Add datapoints - these should NOT be marshaled
	dp := gauge.DataPoints().AppendEmpty()
	dp.SetDoubleValue(42.5)

	// Create observer to capture log output
	core, logs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)

	logger.Info("test", zap.Object("metric", MetricWithoutDataPoints(metric)))

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]

	metricObj := entry.ContextMap()["metric"].(map[string]interface{})
	require.Equal(t, "test.gauge", metricObj["name"])
	require.Equal(t, "Gauge", metricObj["type"])

	// Verify datapoints array is NOT present
	_, hasDatapoints := metricObj["datapoints"]
	require.False(t, hasDatapoints, "MetricWithoutDataPoints should not include datapoints array")
}

func TestMetricWithoutDataPoints_Histogram(t *testing.T) {
	metric := pmetric.NewMetric()
	metric.SetName("test.histogram")
	histogram := metric.SetEmptyHistogram()
	histogram.SetAggregationTemporality(pmetric.AggregationTemporalityDelta)

	// Add datapoints - these should NOT be marshaled
	dp := histogram.DataPoints().AppendEmpty()
	dp.SetCount(10)

	// Create observer to capture log output
	core, logs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)

	logger.Info("test", zap.Object("metric", MetricWithoutDataPoints(metric)))

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]

	metricObj := entry.ContextMap()["metric"].(map[string]interface{})
	require.Equal(t, "test.histogram", metricObj["name"])
	require.Equal(t, "Histogram", metricObj["type"])
	require.Equal(t, "Delta", metricObj["aggregation_temporality"])

	// Verify datapoints array is NOT present
	_, hasDatapoints := metricObj["datapoints"]
	require.False(t, hasDatapoints, "MetricWithoutDataPoints should not include datapoints array")
}

func TestMetricWithoutDataPoints_ExponentialHistogram(t *testing.T) {
	metric := pmetric.NewMetric()
	metric.SetName("test.exp_histogram")
	expHistogram := metric.SetEmptyExponentialHistogram()
	expHistogram.SetAggregationTemporality(pmetric.AggregationTemporalityCumulative)

	// Add datapoints - these should NOT be marshaled
	dp := expHistogram.DataPoints().AppendEmpty()
	dp.SetCount(5)

	// Create observer to capture log output
	core, logs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)

	logger.Info("test", zap.Object("metric", MetricWithoutDataPoints(metric)))

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]

	metricObj := entry.ContextMap()["metric"].(map[string]interface{})
	require.Equal(t, "test.exp_histogram", metricObj["name"])
	require.Equal(t, "ExponentialHistogram", metricObj["type"])
	require.Equal(t, "Cumulative", metricObj["aggregation_temporality"])

	// Verify datapoints array is NOT present
	_, hasDatapoints := metricObj["datapoints"]
	require.False(t, hasDatapoints, "MetricWithoutDataPoints should not include datapoints array")
}

func TestMetricWithoutDataPoints_Summary(t *testing.T) {
	metric := pmetric.NewMetric()
	metric.SetName("test.summary")
	summary := metric.SetEmptySummary()

	// Add datapoints - these should NOT be marshaled
	dp := summary.DataPoints().AppendEmpty()
	dp.SetCount(100)

	// Create observer to capture log output
	core, logs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)

	logger.Info("test", zap.Object("metric", MetricWithoutDataPoints(metric)))

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]

	metricObj := entry.ContextMap()["metric"].(map[string]interface{})
	require.Equal(t, "test.summary", metricObj["name"])
	require.Equal(t, "Summary", metricObj["type"])

	// Verify datapoints array is NOT present
	_, hasDatapoints := metricObj["datapoints"]
	require.False(t, hasDatapoints, "MetricWithoutDataPoints should not include datapoints array")
}

func TestMetric_IncludesDataPoints(t *testing.T) {
	// Verify that the regular Metric marshaler still includes datapoints
	metric := pmetric.NewMetric()
	metric.SetName("test.sum")
	sum := metric.SetEmptySum()

	dp1 := sum.DataPoints().AppendEmpty()
	dp1.SetIntValue(100)

	// Create observer to capture log output
	core, logs := observer.New(zapcore.DebugLevel)
	logger := zap.New(core)

	logger.Info("test", zap.Object("metric", Metric(metric)))

	require.Equal(t, 1, logs.Len())
	entry := logs.All()[0]

	metricObj := entry.ContextMap()["metric"].(map[string]interface{})

	// Verify datapoints array IS present for regular Metric
	datapoints, hasDatapoints := metricObj["datapoints"]
	require.True(t, hasDatapoints, "Regular Metric should include datapoints array")
	require.NotNil(t, datapoints)
}
