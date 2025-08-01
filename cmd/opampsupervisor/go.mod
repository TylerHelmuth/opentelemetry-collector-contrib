module github.com/open-telemetry/opentelemetry-collector-contrib/cmd/opampsupervisor

go 1.23.7

require (
	github.com/google/uuid v1.6.0
	github.com/knadh/koanf/maps v0.1.2
	github.com/knadh/koanf/parsers/yaml v1.1.0
	github.com/knadh/koanf/providers/file v1.2.0
	github.com/knadh/koanf/providers/rawbytes v1.0.0
	github.com/knadh/koanf/v2 v2.2.2
	github.com/open-telemetry/opamp-go v0.20.0
	github.com/open-telemetry/opentelemetry-collector-contrib/testbed v0.131.0
	github.com/stretchr/testify v1.10.0
	go.opentelemetry.io/collector/component v1.37.0
	go.opentelemetry.io/collector/config/configopaque v1.37.0
	go.opentelemetry.io/collector/config/configtelemetry v0.131.0
	go.opentelemetry.io/collector/config/configtls v1.37.0
	go.opentelemetry.io/collector/confmap v1.37.0
	go.opentelemetry.io/collector/confmap/provider/envprovider v1.37.0
	go.opentelemetry.io/collector/confmap/provider/fileprovider v1.37.0
	go.opentelemetry.io/collector/pdata v1.37.0
	go.opentelemetry.io/collector/service v0.131.0
	go.opentelemetry.io/contrib/bridges/otelzap v0.12.0
	go.opentelemetry.io/contrib/otelconf v0.17.0
	go.opentelemetry.io/otel v1.37.0
	go.opentelemetry.io/otel/log v0.13.0
	go.opentelemetry.io/otel/trace v1.37.0
	go.uber.org/goleak v1.3.0
	go.uber.org/multierr v1.11.0
	go.uber.org/zap v1.27.0
	golang.org/x/sys v0.34.0
	google.golang.org/protobuf v1.36.6
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/alecthomas/participle/v2 v2.1.4 // indirect
	github.com/antchfx/xmlquery v1.4.4 // indirect
	github.com/antchfx/xpath v1.3.4 // indirect
	github.com/apache/thrift v0.22.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/elastic/go-grok v0.3.1 // indirect
	github.com/elastic/lunes v0.1.0 // indirect
	github.com/expr-lang/expr v1.17.5 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/foxboron/go-tpm-keyfiles v0.0.0-20250323135004-b31fac66206e // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/go-tpm v0.9.5 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.1 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jaegertracing/jaeger-idl v0.5.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/knadh/koanf/providers/confmap v1.0.0 // indirect
	github.com/leodido/go-syslog/v4 v4.2.0 // indirect
	github.com/leodido/ragel-machinery v0.0.0-20190525184631-5f46317e436b // indirect
	github.com/lightstep/go-expohisto v1.0.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20240909124753-873cd0166683 // indirect
	github.com/magefile/mage v1.15.0 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/mostynb/go-grpc-compression v1.2.3 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/routingconnector v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/connector/spanmetricsconnector v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opencensusexporter v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/stefexporter v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/syslogexporter v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/exporter/zipkinexporter v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/pdatautil v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/core/xidutils v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/stefreceiver v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver v0.131.0 // indirect
	github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver v0.131.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.65.0 // indirect
	github.com/prometheus/procfs v0.17.0 // indirect
	github.com/rs/cors v1.11.1 // indirect
	github.com/shirou/gopsutil/v4 v4.25.6 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/splunk/stef/go/grpc v0.0.6 // indirect
	github.com/splunk/stef/go/otel v0.0.6 // indirect
	github.com/splunk/stef/go/pdata v0.0.6 // indirect
	github.com/splunk/stef/go/pkg v0.0.6 // indirect
	github.com/tklauser/go-sysconf v0.3.14 // indirect
	github.com/tklauser/numcpus v0.9.0 // indirect
	github.com/twmb/murmur3 v1.1.8 // indirect
	github.com/ua-parser/uap-go v0.0.0-20240611065828-3a4781585db6 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/collector v0.131.0 // indirect
	go.opentelemetry.io/collector/client v1.37.0 // indirect
	go.opentelemetry.io/collector/component/componentstatus v0.131.0 // indirect
	go.opentelemetry.io/collector/component/componenttest v0.131.0 // indirect
	go.opentelemetry.io/collector/config/configauth v0.131.0 // indirect
	go.opentelemetry.io/collector/config/configcompression v1.37.0 // indirect
	go.opentelemetry.io/collector/config/configgrpc v0.131.0 // indirect
	go.opentelemetry.io/collector/config/confighttp v0.131.0 // indirect
	go.opentelemetry.io/collector/config/configmiddleware v0.131.0 // indirect
	go.opentelemetry.io/collector/config/confignet v1.37.0 // indirect
	go.opentelemetry.io/collector/config/configoptional v0.131.0 // indirect
	go.opentelemetry.io/collector/config/configretry v1.37.0 // indirect
	go.opentelemetry.io/collector/confmap/xconfmap v0.131.0 // indirect
	go.opentelemetry.io/collector/connector v0.131.0 // indirect
	go.opentelemetry.io/collector/connector/connectortest v0.131.0 // indirect
	go.opentelemetry.io/collector/connector/xconnector v0.131.0 // indirect
	go.opentelemetry.io/collector/consumer v1.37.0 // indirect
	go.opentelemetry.io/collector/consumer/consumererror v0.131.0 // indirect
	go.opentelemetry.io/collector/consumer/consumererror/xconsumererror v0.131.0 // indirect
	go.opentelemetry.io/collector/consumer/consumertest v0.131.0 // indirect
	go.opentelemetry.io/collector/consumer/xconsumer v0.131.0 // indirect
	go.opentelemetry.io/collector/exporter v0.131.0 // indirect
	go.opentelemetry.io/collector/exporter/debugexporter v0.131.0 // indirect
	go.opentelemetry.io/collector/exporter/exporterhelper/xexporterhelper v0.131.0 // indirect
	go.opentelemetry.io/collector/exporter/exportertest v0.131.0 // indirect
	go.opentelemetry.io/collector/exporter/otlpexporter v0.131.0 // indirect
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.131.0 // indirect
	go.opentelemetry.io/collector/exporter/xexporter v0.131.0 // indirect
	go.opentelemetry.io/collector/extension v1.37.0 // indirect
	go.opentelemetry.io/collector/extension/extensionauth v1.37.0 // indirect
	go.opentelemetry.io/collector/extension/extensioncapabilities v0.131.0 // indirect
	go.opentelemetry.io/collector/extension/extensionmiddleware v0.131.0 // indirect
	go.opentelemetry.io/collector/extension/extensiontest v0.131.0 // indirect
	go.opentelemetry.io/collector/extension/xextension v0.131.0 // indirect
	go.opentelemetry.io/collector/extension/zpagesextension v0.131.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.37.0 // indirect
	go.opentelemetry.io/collector/internal/fanoutconsumer v0.131.0 // indirect
	go.opentelemetry.io/collector/internal/memorylimiter v0.131.0 // indirect
	go.opentelemetry.io/collector/internal/sharedcomponent v0.131.0 // indirect
	go.opentelemetry.io/collector/internal/telemetry v0.131.0 // indirect
	go.opentelemetry.io/collector/otelcol v0.131.0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.131.0 // indirect
	go.opentelemetry.io/collector/pdata/testdata v0.131.0 // indirect
	go.opentelemetry.io/collector/pdata/xpdata v0.131.0 // indirect
	go.opentelemetry.io/collector/pipeline v0.131.0 // indirect
	go.opentelemetry.io/collector/pipeline/xpipeline v0.131.0 // indirect
	go.opentelemetry.io/collector/processor v1.37.0 // indirect
	go.opentelemetry.io/collector/processor/batchprocessor v0.131.0 // indirect
	go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.131.0 // indirect
	go.opentelemetry.io/collector/processor/processorhelper v0.131.0 // indirect
	go.opentelemetry.io/collector/processor/processorhelper/xprocessorhelper v0.131.0 // indirect
	go.opentelemetry.io/collector/processor/processortest v0.131.0 // indirect
	go.opentelemetry.io/collector/processor/xprocessor v0.131.0 // indirect
	go.opentelemetry.io/collector/receiver v1.37.0 // indirect
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.131.0 // indirect
	go.opentelemetry.io/collector/receiver/receiverhelper v0.131.0 // indirect
	go.opentelemetry.io/collector/receiver/receivertest v0.131.0 // indirect
	go.opentelemetry.io/collector/receiver/xreceiver v0.131.0 // indirect
	go.opentelemetry.io/collector/service/hostcapabilities v0.131.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.62.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.62.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.37.0 // indirect
	go.opentelemetry.io/contrib/zpages v0.62.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.59.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.13.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.37.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.37.0 // indirect
	go.opentelemetry.io/otel/metric v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk v1.37.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.13.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.37.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.0 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/exp v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	gonum.org/v1/gonum v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/grpc v1.74.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	modernc.org/b/v2 v2.1.0 // indirect
	sigs.k8s.io/yaml v1.5.0 // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/receiver/syslogreceiver => ../../receiver/syslogreceiver

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/common => ../../internal/common

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza => ../../pkg/stanza

replace github.com/open-telemetry/opentelemetry-collector-contrib/extension/storage => ../../extension/storage

replace github.com/open-telemetry/opentelemetry-collector-contrib/connector/routingconnector => ../../connector/routingconnector

replace github.com/open-telemetry/opentelemetry-collector-contrib/receiver/signalfxreceiver => ../../receiver/signalfxreceiver

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheusremotewrite => ../../pkg/translator/prometheusremotewrite

replace github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusexporter => ../../exporter/prometheusexporter

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/prometheus => ../../pkg/translator/prometheus

replace github.com/open-telemetry/opentelemetry-collector-contrib/exporter/zipkinexporter => ../../exporter/zipkinexporter

replace github.com/open-telemetry/opentelemetry-collector-contrib/exporter/syslogexporter => ../../exporter/syslogexporter

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchperresourceattr => ../../pkg/batchperresourceattr

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/zipkin => ../../pkg/translator/zipkin

replace github.com/open-telemetry/opentelemetry-collector-contrib/exporter/splunkhecexporter => ../../exporter/splunkhecexporter

replace github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusreceiver => ../../receiver/prometheusreceiver

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger => ../../pkg/translator/jaeger

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/pdatautil => ../../internal/pdatautil

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/experimentalmetricmetadata => ../../pkg/experimentalmetricmetadata

replace github.com/open-telemetry/opentelemetry-collector-contrib/testbed/mockdatasenders/mockdatadogagentexporter => ../../testbed/mockdatasenders/mockdatadogagentexporter

replace github.com/open-telemetry/opentelemetry-collector-contrib/exporter/carbonexporter => ../../exporter/carbonexporter

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/resourcetotelemetry => ../../pkg/resourcetotelemetry

replace github.com/open-telemetry/opentelemetry-collector-contrib/receiver/zipkinreceiver => ../../receiver/zipkinreceiver

replace github.com/open-telemetry/opentelemetry-collector-contrib/receiver/opencensusreceiver => ../../receiver/opencensusreceiver

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/sharedcomponent => ../../internal/sharedcomponent

replace github.com/open-telemetry/opentelemetry-collector-contrib/testbed => ../../testbed

replace github.com/open-telemetry/opentelemetry-collector-contrib/receiver/datadogreceiver => ../../receiver/datadogreceiver

replace github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver => ../../receiver/carbonreceiver

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/opencensus => ../../pkg/translator/opencensus

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatatest => ../../pkg/pdatatest

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/pdatautil => ../../pkg/pdatautil

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/core/xidutils => ../../pkg/core/xidutils

replace github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opencensusexporter => ../../exporter/opencensusexporter

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk => ../../internal/splunk

replace github.com/open-telemetry/opentelemetry-collector-contrib/receiver/splunkhecreceiver => ../../receiver/splunkhecreceiver

replace github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter => ../../exporter/signalfxexporter

replace github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jaegerreceiver => ../../receiver/jaegerreceiver

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl => ../../pkg/ottl

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/golden => ../../pkg/golden

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal => ../../internal/coreinternal

replace github.com/open-telemetry/opentelemetry-collector-contrib/connector/spanmetricsconnector => ../../connector/spanmetricsconnector

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/signalfx => ../../pkg/translator/signalfx

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/exp/metrics => ../../internal/exp/metrics

replace github.com/open-telemetry/opentelemetry-collector-contrib/extension/ackextension => ../../extension/ackextension

replace github.com/open-telemetry/opentelemetry-collector-contrib/exporter/prometheusremotewriteexporter => ../../exporter/prometheusremotewriteexporter

replace github.com/open-telemetry/opentelemetry-collector-contrib/exporter/stefexporter => ../../exporter/stefexporter

replace github.com/open-telemetry/opentelemetry-collector-contrib/receiver/stefreceiver => ../../receiver/stefreceiver

replace github.com/open-telemetry/opentelemetry-collector-contrib/processor/deltatocumulativeprocessor => ../../processor/deltatocumulativeprocessor

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/gopsutilenv => ../../internal/gopsutilenv

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/metadataproviders => ../../internal/metadataproviders

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig => ../../internal/k8sconfig

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/datadog => ../../internal/datadog

replace github.com/open-telemetry/opentelemetry-collector-contrib/pkg/datadog => ../../pkg/datadog

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/ecsutil => ../../internal/aws/ecsutil
