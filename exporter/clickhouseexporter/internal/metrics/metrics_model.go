// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metrics // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/clickhouseexporter/internal/metrics"

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/column"
	"github.com/ClickHouse/clickhouse-go/v2/lib/column/orderedmap"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	conventions "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/clickhouseexporter/internal/sqltemplates"
)

var supportedMetricTypes = map[pmetric.MetricType]string{
	pmetric.MetricTypeGauge:                sqltemplates.MetricsGaugeCreateTable,
	pmetric.MetricTypeSum:                  sqltemplates.MetricsSumCreateTable,
	pmetric.MetricTypeHistogram:            sqltemplates.MetricsHistogramCreateTable,
	pmetric.MetricTypeExponentialHistogram: sqltemplates.MetricsExpHistogramCreateTable,
	pmetric.MetricTypeSummary:              sqltemplates.MetricsSummaryCreateTable,
}

var logger *zap.Logger

type MetricTablesConfigMapper map[pmetric.MetricType]MetricTypeConfig

type MetricTypeConfig struct {
	Name string `mapstructure:"name"`
}

// MetricsModel is used to group metric data and insert into clickhouse
// any type of metrics need implement it.
type MetricsModel interface {
	// Add used to bind MetricsMetaData to a specific metric then put them into a slice
	Add(resAttr pcommon.Map, resURL string, scopeInstr pcommon.InstrumentationScope, scopeURL string, metrics any, name, description, unit string) error
	// insert is used to insert metric data to clickhouse
	insert(ctx context.Context, db driver.Conn) error
}

// MetricsMetaData contain specific metric data
type MetricsMetaData struct {
	ResAttr    pcommon.Map
	ResURL     string
	ScopeURL   string
	ScopeInstr pcommon.InstrumentationScope
}

// SetLogger set a logger instance
func SetLogger(l *zap.Logger) {
	logger = l
}

// NewMetricsTable create metric tables with an expiry time to storage metric telemetry data
func NewMetricsTable(ctx context.Context, tablesConfig MetricTablesConfigMapper, database, cluster, engine, ttlExpr string, db driver.Conn) error {
	for key, ddlTemplate := range supportedMetricTypes {
		query := fmt.Sprintf(ddlTemplate, database, tablesConfig[key].Name, cluster, engine, ttlExpr)
		if err := db.Exec(ctx, query); err != nil {
			return fmt.Errorf("exec create metrics table sql: %w", err)
		}
	}
	return nil
}

// NewMetricsModel create a model for contain different metric data
func NewMetricsModel(tablesConfig MetricTablesConfigMapper, database string) map[pmetric.MetricType]MetricsModel {
	return map[pmetric.MetricType]MetricsModel{
		pmetric.MetricTypeGauge: &gaugeMetrics{
			insertSQL: fmt.Sprintf(sqltemplates.MetricsGaugeInsert, database, tablesConfig[pmetric.MetricTypeGauge].Name),
		},
		pmetric.MetricTypeSum: &sumMetrics{
			insertSQL: fmt.Sprintf(sqltemplates.MetricsSumInsert, database, tablesConfig[pmetric.MetricTypeSum].Name),
		},
		pmetric.MetricTypeHistogram: &histogramMetrics{
			insertSQL: fmt.Sprintf(sqltemplates.MetricsHistogramInsert, database, tablesConfig[pmetric.MetricTypeHistogram].Name),
		},
		pmetric.MetricTypeExponentialHistogram: &expHistogramMetrics{
			insertSQL: fmt.Sprintf(sqltemplates.MetricsExpHistogramInsert, database, tablesConfig[pmetric.MetricTypeExponentialHistogram].Name),
		},
		pmetric.MetricTypeSummary: &summaryMetrics{
			insertSQL: fmt.Sprintf(sqltemplates.MetricsSummaryInsert, database, tablesConfig[pmetric.MetricTypeSummary].Name),
		},
	}
}

