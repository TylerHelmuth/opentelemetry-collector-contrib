// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package genainormalizerprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/genainormalizerprocessor"

// This file is the single source of truth for the OTel semantic-conventions
// version this processor targets. Bumping versions: update the import path
// below, verify referenced constants still exist, and bump the README.

import conventions "go.opentelemetry.io/otel/semconv/v1.40.0"

// targetSchemaURL is stamped onto ScopeSpans.schema_url when a mapping fires.
const targetSchemaURL = conventions.SchemaURL

// Target attribute keys used by the mapping tables and the value-transform
// dispatch. Kept as plain strings so they can be used directly as pdata keys.
var (
	targetAttrOperationName     = string(conventions.GenAIOperationNameKey)
	targetAttrUsageInputTokens  = string(conventions.GenAIUsageInputTokensKey)
	targetAttrUsageOutputTokens = string(conventions.GenAIUsageOutputTokensKey)
	targetAttrRequestModel      = string(conventions.GenAIRequestModelKey)
	targetAttrProviderName      = string(conventions.GenAIProviderNameKey)
	targetAttrInputMessages     = string(conventions.GenAIInputMessagesKey)
	targetAttrOutputMessages    = string(conventions.GenAIOutputMessagesKey)
	targetAttrToolName          = string(conventions.GenAIToolNameKey)
	targetAttrToolDescription   = string(conventions.GenAIToolDescriptionKey)
	targetAttrToolCallArguments = string(conventions.GenAIToolCallArgumentsKey)
	targetAttrToolCallID        = string(conventions.GenAIToolCallIDKey)
	targetAttrAgentName         = string(conventions.GenAIAgentNameKey)
	targetAttrConversationID    = string(conventions.GenAIConversationIDKey)
)
