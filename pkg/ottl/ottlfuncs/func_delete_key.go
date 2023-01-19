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

package ottlfuncs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"

import (
	"context"
	"reflect"

	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

type DeleteKeyFactory[K any] struct{}

func (f DeleteKeyFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.deleteKey)
}

func (f DeleteKeyFactory[K]) CreateFunction(args []reflect.Value) (ottl.ExprFunc[K], error) {
	return f.deleteKey(args[0].Interface().(ottl.Getter[K]), args[1].String())
}

func (f DeleteKeyFactory[K]) deleteKey(target ottl.Getter[K], key string) (ottl.ExprFunc[K], error) {
	return func(ctx context.Context, tCtx K) (interface{}, error) {
		val, err := target.Get(ctx, tCtx)
		if err != nil {
			return nil, err
		}
		if val == nil {
			return nil, nil
		}

		if attrs, ok := val.(pcommon.Map); ok {
			attrs.Remove(key)
		}
		return nil, nil
	}, nil
}
