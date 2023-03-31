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

package ottlscope

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

func Test_newPathGetSetter(t *testing.T) {
	refIS, refResource := createTelemetry()

	newAttrs := pcommon.NewMap()
	newAttrs.PutStr("hello", "world")

	newCache := pcommon.NewMap()
	newCache.PutStr("temp", "value")

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
		modified func(is pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map)
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
			modified: func(il pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				newCache.CopyTo(cache)
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
			modified: func(is pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				newAttrs.CopyTo(is.Attributes())
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
			modified: func(is pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				is.SetDroppedAttributesCount(20)
			},
		},
		{
			name: "name",
			path: ottl.Path{
				Fields: []string{
					"name",
				},
			},
			orig:   refIS.Name(),
			newVal: "newname",
			modified: func(is pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				is.SetName("newname")
			},
		},
		{
			name: "version",
			path: ottl.Path{
				Fields: []string{
					"version",
				},
			},
			orig:   refIS.Version(),
			newVal: "next",
			modified: func(is pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				is.SetVersion("next")
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
			modified: func(is pcommon.InstrumentationScope, resource pcommon.Resource, cache pcommon.Map) {
				pcommon.NewResource().CopyTo(resource)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessor, err := newPathGetSetter(tt.path)
			assert.NoError(t, err)

			il, resource := createTelemetry()

			tCtx := NewTransformContext(il, resource)
			got, err := accessor.Get(context.Background(), tCtx)
			assert.Nil(t, err)
			assert.Equal(t, tt.orig, got)

			err = accessor.Set(context.Background(), tCtx, tt.newVal)
			assert.Nil(t, err)

			exIl, exRes := createTelemetry()
			exCache := pcommon.NewMap()
			tt.modified(exIl, exRes, exCache)

			assert.Equal(t, exIl, il)
			assert.Equal(t, exRes, resource)
			assert.Equal(t, exCache, tCtx.getCache())
		})
	}
}

func createTelemetry() (pcommon.InstrumentationScope, pcommon.Resource) {
	is := pcommon.NewInstrumentationScope()
	is.SetName("library")
	is.SetVersion("version")
	is.SetDroppedAttributesCount(10)

	is.Attributes().PutStr("str", "val")
	is.Attributes().PutBool("bool", true)
	is.Attributes().PutInt("int", 10)
	is.Attributes().PutDouble("double", 1.2)
	is.Attributes().PutEmptyBytes("bytes").FromRaw([]byte{1, 3, 2})

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

	pMap := is.Attributes().PutEmptyMap("pMap")
	pMap.PutStr("original", "map")

	m := is.Attributes().PutEmptyMap("map")
	m.PutStr("original", "map")

	resource := pcommon.NewResource()
	is.Attributes().CopyTo(resource.Attributes())

	return is, resource
}
