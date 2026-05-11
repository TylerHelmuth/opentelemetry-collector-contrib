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

// openInferenceMappings is the OpenInference attribute-rename table.
// Target keys come from semconv.go (single source of truth for the targeted
// OTel semconv version).
// Reference: https://github.com/Arize-ai/openinference/blob/main/spec/semantic_conventions.md
var openInferenceMappings = []mapping{
	// Token usage
	{from: "llm.token_count.prompt", to: targetAttrUsageInputTokens},
	{from: "llm.token_count.completion", to: targetAttrUsageOutputTokens},

	// Model & provider
	{from: "llm.model_name", to: targetAttrRequestModel},
	{from: "llm.provider", to: targetAttrProviderName},

	// Input/output content
	{from: "llm.input_messages", to: targetAttrInputMessages},
	{from: "llm.output_messages", to: targetAttrOutputMessages},

	// Embeddings
	{from: "embedding.model_name", to: targetAttrRequestModel},

	// Tool
	{from: "tool.name", to: targetAttrToolName},
	{from: "tool.description", to: targetAttrToolDescription},
	{from: "tool_call.function.arguments", to: targetAttrToolCallArguments},
	{from: "tool_call.id", to: targetAttrToolCallID},

	// Reranker model (mutually exclusive with llm/embedding model)
	{from: "reranker.model_name", to: targetAttrRequestModel},

	// Agent & session
	{from: "agent.name", to: targetAttrAgentName},
	{from: "session.id", to: targetAttrConversationID},

	// Span kind -> operation name (value mapping handled in valuemappings.go)
	{from: "openinference.span.kind", to: targetAttrOperationName},
}
