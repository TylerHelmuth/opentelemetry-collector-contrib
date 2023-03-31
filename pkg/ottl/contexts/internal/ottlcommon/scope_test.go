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
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

func TestScopePathGetSetter(t *testing.T) {
	refIS := createInstrumentationScope()

	newAttrs := pcommon.NewMap()
	newAttrs.PutStr("hello", "world")
	tests := []struct {
		name     string
		path     ottl.Path
		orig     interface{}
		newVal   interface{}
		modified func(is pcommon.InstrumentationScope)
	}{
		{
			name:   "instrumentation_scope",
			path:   ottl.Path{},
			orig:   refIS,
			newVal: pcommon.NewInstrumentationScope(),
			modified: func(is pcommon.InstrumentationScope) {
				pcommon.NewInstrumentationScope().CopyTo(is)
			},
		},
		{
			name: "instrumentation_scope name",
			path: ottl.Path{
				Fields: []string{
					"name",
				},
			},
			orig:   refIS.Name(),
			newVal: "newname",
			modified: func(is pcommon.InstrumentationScope) {
				is.SetName("newname")
			},
		},
		{
			name: "instrumentation_scope version",
			path: ottl.Path{
				Fields: []string{
					"version",
				},
			},
			orig:   refIS.Version(),
			newVal: "next",
			modified: func(is pcommon.InstrumentationScope) {
				is.SetVersion("next")
			},
		},
		{
			name: "attributes",
			path: ottl.Path{
				Fields: []string{
					"attributes",
				},
			},
			orig:   refIS.Attributes(),
			newVal: newAttrs,
			modified: func(is pcommon.InstrumentationScope) {
				newAttrs.CopyTo(is.Attributes())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessor, err := ScopePathGetSetter[*instrumentationScopeContext](tt.path)
			assert.NoError(t, err)

			is := createInstrumentationScope()

			got, err := accessor.Get(context.Background(), newInstrumentationScopeContext(is))
			assert.Nil(t, err)
			assert.Equal(t, tt.orig, got)

			err = accessor.Set(context.Background(), newInstrumentationScopeContext(is), tt.newVal)
			assert.Nil(t, err)

			expectedIS := createInstrumentationScope()
			tt.modified(expectedIS)

			assert.Equal(t, expectedIS, is)
		})
	}
}

func createInstrumentationScope() pcommon.InstrumentationScope {
	is := pcommon.NewInstrumentationScope()
	is.SetName("library")
	is.SetVersion("version")
	is.SetDroppedAttributesCount(10)

	is.Attributes().PutStr("str", "val")
	is.Attributes().PutBool("bool", true)
	is.Attributes().PutInt("int", 10)
	is.Attributes().PutDouble("double", 1.2)
	is.Attributes().PutEmptyBytes("bytes").FromRaw([]byte{1, 3, 2})

	is.Attributes().PutEmptySlice("arr_empty")

	arrStr := is.Attributes().PutEmptySlice("arr_str")
	arrStr.AppendEmpty().SetStr("one")
	arrStr.AppendEmpty().SetStr("two")

	arrBool := is.Attributes().PutEmptySlice("arr_bool")
	arrBool.AppendEmpty().SetBool(true)
	arrBool.AppendEmpty().SetBool(false)

	arrInt := is.Attributes().PutEmptySlice("arr_int")
	arrInt.AppendEmpty().SetInt(2)
	arrInt.AppendEmpty().SetInt(3)

	arrFloat := is.Attributes().PutEmptySlice("arr_float")
	arrFloat.AppendEmpty().SetDouble(1.0)
	arrFloat.AppendEmpty().SetDouble(2.0)

	arrBytes := is.Attributes().PutEmptySlice("arr_bytes")
	arrBytes.AppendEmpty().SetEmptyBytes().FromRaw([]byte{1, 2, 3})
	arrBytes.AppendEmpty().SetEmptyBytes().FromRaw([]byte{2, 3, 4})

	return is
}

type instrumentationScopeContext struct {
	is pcommon.InstrumentationScope
}

func (r *instrumentationScopeContext) GetInstrumentationScope() pcommon.InstrumentationScope {
	return r.is
}

func newInstrumentationScopeContext(is pcommon.InstrumentationScope) *instrumentationScopeContext {
	return &instrumentationScopeContext{is: is}
}
