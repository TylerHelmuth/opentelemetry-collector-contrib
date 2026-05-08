// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package genainormalizerprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/genainormalizerprocessor"

import (
	conventions "go.opentelemetry.io/otel/semconv/v1.40.0"
)

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
// Reference: https://github.com/Arize-ai/openinference/blob/main/spec/semantic_conventions.md
var openInferenceMappings = []mapping{
	// Token usage
	{from: "llm.token_count.prompt", to: string(conventions.GenAIUsageInputTokensKey)},
	{from: "llm.token_count.completion", to: string(conventions.GenAIUsageOutputTokensKey)},

	// Model & provider
	{from: "llm.model_name", to: string(conventions.GenAIRequestModelKey)},
	{from: "llm.provider", to: string(conventions.GenAIProviderNameKey)},

	// Input/output content
	{from: "llm.input_messages", to: string(conventions.GenAIInputMessagesKey)},
	{from: "llm.output_messages", to: string(conventions.GenAIOutputMessagesKey)},

	// Embeddings
	{from: "embedding.model_name", to: string(conventions.GenAIRequestModelKey)},

	// Tool
	{from: "tool.name", to: string(conventions.GenAIToolNameKey)},
	{from: "tool.description", to: string(conventions.GenAIToolDescriptionKey)},
	{from: "tool_call.function.arguments", to: string(conventions.GenAIToolCallArgumentsKey)},
	{from: "tool_call.id", to: string(conventions.GenAIToolCallIDKey)},

	// Reranker model (mutually exclusive with llm/embedding model)
	{from: "reranker.model_name", to: string(conventions.GenAIRequestModelKey)},

	// Agent & session
	{from: "agent.name", to: string(conventions.GenAIAgentNameKey)},
	{from: "session.id", to: string(conventions.GenAIConversationIDKey)},

	// Span kind -> operation name (value mapping handled in valuemappings.go)
	{from: "openinference.span.kind", to: string(conventions.GenAIOperationNameKey)},
}
