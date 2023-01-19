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

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

type SetFactory[K any] struct{}

func (f SetFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.set)
}

func (f SetFactory[K]) CreateFunction(args []reflect.Value) (ottl.ExprFunc[K], error) {
	return f.set(args[0].Interface().(ottl.Setter[K]), args[1].Interface().(ottl.Getter[K]))
}
func (f SetFactory[K]) set(target ottl.Setter[K], value ottl.Getter[K]) (ottl.ExprFunc[K], error) {
	return func(ctx context.Context, tCtx K) (interface{}, error) {
		val, err := value.Get(ctx, tCtx)
		if err != nil {
			return nil, err
		}

		// No fields currently support `null` as a valid type.
		if val != nil {
			err = target.Set(ctx, tCtx, val)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	}, nil
}
