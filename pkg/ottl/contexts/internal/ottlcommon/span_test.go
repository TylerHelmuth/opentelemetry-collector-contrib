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

package ottlcommon

import (
	"context"
	"encoding/hex"
	"go.opentelemetry.io/otel/trace"
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

func TestSpanPathGetSetter(t *testing.T) {
	refSpan := createSpan()

	refTraceState, _ := trace.ParseTraceState(refSpan.TraceState().AsRaw())
	newTraceState, _ := trace.ParseTraceState("key1=newVal,key2=val2")

	newAttrs := pcommon.NewMap()
	newAttrs.PutStr("hello", "world")

	newEvents := ptrace.NewSpanEventSlice()
	newEvents.AppendEmpty().SetName("new event")

	newLinks := ptrace.NewSpanLinkSlice()
	newLinks.AppendEmpty().SetSpanID(spanID2)

	newStatus := ptrace.NewStatus()
	newStatus.SetMessage("new status")

	tests := []struct {
		name     string
		path     ottl.Path
		orig     interface{}
		newVal   interface{}
		modified func(span ptrace.Span)
	}{
		{
			name: "trace_id",
			path: ottl.Path{
				Fields: []string{
					"trace_id",
				},
			},
			orig:   pcommon.TraceID(traceID),
			newVal: pcommon.TraceID(traceID2),
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
				span.TraceState().FromRaw("key=newVal")
			},
		},
		{
			name: "trace_state key",
			path: ottl.Path{
				Fields: []string{
					"trace_state",
				},
				Keys: []ottl.Key{
					{
						Map: ottltest.Strp("key1"),
					},
				},
			},
			orig:   refTraceState,
			newVal: newTraceState,
			modified: func(span ptrace.Span) {
				span.TraceState().FromRaw("key1=newVal,key2=val2")
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
			modified: func(span ptrace.Span) {
				span.SetParentSpanID(spanID)
			},
		},
		{
			name: "parent_span_id string",
			path: ottl.Path{
				Fields: []string{
					"parent_span_id",
					"string",
				},
			},
			orig:   hex.EncodeToString(spanID2[:]),
			newVal: hex.EncodeToString(spanID[:]),
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
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
			modified: func(span ptrace.Span) {
				span.Status().SetMessage("bad span")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessor, err := SpanPathGetSetter[*spanContext](tt.path)
			assert.NoError(t, err)

			span := createSpan()

			got, err := accessor.Get(context.Background(), newSpanContext(span))
			assert.NoError(t, err)
			assert.Equal(t, tt.orig, got)

			err = accessor.Set(context.Background(), newSpanContext(span), tt.newVal)
			assert.NoError(t, err)

			expectedSpan := createSpan()
			tt.modified(expectedSpan)

			assert.Equal(t, expectedSpan, span)
		})
	}
}

func createSpan() ptrace.Span {
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

	span.Attributes().PutEmptySlice("arr_empty")

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

	span.SetDroppedAttributesCount(10)

	span.Events().AppendEmpty().SetName("event")
	span.SetDroppedEventsCount(20)

	span.Links().AppendEmpty().SetTraceID(traceID)
	span.SetDroppedLinksCount(30)

	span.Status().SetCode(ptrace.StatusCodeOk)
	span.Status().SetMessage("good span")

	return span
}

type spanContext struct {
	span ptrace.Span
}

func (r *spanContext) GetSpan() ptrace.Span {
	return r.span
}

func newSpanContext(span ptrace.Span) *spanContext {
	return &spanContext{span: span}
}