// InsertMetrics insert metric data into clickhouse concurrently
func InsertMetrics(ctx context.Context, db driver.Conn, metricsMap map[pmetric.MetricType]MetricsModel) error {
	errsChan := make(chan error, len(supportedMetricTypes))
	wg := &sync.WaitGroup{}
	for _, m := range metricsMap {
		wg.Add(1)
		go func(m MetricsModel, wg *sync.WaitGroup) {
			errsChan <- m.insert(ctx, db)
			wg.Done()
		}(m, wg)
	}
	wg.Wait()
	close(errsChan)
	var errs error
	for err := range errsChan {
		errs = errors.Join(errs, err)
	}
	return errs
}

func convertExemplars(exemplars pmetric.ExemplarSlice) (clickhouse.ArraySet, clickhouse.ArraySet, clickhouse.ArraySet, clickhouse.ArraySet, clickhouse.ArraySet) {
	var (
		attrs    clickhouse.ArraySet
		times    clickhouse.ArraySet
		values   clickhouse.ArraySet
		traceIDs clickhouse.ArraySet
		spanIDs  clickhouse.ArraySet
	)
	for i := 0; i < exemplars.Len(); i++ {
		exemplar := exemplars.At(i)
		attrs = append(attrs, AttributesToMap(exemplar.FilteredAttributes()))
		times = append(times, exemplar.Timestamp().AsTime())
		values = append(values, getValue(exemplar.IntValue(), exemplar.DoubleValue(), exemplar.ValueType()))

		traceID, spanID := exemplar.TraceID(), exemplar.SpanID()
		traceIDs = append(traceIDs, hex.EncodeToString(traceID[:]))
		spanIDs = append(spanIDs, hex.EncodeToString(spanID[:]))
	}
	return attrs, times, values, traceIDs, spanIDs
}

// https://github.com/open-telemetry/opentelemetry-proto/blob/main/opentelemetry/proto/metrics/v1/metrics.proto#L358
// define two types for one datapoint value, clickhouse only use one value of float64 to store them
func getValue(intValue int64, floatValue float64, dataType any) float64 {
	switch t := dataType.(type) {
	case pmetric.ExemplarValueType:
		switch t {
		case pmetric.ExemplarValueTypeDouble:
			return floatValue
		case pmetric.ExemplarValueTypeInt:
			return float64(intValue)
		case pmetric.ExemplarValueTypeEmpty:
			return 0.0
		default:
			logger.Warn("Can't find a suitable value for ExemplarValueType, use 0.0 as default")
			return 0.0
		}
	case pmetric.NumberDataPointValueType:
		switch t {
		case pmetric.NumberDataPointValueTypeDouble:
			return floatValue
		case pmetric.NumberDataPointValueTypeInt:
			return float64(intValue)
		case pmetric.NumberDataPointValueTypeEmpty:
			return 0.0
		default:
			logger.Warn("Can't find a suitable value for NumberDataPointValueType, use 0.0 as default")
			return 0.0
		}
	default:
		logger.Warn("unsupported ValueType, current support: ExemplarValueType, NumberDataPointValueType, ues 0.0 as default")
		return 0.0
	}
}

func AttributesToMap(attributes pcommon.Map) column.IterableOrderedMap {
	return orderedmap.CollectN(func(yield func(string, string) bool) {
		for k, v := range attributes.All() {
			yield(k, v.AsString())
		}
	}, attributes.Len())
}

func GetServiceName(resAttr pcommon.Map) string {
	if v, ok := resAttr.Get(string(conventions.ServiceNameKey)); ok {
		return v.AsString()
	}

	return ""
}

func convertSliceToArraySet[T any](slice []T) clickhouse.ArraySet {
	var set clickhouse.ArraySet
	for _, item := range slice {
		set = append(set, item)
	}
	return set
}

func convertValueAtQuantile(valueAtQuantile pmetric.SummaryDataPointValueAtQuantileSlice) (clickhouse.ArraySet, clickhouse.ArraySet) {
	var (
		quantiles clickhouse.ArraySet
		values    clickhouse.ArraySet
	)
	for i := 0; i < valueAtQuantile.Len(); i++ {
		value := valueAtQuantile.At(i)
		quantiles = append(quantiles, value.Quantile())
		values = append(values, value.Value())
	}
	return quantiles, values
}

func newPlaceholder(count int) *string {
	var b strings.Builder
	for i := 0; i < count; i++ {
		b.WriteString(",?")
	}
	b.WriteString("),")
	placeholder := strings.Replace(b.String(), ",", "(", 1)
	return &placeholder
}
