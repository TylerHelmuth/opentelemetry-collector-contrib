package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/open-telemetry/opamp-go/client"
	"github.com/open-telemetry/opamp-go/client/types"
	"github.com/open-telemetry/opamp-go/protobufs"
	"go.opentelemetry.io/collector/confmap"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

func NewFactory() confmap.ProviderFactory {
	return confmap.NewProviderFactory(newProvider)
}

func newProvider(confmap.ProviderSettings) confmap.Provider {
	return newOpampProvider()
}

type opampProvider struct {
	instanceID uuid.UUID

	opampClient client.OpAMPClient

	remoteConfigChan chan string
	remoteConfig     string

	logger *zap.Logger
}

func newOpampProvider() *opampProvider {
	return &opampProvider{
		instanceID:       uuid.New(),
		remoteConfigChan: make(chan string, 1),
		logger:           zap.NewNop(),
	}
}

func (o *opampProvider) Start() {
	o.opampClient = client.NewWebSocket(nil)
	settings := types.StartSettings{
		Header:         http.Header{},
		OpAMPServerURL: "ws://0.0.0.0:4320/v1/opamp",
		InstanceUid:    types.InstanceUid(o.instanceID),
		TLSConfig:      nil,
		Callbacks: types.Callbacks{
			OnConnect: func(_ context.Context) {
				o.logger.Debug("Connected to the OpAMP server")
			},
			OnConnectFailed: func(_ context.Context, err error) {
				o.logger.Error("Failed to connect to the OpAMP server", zap.Error(err))
			},
			OnError: func(_ context.Context, err *protobufs.ServerErrorResponse) {
				o.logger.Error("OpAMP server returned an error response", zap.String("message", err.ErrorMessage))
			},
			GetEffectiveConfig: func(_ context.Context) (*protobufs.EffectiveConfig, error) {
				conf, _ := yaml.Marshal(map[string]any{
					"exporters": map[string]any{
						"debug/test": map[string]any{
							"verbosity": "detailed",
						},
					},
				})
				return &protobufs.EffectiveConfig{
					ConfigMap: &protobufs.AgentConfigMap{
						ConfigMap: map[string]*protobufs.AgentConfigFile{
							"": {
								Body:        conf,
								ContentType: "text/yaml",
							},
						},
					},
				}, nil
			},
			OnMessage: o.onMessage,
		},
	}
	_ = o.opampClient.SetAgentDescription(&protobufs.AgentDescription{
		IdentifyingAttributes: []*protobufs.KeyValue{
			{
				Key: string(semconv.ServiceInstanceIDKey),
				Value: &protobufs.AnyValue{
					Value: &protobufs.AnyValue_StringValue{StringValue: o.instanceID.String()},
				},
			},
			{
				Key: string(semconv.ServiceNameKey),
				Value: &protobufs.AnyValue{
					Value: &protobufs.AnyValue_StringValue{StringValue: "primary-collector"},
				},
			},
			{
				Key: string(semconv.ServiceVersionKey),
				Value: &protobufs.AnyValue{
					Value: &protobufs.AnyValue_StringValue{StringValue: "v0.141.0"},
				},
			},
		},
	})
	capabilities := protobufs.AgentCapabilities_AgentCapabilities_AcceptsRemoteConfig | protobufs.AgentCapabilities_AgentCapabilities_ReportsEffectiveConfig
	_ = o.opampClient.SetCapabilities(&capabilities)
	_ = o.opampClient.Start(context.Background(), settings)
}

func (o *opampProvider) Shutdown(ctx context.Context) error {
	return o.opampClient.Stop(ctx)
}

func (*opampProvider) Scheme() string {
	return "opamp"
}

func (o *opampProvider) onMessage(_ context.Context, msg *types.MessageData) {
	if msg.RemoteConfig != nil && o.remoteConfig != msg.RemoteConfig.String() {
		o.remoteConfig = msg.RemoteConfig.String()
		o.remoteConfigChan <- o.remoteConfig
	}
}

func (o *opampProvider) Retrieve(_ context.Context, _ string, watcher confmap.WatcherFunc) (*confmap.Retrieved, error) {
	closeCh := make(chan struct{})
	startWatcher(o.remoteConfigChan, closeCh, watcher)
	o.Start()
	return confmap.NewRetrieved(map[string]any{
		"exporters": map[string]any{
			"debug/test": map[string]any{
				"verbosity": "detailed",
			},
		},
	}, confmap.WithRetrievedClose(func(_ context.Context) error {
		close(closeCh)
		_ = o.Shutdown(context.Background())
		return nil
	}))
}

func startWatcher(watchCh <-chan string, closeCh <-chan struct{}, watcher confmap.WatcherFunc) {
	go func() {
		select {
		case <-closeCh:
			return
		case _, ok := <-watchCh:
			if !ok {
				// Channel close without any event, connection must have been closed.
				return
			}
			watcher(&confmap.ChangeEvent{})
			return
		}
	}()
}
