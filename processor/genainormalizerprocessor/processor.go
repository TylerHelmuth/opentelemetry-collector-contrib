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
	// mappings holds the ordered (from, to) pairs for this source. Stored as
	// a slice rather than a map because (a) iteration order is stable, so
	// renames are applied deterministically when two sources map to the same
	// target (e.g. llm.model_name and embedding.model_name both → request.model),
	// and (b) for the small K we care about (≤ ~20), a linear walk is faster
	// and allocation-free compared to building/iterating a hash map.
	mappings []mapping
	// fromSet is the set of from-keys, used only for the cheap probe that
	// short-circuits non-matching spans before we do the per-mapping Get work.
	fromSet         map[string]struct{}
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
		m := builtInMappings(src.Name)
		fromSet := make(map[string]struct{}, len(m))
		for _, mp := range m {
			fromSet[mp.from] = struct{}{}
		}
		p.sources = append(p.sources, sourceNormalizer{
			mappings:        m,
			fromSet:         fromSet,
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
			ss := ilss.At(j)
			spans := ss.Spans()
			scopeWrote := false
			for k := 0; k < spans.Len(); k++ {
				attrs := spans.At(k).Attributes()
				for s := range p.sources {
					if p.sources[s].normalizeAttributes(attrs) {
						scopeWrote = true
					}
				}
			}
			if scopeWrote {
				// We rewrote at least one span's attributes in this scope,
				// so the scope now conforms to our target OTel semconv version.
				ss.SetSchemaUrl(targetSchemaURL)
			}
		}
	}
	return td, nil
}

// normalizeAttributes applies the source's rename rules to attrs. It returns
// true if at least one attribute was written.
//
// Hot path: iterate the mappings slice and probe attrs for each "from" key.
// This is faster than the inverse — Range over attrs + hashtable lookup —
// for match cases because:
//   - It avoids the Range closure heap allocation.
//   - It avoids a temporary "renames to apply" slice.
//
// For non-matching spans (no GenAI attributes), a single Range with hash-set
// lookup is much cheaper than K Gets through the mappings slice, so we gate
// the slow path on a fast probe.
func (sn *sourceNormalizer) normalizeAttributes(attrs pcommon.Map) bool {
	if attrs.Len() == 0 || len(sn.mappings) == 0 {
		return false
	}
	// Probe: do any attrs keys belong to this source? Range over attrs (≤ N
	// hash lookups) is cheaper than blindly running the mappings loop (K Gets
	// over attrs) when the span has nothing to rename — which is most spans
	// flowing through a non-GenAI pipeline.
	mightMatch := false
	attrs.Range(func(k string, _ pcommon.Value) bool {
		if _, ok := sn.fromSet[k]; ok {
			mightMatch = true
			return false
		}
		return true
	})
	if !mightMatch {
		return false
	}

	wrote := false
	for _, m := range sn.mappings {
		val, ok := attrs.Get(m.from)
		if !ok {
			continue
		}
		// GetOrPutEmpty folds the target existence check and the slot
		// allocation into a single scan of attrs. The prior Get(to) + Put*
		// pair scanned twice (Put* itself does an internal Get).
		dest, existed := attrs.GetOrPutEmpty(m.to)
		if existed && !sn.overwrite {
			continue
		}
		// val points into attrs's backing slice. GetOrPutEmpty just above may
		// have reallocated that slice, but val's wrapper holds the *AnyValue
		// pointer by value — the pointer itself was copied during slice grow
		// and still resolves to the original key/value data, so reads remain
		// valid.
		switch val.Type() {
		case pcommon.ValueTypeStr:
			dest.SetStr(transformValue(m.to, val.Str()))
		case pcommon.ValueTypeInt:
			dest.SetInt(val.Int())
		case pcommon.ValueTypeDouble:
			dest.SetDouble(val.Double())
		case pcommon.ValueTypeBool:
			dest.SetBool(val.Bool())
		default:
			// Map / Slice / Bytes: val and dest reference independent
			// AnyValue regions, so CopyTo is safe without a snapshot.
			val.CopyTo(dest)
		}
		wrote = true

		if sn.removeOriginals {
			attrs.Remove(m.from)
		}
	}
	return wrote
}
