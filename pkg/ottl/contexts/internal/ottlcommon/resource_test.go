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

func TestResourcePathGetSetter(t *testing.T) {
	refResource := createResource()

	newAttrs := pcommon.NewMap()
	newAttrs.PutStr("hello", "world")

	tests := []struct {
		name     string
		path     ottl.Path
		orig     interface{}
		newVal   interface{}
		modified func(resource pcommon.Resource)
	}{
		{
			name:   "resource",
			path:   ottl.Path{},
			orig:   refResource,
			newVal: pcommon.NewResource(),
			modified: func(resource pcommon.Resource) {
				pcommon.NewResource().CopyTo(resource)
			},
		},
		{
			name: "attributes",
			path: ottl.Path{
				Fields: []string{
					"attributes",
				},
			},
			orig:   refResource.Attributes(),
			newVal: newAttrs,
			modified: func(resource pcommon.Resource) {
				newAttrs.CopyTo(resource.Attributes())
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
			modified: func(resource pcommon.Resource) {
				resource.SetDroppedAttributesCount(20)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessor, err := ResourcePathGetSetter[*resourceContext](tt.path)
			assert.NoError(t, err)

			resource := createResource()

			got, err := accessor.Get(context.Background(), newResourceContext(resource))
			assert.Nil(t, err)
			assert.Equal(t, tt.orig, got)

			err = accessor.Set(context.Background(), newResourceContext(resource), tt.newVal)
			assert.Nil(t, err)

			expectedResource := createResource()
			tt.modified(expectedResource)

			assert.Equal(t, expectedResource, resource)
		})
	}
}

func createResource() pcommon.Resource {
	resource := pcommon.NewResource()
	resource.Attributes().PutStr("str", "val")
	resource.Attributes().PutBool("bool", true)
	resource.Attributes().PutInt("int", 10)
	resource.Attributes().PutDouble("double", 1.2)
	resource.Attributes().PutEmptyBytes("bytes").FromRaw([]byte{1, 3, 2})

	resource.Attributes().PutEmptySlice("arr_empty")

	arrStr := resource.Attributes().PutEmptySlice("arr_str")
	arrStr.AppendEmpty().SetStr("one")
	arrStr.AppendEmpty().SetStr("two")

	arrBool := resource.Attributes().PutEmptySlice("arr_bool")
	arrBool.AppendEmpty().SetBool(true)
	arrBool.AppendEmpty().SetBool(false)

	arrInt := resource.Attributes().PutEmptySlice("arr_int")
	arrInt.AppendEmpty().SetInt(2)
	arrInt.AppendEmpty().SetInt(3)

	arrFloat := resource.Attributes().PutEmptySlice("arr_float")
	arrFloat.AppendEmpty().SetDouble(1.0)
	arrFloat.AppendEmpty().SetDouble(2.0)

	arrBytes := resource.Attributes().PutEmptySlice("arr_bytes")
	arrBytes.AppendEmpty().SetEmptyBytes().FromRaw([]byte{1, 2, 3})
	arrBytes.AppendEmpty().SetEmptyBytes().FromRaw([]byte{2, 3, 4})

	resource.SetDroppedAttributesCount(10)

	il := pcommon.NewInstrumentationScope()
	il.SetName("library")
	il.SetVersion("version")

	return resource
}

type resourceContext struct {
	resource pcommon.Resource
}

func (r *resourceContext) GetResource() pcommon.Resource {
	return r.resource
}

func newResourceContext(resource pcommon.Resource) *resourceContext {
	return &resourceContext{resource: resource}
}
