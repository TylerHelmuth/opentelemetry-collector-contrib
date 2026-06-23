// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ambientweatherreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver"

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componentstatus"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receiverhelper"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver/internal/metadata"
)

var (
	_ receiver.Metrics = (*ambientWeatherReceiver)(nil)

	errInvalidMethod = errors.New("only GET and POST methods are supported")
)

// ambientWeatherReceiver hosts an HTTP server that accepts AMBWeather-format
// uploads pushed by an Ambient Weather console in "custom server" mode.
type ambientWeatherReceiver struct {
	cfg          *Config
	settings     receiver.Settings
	logger       *zap.Logger
	nextConsumer consumer.Metrics
	obsrecv      *receiverhelper.ObsReport
	server       *http.Server
	shutdownWG   sync.WaitGroup
}

func newAmbientWeatherReceiver(cfg *Config, settings receiver.Settings, nextConsumer consumer.Metrics) *ambientWeatherReceiver {
	return &ambientWeatherReceiver{
		cfg:          cfg,
		settings:     settings,
		logger:       settings.Logger,
		nextConsumer: nextConsumer,
	}
}

// Start brings up the HTTP server that listens for console uploads.
func (r *ambientWeatherReceiver) Start(ctx context.Context, host component.Host) error {
	var err error
	r.obsrecv, err = receiverhelper.NewObsReport(receiverhelper.ObsReportSettings{
		ReceiverID:             r.settings.ID,
		Transport:              "http",
		ReceiverCreateSettings: r.settings,
	})
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc(r.cfg.Path, r.handleReport)

	r.server, err = r.cfg.ToServer(ctx, host.GetExtensions(), r.settings.TelemetrySettings, mux)
	if err != nil {
		return err
	}

	listener, err := r.cfg.ToListener(ctx)
	if err != nil {
		return err
	}

	r.shutdownWG.Go(func() {
		if errServe := r.server.Serve(listener); errServe != nil && !errors.Is(errServe, http.ErrServerClosed) {
			componentstatus.ReportStatus(host, componentstatus.NewFatalErrorEvent(errServe))
		}
	})
	return nil
}

// Shutdown stops the HTTP server.
func (r *ambientWeatherReceiver) Shutdown(ctx context.Context) error {
	if r.server == nil {
		return nil
	}
	err := r.server.Shutdown(ctx)
	r.shutdownWG.Wait()
	return err
}

// handleReport parses a console upload and forwards the resulting metrics.
func (r *ambientWeatherReceiver) handleReport(w http.ResponseWriter, req *http.Request) {
	ctx := r.obsrecv.StartMetricsOp(req.Context())

	if req.Method != http.MethodGet && req.Method != http.MethodPost {
		r.obsrecv.EndMetricsOp(ctx, metadata.Type.String(), 0, errInvalidMethod)
		http.Error(w, errInvalidMethod.Error(), http.StatusMethodNotAllowed)
		return
	}

	fields, err := r.extractParams(req)
	if err != nil {
		r.obsrecv.EndMetricsOp(ctx, metadata.Type.String(), 0, err)
		http.Error(w, "unable to parse request", http.StatusBadRequest)
		return
	}

	md := r.buildMetrics(fields)
	dataPoints := md.DataPointCount()

	if dataPoints == 0 {
		// Acknowledge the upload so the console does not retry, but forward nothing.
		r.logger.Debug("no recognized fields in upload", zap.String("path", req.URL.Path))
		r.obsrecv.EndMetricsOp(ctx, metadata.Type.String(), 0, nil)
		_, _ = w.Write([]byte("OK"))
		return
	}

	if err := r.nextConsumer.ConsumeMetrics(ctx, md); err != nil {
		r.obsrecv.EndMetricsOp(ctx, metadata.Type.String(), dataPoints, err)
		http.Error(w, "failed to consume metrics", http.StatusInternalServerError)
		return
	}

	r.obsrecv.EndMetricsOp(ctx, metadata.Type.String(), dataPoints, nil)
	_, _ = w.Write([]byte("OK"))
}

// extractParams returns the AMBWeather key/value pairs from a console upload.
//
// Ambient Weather consoles build the upload URL by concatenating the configured
// "Path" with the parameters. When the path is set with a trailing "?" the params
// arrive as a normal query string; when it is not (the common real-world case) the
// params are glued directly onto the path with no separator, leaving the query
// string empty (e.g. "GET /data/report/tempf=70&humidity=55"). Both forms are
// handled here. A urlencoded POST body (Ecowitt-style) is also merged in; it is
// parsed explicitly rather than via http.Request.ParseForm because consoles often
// omit or misreport the Content-Type header.
func (r *ambientWeatherReceiver) extractParams(req *http.Request) (url.Values, error) {
	rawQuery := req.URL.RawQuery
	if rawQuery == "" {
		// No "?" in the request: recover the params the console glued onto the path.
		// The raw request target is used so percent/plus encoding is preserved for
		// url.ParseQuery to decode.
		target := req.RequestURI
		if i := strings.IndexByte(target, '?'); i >= 0 {
			rawQuery = target[i+1:]
		} else {
			rawQuery = strings.TrimPrefix(target, r.cfg.Path)
		}
	}

	fields, err := url.ParseQuery(rawQuery)
	if err != nil {
		return nil, err
	}

	if req.Method == http.MethodPost {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if len(body) > 0 {
			bodyValues, err := url.ParseQuery(string(body))
			if err != nil {
				return nil, err
			}
			for key, values := range bodyValues {
				for _, value := range values {
					fields.Add(key, value)
				}
			}
		}
	}

	return fields, nil
}
