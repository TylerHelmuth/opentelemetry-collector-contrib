// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package genainormalizerprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/genainormalizerprocessor"

import (
	"context"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// sourceNormalizer holds per-source state used during normalization.
type sourceNormalizer struct {
	lookupTable     map[string]string
	removeOriginals bool
	overwrite       bool
}

// genaiNormalizerProcessor normalizes span attributes for each configured source.
type genaiNormalizerProcessor struct {
	sources []sourceNormalizer
}

// newGenaiNormalizerProcessor builds a processor from a validated Config.
// Sources are applied in the order specified in the configuration.
func newGenaiNormalizerProcessor(cfg *Config) *genaiNormalizerProcessor {
	p := &genaiNormalizerProcessor{sources: make([]sourceNormalizer, 0, len(cfg.Sources))}
	for _, src := range cfg.Sources {
		p.sources = append(p.sources, sourceNormalizer{
			lookupTable:     buildLookupTable(src.Name),
			removeOriginals: src.RemoveOriginals,
			overwrite:       src.Overwrite,
		})
	}
	return p
}

func (p *genaiNormalizerProcessor) processTraces(_ context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		ilss := rss.At(i).ScopeSpans()
		for j := 0; j < ilss.Len(); j++ {
			spans := ilss.At(j).Spans()
			for k := 0; k < spans.Len(); k++ {
				attrs := spans.At(k).Attributes()
				for s := range p.sources {
					p.sources[s].normalizeAttributes(attrs)
				}
			}
		}
	}
	return td, nil
}

func (sn *sourceNormalizer) normalizeAttributes(attrs pcommon.Map) {
	type rename struct {
		from string
		to   string
	}
	var renames []rename

	attrs.Range(func(k string, _ pcommon.Value) bool {
		if target, ok := sn.lookupTable[k]; ok {
			renames = append(renames, rename{k, target})
		}
		return true
	})

	if len(renames) == 0 {
		return
	}

	for _, r := range renames {
		val, ok := attrs.Get(r.from)
		if !ok {
			continue
		}
		if _, exists := attrs.Get(r.to); exists && !sn.overwrite {
			continue
		}

		// Read scalar values up front so they remain valid across subsequent
		// Put* calls that may reallocate the backing slice.
		switch val.Type() {
		case pcommon.ValueTypeStr:
			attrs.PutStr(r.to, transformValue(r.to, val.Str()))
		case pcommon.ValueTypeInt:
			attrs.PutInt(r.to, val.Int())
		case pcommon.ValueTypeDouble:
			attrs.PutDouble(r.to, val.Double())
		case pcommon.ValueTypeBool:
			attrs.PutBool(r.to, val.Bool())
		default:
			// Map / Slice / Bytes: snapshot before writing so the source read
			// is independent of any reallocation the destination write causes.
			snap := pcommon.NewValueEmpty()
			val.CopyTo(snap)
			snap.CopyTo(attrs.PutEmpty(r.to))
		}

		if sn.removeOriginals {
			attrs.Remove(r.from)
		}
	}
}
