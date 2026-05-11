// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package genainormalizerprocessor

import (
	"strconv"
	"testing"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

// populator builds a fresh pcommon.Map representing a single span's attributes
// in a given test shape. Each scenario has its own populator so the benchmark
// inputs are deterministic and independent of the data tables.
type populator func(pcommon.Map)

// runNormalizeBench drives normalizeAttributes against a freshly populated
// attribute map on every iteration. We include the fresh-state setup inside
// the timed region (so b.N converges) and subtract the populate-only baseline
// to recover the normalize cost. StopTimer/StartTimer per-iter is avoided
// because it makes b.N grow unboundedly when normalize is fast relative to
// setup, and adds its own per-call overhead.
func runNormalizeBench(b *testing.B, sn sourceNormalizer, populate populator) {
	b.Helper()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		attrs := pcommon.NewMap()
		populate(attrs)
		sn.normalizeAttributes(attrs)
	}
}

// openInferenceTypical mimics a real OpenInference LLM span: 7 attrs that all
// match openInferenceMappings, scalars only. This is the happy-path shape we
// most want to be fast.
func openInferenceTypical(attrs pcommon.Map) {
	attrs.PutInt("llm.token_count.prompt", 100)
	attrs.PutInt("llm.token_count.completion", 200)
	attrs.PutStr("llm.model_name", "claude-sonnet-4")
	attrs.PutStr("llm.provider", "anthropic")
	attrs.PutStr("openinference.span.kind", "LLM")
	attrs.PutStr("agent.name", "research-agent")
	attrs.PutStr("session.id", "sess-123")
}

// openInferenceWithBigSlice adds a large input_messages slice to the typical
// shape. Exercises the non-scalar snapshot path.
func openInferenceWithBigSlice(attrs pcommon.Map) {
	openInferenceTypical(attrs)
	s := attrs.PutEmptySlice("llm.input_messages")
	for i := 0; i < 8; i++ {
		m := s.AppendEmpty().SetEmptyMap()
		m.PutStr("role", "user")
		m.PutStr("content", "what is the capital of france and please be detailed about it")
	}
}

// noMatchSpan is a non-GenAI span (e.g. an HTTP server span). normalize must
// pay the Range cost but make no renames. This is the cost the processor adds
// to every passthrough span.
func noMatchSpan(attrs pcommon.Map) {
	attrs.PutStr("http.method", "GET")
	attrs.PutStr("http.url", "https://example.com/api/v1/things")
	attrs.PutInt("http.status_code", 200)
	attrs.PutStr("net.peer.name", "example.com")
	attrs.PutInt("net.peer.port", 443)
	attrs.PutStr("user_agent.original", "Mozilla/5.0")
	attrs.PutStr("server.address", "example.com")
}

// fewMatchesManyAttrs simulates a large span (20 attrs) where only 3 keys
// match. Measures the cost of Range over many non-matching keys plus a few
// renames.
func fewMatchesManyAttrs(attrs pcommon.Map) {
	// 3 matching keys.
	attrs.PutStr("llm.model_name", "claude-sonnet-4")
	attrs.PutInt("llm.token_count.prompt", 50)
	attrs.PutStr("llm.provider", "anthropic")
	// 17 non-matching keys.
	for i := 0; i < 17; i++ {
		attrs.PutStr("noise.attr."+strconv.Itoa(i), "v"+strconv.Itoa(i))
	}
}

// overwriteFalseAllCollide hits the overwrite=false skip path on every rename
// (every target key already exists). Stresses the per-rename target-exists
// lookup.
func overwriteFalseAllCollide(attrs pcommon.Map) {
	openInferenceTypical(attrs)
	// Pre-populate every target.
	attrs.PutInt(targetAttrUsageInputTokens, 1)
	attrs.PutInt(targetAttrUsageOutputTokens, 1)
	attrs.PutStr(targetAttrRequestModel, "x")
	attrs.PutStr(targetAttrProviderName, "x")
	attrs.PutStr(targetAttrOperationName, "x")
	attrs.PutStr(targetAttrAgentName, "x")
	attrs.PutStr(targetAttrConversationID, "x")
}

func benchSource(removeOriginals, overwrite bool) sourceNormalizer {
	m := builtInMappings(SourceOpenInference)
	fromSet := make(map[string]struct{}, len(m))
	for _, mp := range m {
		fromSet[mp.from] = struct{}{}
	}
	return sourceNormalizer{
		mappings:        m,
		fromSet:         fromSet,
		removeOriginals: removeOriginals,
		overwrite:       overwrite,
	}
}

func BenchmarkNormalizeAttributes(b *testing.B) {
	cases := []struct {
		name            string
		populate        populator
		removeOriginals bool
		overwrite       bool
	}{
		{"openinference_typical_remove", openInferenceTypical, true, false},
		{"openinference_typical_keep", openInferenceTypical, false, false},
		{"openinference_with_big_slice", openInferenceWithBigSlice, true, false},
		{"no_match_passthrough", noMatchSpan, true, false},
		{"few_matches_many_attrs", fewMatchesManyAttrs, true, false},
		{"overwrite_false_all_collide", overwriteFalseAllCollide, false, false},
	}
	for _, c := range cases {
		sn := benchSource(c.removeOriginals, c.overwrite)
		b.Run(c.name, func(b *testing.B) {
			runNormalizeBench(b, sn, c.populate)
		})
	}
}

// BenchmarkPopulateOnly establishes the noise floor for runNormalizeBench:
// the cost of building a fresh attribute map plus StopTimer/StartTimer
// bookkeeping, with no normalize work. Subtract from the cases above to
// approximate the pure normalize cost.
func BenchmarkPopulateOnly(b *testing.B) {
	cases := []struct {
		name     string
		populate populator
	}{
		{"openinference_typical", openInferenceTypical},
		{"openinference_with_big_slice", openInferenceWithBigSlice},
		{"no_match_passthrough", noMatchSpan},
		{"few_matches_many_attrs", fewMatchesManyAttrs},
		{"overwrite_false_all_collide", overwriteFalseAllCollide},
	}
	noopSrc := sourceNormalizer{}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			runNormalizeBench(b, noopSrc, c.populate)
		})
	}
}
