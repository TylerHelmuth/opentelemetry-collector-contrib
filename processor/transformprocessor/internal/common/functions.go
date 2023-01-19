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

package common // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor/internal/common"

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlresource"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlscope"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
)

func Functions[K any]() ottl.FunctionFactoryMap[K] {
	return ottl.FunctionFactoryMap[K]{
		"TraceID":              ottlfuncs.TraceIDFactory[K]{},
		"SpanID":               ottlfuncs.SpanIDFactory[K]{},
		"IsMatch":              ottlfuncs.IsMatchFactory[K]{},
		"Concat":               ottlfuncs.ConcatFactory[K]{},
		"Split":                ottlfuncs.SplitFactory[K]{},
		"Int":                  ottlfuncs.IntFactory[K]{},
		"ConvertCase":          ottlfuncs.ConvertCaseFactory[K]{},
		"ParseJSON":            ottlfuncs.ParseJSONFactory[K]{},
		"Substring":            ottlfuncs.SubstringFactory[K]{},
		"keep_keys":            ottlfuncs.KeepKeysFactory[K]{},
		"set":                  ottlfuncs.SetFactory[K]{},
		"truncate_all":         ottlfuncs.TruncateAllFactory[K]{},
		"limit":                ottlfuncs.LimitFactory[K]{},
		"replace_match":        ottlfuncs.ReplaceMatchFactory[K]{},
		"replace_all_matches":  ottlfuncs.ReplaceAllMatchesFactory[K]{},
		"replace_pattern":      ottlfuncs.ReplacePatternFactory[K]{},
		"replace_all_patterns": ottlfuncs.ReplaceAllPatternsFactory[K]{},
		"delete_key":           ottlfuncs.DeleteKeyFactory[K]{},
		"delete_matching_keys": ottlfuncs.DeleteMatchingKeysFactory[K]{},
		"merge_maps":           ottlfuncs.MergeMapsFactory[K]{},
	}
}

func ResourceFunctions() ottl.FunctionFactoryMap[ottlresource.TransformContext] {
	return Functions[ottlresource.TransformContext]()
}

func ScopeFunctions() ottl.FunctionFactoryMap[ottlscope.TransformContext] {
	return Functions[ottlscope.TransformContext]()
}
