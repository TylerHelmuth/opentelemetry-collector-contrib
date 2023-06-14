// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/internal"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

type ResourceContext interface {
	GetResource() pcommon.Resource
}

func ResourcePathGetSetter[K ResourceContext](path ottl.Path) (ottl.GetSetter[K], error) {
	switch path.Name() {
	case "":
		return accessResource[K](), nil
	case "attributes":
		mapKeys := path.Keys()
		if mapKeys == nil {
			return accessResourceAttributes[K](), nil
		}
		return accessResourceAttributesKey[K](*mapKeys), nil
	case "dropped_attributes_count":
		return accessResourceDroppedAttributesCount[K](), nil
	default:
		return nil, fmt.Errorf("invalid resource path expression %v", path)
	}
}

func accessResource[K ResourceContext]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(ctx context.Context, tCtx K) (interface{}, error) {
			return tCtx.GetResource(), nil
		},
		Setter: func(ctx context.Context, tCtx K, val interface{}) error {
			if newRes, ok := val.(pcommon.Resource); ok {
				newRes.CopyTo(tCtx.GetResource())
			}
			return nil
		},
	}
}

func accessResourceAttributes[K ResourceContext]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(ctx context.Context, tCtx K) (interface{}, error) {
			return tCtx.GetResource().Attributes(), nil
		},
		Setter: func(ctx context.Context, tCtx K, val interface{}) error {
			if attrs, ok := val.(pcommon.Map); ok {
				attrs.CopyTo(tCtx.GetResource().Attributes())
			}
			return nil
		},
	}
}

func accessResourceAttributesKey[K ResourceContext](keys ottl.Key) ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(ctx context.Context, tCtx K) (interface{}, error) {
			return GetMapValue(tCtx.GetResource().Attributes(), keys)
		},
		Setter: func(ctx context.Context, tCtx K, val interface{}) error {
			return SetMapValue(tCtx.GetResource().Attributes(), keys, val)
		},
	}
}

func accessResourceDroppedAttributesCount[K ResourceContext]() ottl.StandardGetSetter[K] {
	return ottl.StandardGetSetter[K]{
		Getter: func(ctx context.Context, tCtx K) (interface{}, error) {
			return int64(tCtx.GetResource().DroppedAttributesCount()), nil
		},
		Setter: func(ctx context.Context, tCtx K, val interface{}) error {
			if i, ok := val.(int64); ok {
				tCtx.GetResource().SetDroppedAttributesCount(uint32(i))
			}
			return nil
		},
	}
}
