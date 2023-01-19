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
	"fmt"
	"reflect"
	"regexp"

	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

type DeleteMatchingKeysFactory[K any] struct{}

func (f DeleteMatchingKeysFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.deleteMatchingKeys)
}

func (f DeleteMatchingKeysFactory[K]) CreateFunction(args []reflect.Value) (ottl.ExprFunc[K], error) {
	return f.deleteMatchingKeys(args[0].Interface().(ottl.Getter[K]), args[1].String())
}

func (f DeleteMatchingKeysFactory[K]) deleteMatchingKeys(target ottl.Getter[K], pattern string) (ottl.ExprFunc[K], error) {
	compiledPattern, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("the regex pattern supplied to delete_matching_keys is not a valid pattern: %w", err)
	}
	return func(ctx context.Context, tCtx K) (interface{}, error) {
		val, err := target.Get(ctx, tCtx)
		if err != nil {
			return nil, err
		}
		if val == nil {
			return nil, nil
		}

		if attrs, ok := val.(pcommon.Map); ok {
			attrs.RemoveIf(func(key string, _ pcommon.Value) bool {
				return compiledPattern.MatchString(key)
			})
		}
		return nil, nil
	}, nil
}
