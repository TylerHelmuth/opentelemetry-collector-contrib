// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package ottlfuncs // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"

import (
	"context"
	"errors"
	"reflect"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

type IsEmptyArguments[K any] struct {
	Target ottl.Getter[K]
}

func NewIsEmptyFactory[K any]() ottl.Factory[K] {
	return ottl.NewFactory("IsEmpty", &IsEmptyArguments[K]{}, createIsEmptyFunction[K])
}

func createIsEmptyFunction[K any](_ ottl.FunctionContext, oArgs ottl.Arguments) (ottl.ExprFunc[K], error) {
	args, ok := oArgs.(*IsEmptyArguments[K])

	if !ok {
		return nil, errors.New("IsEmptyFactory args must be of type *IsEmptyArguments[K]")
	}

	return isEmpty(args.Target), nil
}

func isEmpty[K any](target ottl.Getter[K]) ottl.ExprFunc[K] {
	return func(ctx context.Context, tCtx K) (any, error) {
		val, err := target.Get(ctx, tCtx)
		if err != nil {
			return false, err
		}
		if val == nil {
			return true, nil
		}
		switch v := val.(type) {
		case pcommon.Value:
			return v.Type() == pcommon.ValueTypeEmpty, nil
		case pcommon.Map:
			return v.Len() == 0, nil
		case pcommon.Slice:
			return v.Len() == 0, nil
		default:
			// Handle native collection types (slices such as []any, []byte or
			// []string, and maps such as map[string]any) by length. Fixed-size
			// arrays have a constant length, so for them - and any other type -
			// emptiness means the zero value. This treats an unset (all-zero)
			// pcommon.TraceID/SpanID as empty, matching pdata's own IsEmpty.
			rv := reflect.ValueOf(v)
			switch rv.Kind() {
			case reflect.Slice, reflect.Map:
				return rv.Len() == 0, nil
			default:
				return rv.IsZero(), nil
			}
		}
	}
}
