// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ottlspan

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottltest"
)

var (
	traceID  = [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	traceID2 = [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	spanID   = [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
	spanID2  = [8]byte{8, 7, 6, 5, 4, 3, 2, 1}
)

func Test_newPathGetSetter(t *testing.T) {
	refSpan, refIS, refResource := createTelemetry()

	newAttrs := pcommon.NewMap()
	newAttrs.PutStr("hello", "world")

	newCache := pcommon.NewMap()
	newCache.PutStr("temp", "value")

	newEvents := ptrace.NewSpanEventSlice()
	newEvents.AppendEmpty().SetName("new event")

	newLinks := ptrace.NewSpanLinkSlice()
	newLinks.AppendEmpty().SetSpanID(spanID2)

	newStatus := ptrace.NewStatus()
	newStatus.SetMessage("new status")

	newPMap := pcommon.NewMap()
	pMap2 := newPMap.PutEmptyMap("k2")
	pMap2.PutStr("k1", "string")

	newMap := make(map[string]interface{})
	newMap2 := make(map[string]interface{})
	newMap2["k1"] = "string"
	newMap["k2"] = newMap2

	tests := []struct {
		name     string
		path     ottl.Path
		orig     interface{}
		newVal   interface{}
		modified func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map)
	}{
		{
			name: "cache",
			path: ottl.Path{
				Fields: []string{
					"cache",
				},
			},
			orig:   pcommon.NewMap(),
			newVal: newCache,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				newCache.CopyTo(cache)
			},
		},
		{
			name: "trace_id",
			path: ottl.Path{
				Fields: []string{
					"trace_id",
				},
			},
			orig:   pcommon.TraceID(traceID),
			newVal: pcommon.TraceID(traceID2),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetTraceID(traceID2)
			},
		},
		{
			name: "span_id",
			path: ottl.Path{
				Fields: []string{
					"span_id",
				},
			},
			orig:   pcommon.SpanID(spanID),
			newVal: pcommon.SpanID(spanID2),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetSpanID(spanID2)
			},
		},
		{
			name: "trace_id string",
			path: ottl.Path{
				Fields: []string{
					"trace_id",
					"string",
				},
			},
			orig:   hex.EncodeToString(traceID[:]),
			newVal: hex.EncodeToString(traceID2[:]),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetTraceID(traceID2)
			},
		},
		{
			name: "span_id string",
			path: ottl.Path{
				Fields: []string{
					"span_id",
					"string",
				},
			},
			orig:   hex.EncodeToString(spanID[:]),
			newVal: hex.EncodeToString(spanID2[:]),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetSpanID(spanID2)
			},
		},
		{
			name: "trace_state",
			path: ottl.Path{
				Fields: []string{
					"trace_state",
				},
			},
			orig:   "key1=val1,key2=val2",
			newVal: "key=newVal",
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.TraceState().FromRaw("key=newVal")
			},
		},
		{
			name: "parent_span_id",
			path: ottl.Path{
				Fields: []string{
					"parent_span_id",
				},
			},
			orig:   pcommon.SpanID(spanID2),
			newVal: pcommon.SpanID(spanID),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetParentSpanID(spanID)
			},
		},
		{
			name: "name",
			path: ottl.Path{
				Fields: []string{
					"name",
				},
			},
			orig:   "bear",
			newVal: "cat",
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetName("cat")
			},
		},
		{
			name: "kind",
			path: ottl.Path{
				Fields: []string{
					"kind",
				},
			},
			orig:   int64(2),
			newVal: int64(3),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetKind(ptrace.SpanKindClient)
			},
		},
		{
			name: "start_time_unix_nano",
			path: ottl.Path{
				Fields: []string{
					"start_time_unix_nano",
				},
			},
			orig:   int64(100_000_000),
			newVal: int64(200_000_000),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetStartTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(200)))
			},
		},
		{
			name: "end_time_unix_nano",
			path: ottl.Path{
				Fields: []string{
					"end_time_unix_nano",
				},
			},
			orig:   int64(500_000_000),
			newVal: int64(200_000_000),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetEndTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(200)))
			},
		},
		{
			name: "attributes",
			path: ottl.Path{
				Fields: []string{
					"attributes",
				},
			},
			orig:   refSpan.Attributes(),
			newVal: newAttrs,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				newAttrs.CopyTo(span.Attributes())
			},
		},
		{
			name: "dropped_attributes_count",
			path: ottl.Path{
				Fields: []string{
					"dropped_attributes_count",
				},
			},
			orig:   int64(10),
			newVal: int64(20),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetDroppedAttributesCount(20)
			},
		},
		{
			name: "events",
			path: ottl.Path{
				Fields: []string{
					"events",
				},
			},
			orig:   refSpan.Events(),
			newVal: newEvents,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.Events().RemoveIf(func(_ ptrace.SpanEvent) bool {
					return true
				})
				newEvents.CopyTo(span.Events())
			},
		},
		{
			name: "dropped_events_count",
			path: ottl.Path{
				Fields: []string{
					"dropped_events_count",
				},
			},
			orig:   int64(20),
			newVal: int64(30),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetDroppedEventsCount(30)
			},
		},
		{
			name: "links",
			path: ottl.Path{
				Fields: []string{
					"links",
				},
			},
			orig:   refSpan.Links(),
			newVal: newLinks,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.Links().RemoveIf(func(_ ptrace.SpanLink) bool {
					return true
				})
				newLinks.CopyTo(span.Links())
			},
		},
		{
			name: "dropped_links_count",
			path: ottl.Path{
				Fields: []string{
					"dropped_links_count",
				},
			},
			orig:   int64(30),
			newVal: int64(40),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.SetDroppedLinksCount(40)
			},
		},
		{
			name: "status",
			path: ottl.Path{
				Fields: []string{
					"status",
				},
			},
			orig:   refSpan.Status(),
			newVal: newStatus,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				newStatus.CopyTo(span.Status())
			},
		},
		{
			name: "status code",
			path: ottl.Path{
				Fields: []string{
					"status",
					"code",
				},
			},
			orig:   int64(ptrace.StatusCodeOk),
			newVal: int64(ptrace.StatusCodeError),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.Status().SetCode(ptrace.StatusCodeError)
			},
		},
		{
			name: "status message",
			path: ottl.Path{
				Fields: []string{
					"status",
					"message",
				},
			},
			orig:   "good span",
			newVal: "bad span",
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				span.Status().SetMessage("bad span")
			},
		},
		{
			name: "instrumentation_scope",
			path: ottl.Path{
				Fields: []string{
					"instrumentation_scope",
				},
			},
			orig:   refIS,
			newVal: pcommon.NewInstrumentationScope(),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				pcommon.NewInstrumentationScope().CopyTo(il)
			},
		},
		{
			name: "resource",
			path: ottl.Path{
				Fields: []string{
					"resource",
				},
			},
			orig:   refResource,
			newVal: pcommon.NewResource(),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				pcommon.NewResource().CopyTo(resource)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessor, err := newPathGetSetter(tt.path)
			assert.NoError(t, err)

			span, il, resource := createTelemetry()

			tCtx := NewTransformContext(span, il, resource)

			got, err := accessor.Get(context.Background(), tCtx)
			assert.Nil(t, err)
			assert.Equal(t, tt.orig, got)

			err = accessor.Set(context.Background(), tCtx, tt.newVal)
			assert.Nil(t, err)

			exSpan, exIl, exRes := createTelemetry()
			exCache := pcommon.NewMap()
			tt.modified(exSpan, exIl, exRes, exCache)

			assert.Equal(t, exSpan, span)
			assert.Equal(t, exIl, il)
			assert.Equal(t, exRes, resource)
			assert.Equal(t, exCache, tCtx.getCache())
		})
	}
}

