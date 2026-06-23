// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ambientweatherreceiver

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
	"go.opentelemetry.io/collector/receiver/receivertest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver/internal/metadata"
)

func newHandlerReceiver(t *testing.T, next consumer.Metrics) *ambientWeatherReceiver {
	t.Helper()
	set := receivertest.NewNopSettings(metadata.Type)
	r := newAmbientWeatherReceiver(createDefaultConfig().(*Config), set, next)
	obs, err := receiverhelper.NewObsReport(receiverhelper.ObsReportSettings{
		ReceiverID:             set.ID,
		Transport:              "http",
		ReceiverCreateSettings: set,
	})
	require.NoError(t, err)
	r.obsrecv = obs
	return r
}

func sampleQuery() url.Values {
	q := url.Values{}
	q.Set("MAC", "AA:BB:CC:DD")
	q.Set("stationtype", "AMBWeatherV4.3.4")
	q.Set("tempf", "68")
	q.Set("humidity", "55")
	return q
}

func TestHandleReportGET(t *testing.T) {
	sink := new(consumertest.MetricsSink)
	r := newHandlerReceiver(t, sink)

	req := httptest.NewRequest(http.MethodGet, "/data/report/?"+sampleQuery().Encode(), http.NoBody)
	rec := httptest.NewRecorder()
	r.handleReport(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
	require.Len(t, sink.AllMetrics(), 1)
	assert.Positive(t, sink.AllMetrics()[0].DataPointCount())
}

// Consoles routinely POST the readings as a urlencoded body without a Content-Type
// header; the body must still be parsed.
func TestHandleReportPOSTBodyWithoutContentType(t *testing.T) {
	sink := new(consumertest.MetricsSink)
	r := newHandlerReceiver(t, sink)

	req := httptest.NewRequest(http.MethodPost, "/data/report/", strings.NewReader(sampleQuery().Encode()))
	rec := httptest.NewRecorder()
	r.handleReport(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, sink.AllMetrics(), 1)
	assert.Positive(t, sink.AllMetrics()[0].DataPointCount())
}

// Guards the deliberate query+body merge: a console may send some fields in the URL and
// others in a urlencoded POST body. Both must be recorded, and parsing must not depend on
// http.Request.ParseForm (which would only succeed here because of the Content-Type).
func TestHandleReportMergesQueryAndBody(t *testing.T) {
	sink := new(consumertest.MetricsSink)
	r := newHandlerReceiver(t, sink)

	body := url.Values{"tempf": {"68"}}.Encode()
	req := httptest.NewRequest(http.MethodPost, "/data/report/?MAC=AA:BB&humidity=55", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	r.handleReport(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, sink.AllMetrics(), 1)

	byName := metricsByName(sink.AllMetrics()[0])
	require.Contains(t, byName, "ambientweather.temperature", "tempf from POST body should be recorded")
	require.Contains(t, byName, "ambientweather.humidity", "humidity from query string should be recorded")
	assert.InDelta(t, 20.0, dataPointByAttr(t, byName["ambientweather.temperature"], "sensor", "outdoor").DoubleValue(), 1e-6)
	assert.InDelta(t, 55.0, dataPointByAttr(t, byName["ambientweather.humidity"], "sensor", "outdoor").DoubleValue(), 1e-6)
}

// Real Ambient Weather consoles concatenate the parameters directly onto the configured
// path with no "?", so they arrive in the URL path rather than the query string. The
// receiver must still recover and record them.
func TestHandleReportGluedPathParams(t *testing.T) {
	sink := new(consumertest.MetricsSink)
	r := newHandlerReceiver(t, sink)

	target := "/data/report/stationtype=AMBWeatherV4.3.7&PASSKEY=84:F3:EB:9C:90:90&dateutc=2026-06-23+17:53:38&tempf=78.8&humidity=40&windspeedmph=3.8"
	req := httptest.NewRequest(http.MethodGet, target, http.NoBody)
	rec := httptest.NewRecorder()
	r.handleReport(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Len(t, sink.AllMetrics(), 1)
	md := sink.AllMetrics()[0]
	require.Positive(t, md.DataPointCount())

	// PASSKEY (no MAC field) becomes the station identity.
	assert.Equal(t, "84:F3:EB:9C:90:90", resourceMAC(t, md))

	byName := metricsByName(md)
	// tempf=78.8 -> (78.8-32)*5/9 = 26.0 Cel
	assert.InDelta(t, 26.0, dataPointByAttr(t, byName["ambientweather.temperature"], "sensor", "outdoor").DoubleValue(), 1e-6)
	assert.InDelta(t, 40.0, dataPointByAttr(t, byName["ambientweather.humidity"], "sensor", "outdoor").DoubleValue(), 1e-6)
}

func TestHandleReportRejectsUnsupportedMethod(t *testing.T) {
	sink := new(consumertest.MetricsSink)
	r := newHandlerReceiver(t, sink)

	req := httptest.NewRequest(http.MethodPut, "/data/report/", http.NoBody)
	rec := httptest.NewRecorder()
	r.handleReport(rec, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	assert.Empty(t, sink.AllMetrics())
}

func TestHandleReportEmptyPayload(t *testing.T) {
	sink := new(consumertest.MetricsSink)
	r := newHandlerReceiver(t, sink)

	req := httptest.NewRequest(http.MethodGet, "/data/report/", http.NoBody)
	rec := httptest.NewRecorder()
	r.handleReport(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "OK", rec.Body.String())
	assert.Empty(t, sink.AllMetrics())
}

func TestHandleReportConsumerError(t *testing.T) {
	r := newHandlerReceiver(t, consumertest.NewErr(errors.New("boom")))

	req := httptest.NewRequest(http.MethodGet, "/data/report/?"+sampleQuery().Encode(), http.NoBody)
	rec := httptest.NewRecorder()
	r.handleReport(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestHandleReportMultipleStations(t *testing.T) {
	sink := new(consumertest.MetricsSink)
	r := newHandlerReceiver(t, sink)

	for _, mac := range []string{"AA:AA", "BB:BB"} {
		q := sampleQuery()
		q.Set("MAC", mac)
		req := httptest.NewRequest(http.MethodGet, "/data/report/?"+q.Encode(), http.NoBody)
		rec := httptest.NewRecorder()
		r.handleReport(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	}

	require.Len(t, sink.AllMetrics(), 2)
	macs := []string{resourceMAC(t, sink.AllMetrics()[0]), resourceMAC(t, sink.AllMetrics()[1])}
	assert.ElementsMatch(t, []string{"AA:AA", "BB:BB"}, macs)
}

func resourceMAC(t *testing.T, md pmetric.Metrics) string {
	t.Helper()
	require.Equal(t, 1, md.ResourceMetrics().Len())
	v, ok := md.ResourceMetrics().At(0).Resource().Attributes().Get("ambientweather.station.mac")
	require.True(t, ok)
	return v.Str()
}

func TestStartShutdown(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.NetAddr.Endpoint = "localhost:0"
	r := newAmbientWeatherReceiver(cfg, receivertest.NewNopSettings(metadata.Type), consumertest.NewNop())

	require.NoError(t, r.Start(t.Context(), componenttest.NewNopHost()))
	require.NoError(t, r.Shutdown(t.Context()))
}
