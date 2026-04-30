// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package genainormalizerprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/genainormalizerprocessor"

// mapping defines a source attribute key and its normalized target.
type mapping struct {
	from string
	to   string
}

// builtInMappings returns the built-in mappings for a given source.
func builtInMappings(name SourceName) []mapping {
	if name == SourceOpenInference {
		return openInferenceMappings
	}
	return nil
}

// buildLookupTable assembles a per-source lookup from built-in mappings.
func buildLookupTable(name SourceName) map[string]string {
	table := make(map[string]string)
	for _, m := range builtInMappings(name) {
		table[m.from] = m.to
	}
	return table
}

// openInferenceMappings is the OpenInference attribute-rename table.
// Populated in a follow-up PR.
// Reference: https://github.com/Arize-ai/openinference/blob/main/spec/semantic_conventions.md
var openInferenceMappings = []mapping{}
