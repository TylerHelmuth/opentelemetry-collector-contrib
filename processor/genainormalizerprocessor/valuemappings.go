// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package genainormalizerprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/genainormalizerprocessor"

import "strings"

// operationNameValues maps OpenInference openinference.span.kind values to
// OTel GenAI operation names. Keys are lowercased; transformValue lowercases
// the input before lookup. Target values follow the enum from
// go.opentelemetry.io/otel/semconv (see semconv.go).
//
// Reference: https://github.com/Arize-ai/openinference/blob/main/spec/semantic_conventions.md#span-kinds
var operationNameValues = map[string]string{
	"llm":       "chat",            // GenAIOperationNameChat
	"embedding": "embeddings",      // GenAIOperationNameEmbeddings
	"chain":     "invoke_agent",    // GenAIOperationNameInvokeAgent
	"retriever": "retrieval",       // GenAIOperationNameRetrieval
	"reranker":  "retrieval",       // GenAIOperationNameRetrieval
	"tool":      "execute_tool",    // GenAIOperationNameExecuteTool
	"agent":     "invoke_agent",    // GenAIOperationNameInvokeAgent
	"prompt":    "text_completion", // GenAIOperationNameTextCompletion
}

// transformValue applies value mapping for a given target attribute key.
// Returns the transformed value, or the original if no mapping exists.
func transformValue(targetKey, value string) string {
	if targetKey != targetAttrOperationName {
		return value
	}
	if mapped, ok := operationNameValues[strings.ToLower(value)]; ok {
		return mapped
	}
	return value
}