func createTelemetry() (ptrace.Span, pcommon.InstrumentationScope, pcommon.Resource) {
	span := ptrace.NewSpan()
	span.SetTraceID(traceID)
	span.SetSpanID(spanID)
	span.TraceState().FromRaw("key1=val1,key2=val2")
	span.SetParentSpanID(spanID2)
	span.SetName("bear")
	span.SetKind(ptrace.SpanKindServer)
	span.SetStartTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(100)))
	span.SetEndTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(500)))
	span.Attributes().PutStr("str", "val")
	span.Attributes().PutBool("bool", true)
	span.Attributes().PutInt("int", 10)
	span.Attributes().PutDouble("double", 1.2)
	span.Attributes().PutEmptyBytes("bytes").FromRaw([]byte{1, 3, 2})

	arrStr := span.Attributes().PutEmptySlice("arr_str")
	arrStr.AppendEmpty().SetStr("one")
	arrStr.AppendEmpty().SetStr("two")

	arrBool := span.Attributes().PutEmptySlice("arr_bool")
	arrBool.AppendEmpty().SetBool(true)
	arrBool.AppendEmpty().SetBool(false)

	arrInt := span.Attributes().PutEmptySlice("arr_int")
	arrInt.AppendEmpty().SetInt(2)
	arrInt.AppendEmpty().SetInt(3)

	arrFloat := span.Attributes().PutEmptySlice("arr_float")
	arrFloat.AppendEmpty().SetDouble(1.0)
	arrFloat.AppendEmpty().SetDouble(2.0)

	arrBytes := span.Attributes().PutEmptySlice("arr_bytes")
	arrBytes.AppendEmpty().SetEmptyBytes().FromRaw([]byte{1, 2, 3})
	arrBytes.AppendEmpty().SetEmptyBytes().FromRaw([]byte{2, 3, 4})

	pMap := span.Attributes().PutEmptyMap("pMap")
	pMap.PutStr("original", "map")

	m := span.Attributes().PutEmptyMap("map")
	m.PutStr("original", "map")

	span.SetDroppedAttributesCount(10)

	span.Events().AppendEmpty().SetName("event")
	span.SetDroppedEventsCount(20)

	span.Links().AppendEmpty().SetTraceID(traceID)
	span.SetDroppedLinksCount(30)

	span.Status().SetCode(ptrace.StatusCodeOk)
	span.Status().SetMessage("good span")

	il := pcommon.NewInstrumentationScope()
	il.SetName("library")
	il.SetVersion("version")

	resource := pcommon.NewResource()
	span.Attributes().CopyTo(resource.Attributes())

	return span, il, resource
}

