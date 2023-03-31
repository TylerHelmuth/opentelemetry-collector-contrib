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

package ottlcommon // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/internal/ottlcommon"

import (
	"fmt"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func GetMapValue(m pcommon.Map, keys []ottl.Key) (interface{}, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("cannot get map value without key")
	}
	if keys[0].Map == nil {
		return nil, fmt.Errorf("non-string indexing is not supported")
	}

	val, ok := m.Get(*keys[0].Map)
	if !ok {
		return nil, nil
	}
	for i := 1; i < len(keys); i++ {
		if keys[i].Map == nil {
			return nil, fmt.Errorf("non-string indexing is not supported")
		}
		switch val.Type() {
		case pcommon.ValueTypeMap:
			val, ok = val.Map().Get(*keys[i].Map)
			if !ok {
				return nil, fmt.Errorf("key not found in map")
			}
		default:
			return nil, fmt.Errorf("type %v does not support string indexing", val.Type())
		}
	}
	return GetValue(val), nil
}

func SetMapValue(m pcommon.Map, keys []ottl.Key, val interface{}) error {
	if len(keys) == 0 {
		return fmt.Errorf("cannot set map value without key")
	}

	var value pcommon.Value
	switch val.(type) {
	case []string, []bool, []int64, []float64, [][]byte, []any:
		value = pcommon.NewValueSlice()
	default:
		value = pcommon.NewValueEmpty()
	}
	SetValue(value, val)

	currentMap := m
	lastKeyIndex := len(keys) - 1
	for i := 0; i < lastKeyIndex; i++ {
		if keys[i].Map == nil {
			return fmt.Errorf("non-string indexing is not supported")
		}
		v, ok := currentMap.Get(*keys[i].Map)
		if !ok || v.Type() != pcommon.ValueTypeMap {
			currentMap = currentMap.PutEmptyMap(*keys[i].Map)
		} else {
			currentMap = v.Map()
		}
	}
	if keys[lastKeyIndex].Map == nil {
		return fmt.Errorf("non-string indexing is not supported")
	}
	value.CopyTo(currentMap.PutEmpty(*keys[lastKeyIndex].Map))
	return nil
}
