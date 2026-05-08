// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package genainormalizerprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/genainormalizerprocessor"

import (
	"strings"

	conventions "go.opentelemetry.io/otel/semconv/v1.40.0"
)

// targetOperationName is the OTel GenAI operation-name attribute key.
// Sourced from go.opentelemetry.io/otel/semconv for alignment with upstream.
var targetOperationName = string(conventions.GenAIOperationNameKey)

// operationNameValues maps OpenInference openinference.span.kind values to
// OTel GenAI operation names. Keys are lowercased; transformValue lowercases
// the input before lookup. Target values follow the enum from
// go.opentelemetry.io/otel/semconv/v1.40.0.GenAIOperationName*.
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
	if targetKey != targetOperationName {
		return value
	}
	if mapped, ok := operationNameValues[strings.ToLower(value)]; ok {
		return mapped
	}
	return value
}