func Test_ParseEnum(t *testing.T) {
	tests := []struct {
		name string
		want ottl.Enum
	}{
		{
			name: "SPAN_KIND_UNSPECIFIED",
			want: ottl.Enum(ptrace.SpanKindUnspecified),
		},
		{
			name: "SPAN_KIND_INTERNAL",
			want: ottl.Enum(ptrace.SpanKindInternal),
		},
		{
			name: "SPAN_KIND_SERVER",
			want: ottl.Enum(ptrace.SpanKindServer),
		},
		{
			name: "SPAN_KIND_CLIENT",
			want: ottl.Enum(ptrace.SpanKindClient),
		},
		{
			name: "SPAN_KIND_PRODUCER",
			want: ottl.Enum(ptrace.SpanKindProducer),
		},
		{
			name: "SPAN_KIND_CONSUMER",
			want: ottl.Enum(ptrace.SpanKindConsumer),
		},
		{
			name: "STATUS_CODE_UNSET",
			want: ottl.Enum(ptrace.StatusCodeUnset),
		},
		{
			name: "STATUS_CODE_OK",
			want: ottl.Enum(ptrace.StatusCodeOk),
		},
		{
			name: "STATUS_CODE_ERROR",
			want: ottl.Enum(ptrace.StatusCodeError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parseEnum((*ottl.EnumSymbol)(ottltest.Strp(tt.name)))
			assert.NoError(t, err)
			assert.Equal(t, *actual, tt.want)
		})
	}
}

func Test_ParseEnum_False(t *testing.T) {
	tests := []struct {
		name       string
		enumSymbol *ottl.EnumSymbol
	}{
		{
			name:       "unknown enum symbol",
			enumSymbol: (*ottl.EnumSymbol)(ottltest.Strp("not an enum")),
		},
		{
			name:       "nil enum symbol",
			enumSymbol: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := parseEnum(tt.enumSymbol)
			assert.Error(t, err)
			assert.Nil(t, actual)
		})
	}
}
