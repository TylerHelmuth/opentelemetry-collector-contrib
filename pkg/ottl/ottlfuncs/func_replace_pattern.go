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

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

type ReplacePatternFactory[K any] struct{}

func (f ReplacePatternFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.replacePattern)
}

func (f ReplacePatternFactory[K]) CreateFunction(args []reflect.Value) (ottl.ExprFunc[K], error) {
	return f.replacePattern(args[0].Interface().(ottl.GetSetter[K]), args[1].String(), args[2].String())
}

func (f ReplacePatternFactory[K]) replacePattern(target ottl.GetSetter[K], regexPattern string, replacement string) (ottl.ExprFunc[K], error) {
	compiledPattern, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, fmt.Errorf("the regex pattern supplied to replace_pattern is not a valid pattern: %w", err)
	}
	return func(ctx context.Context, tCtx K) (interface{}, error) {
		originalVal, err := target.Get(ctx, tCtx)
		if err != nil {
			return nil, err
		}
		if originalVal == nil {
			return nil, nil
		}
		if originalValStr, ok := originalVal.(string); ok {
			if compiledPattern.MatchString(originalValStr) {
				updatedStr := compiledPattern.ReplaceAllLiteralString(originalValStr, replacement)
				err = target.Set(ctx, tCtx, updatedStr)
				if err != nil {
					return nil, err
				}
			}
		}
		return nil, nil
	}, nil
}
