// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package genainormalizerprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformValue(t *testing.T) {
	tests := []struct {
		name      string
		targetKey string
		value     string
		want      string
	}{
		{
			name:      "non-operation-name target passes through unchanged",
			targetKey: "gen_ai.request.model",
			value:     "LLM",
			want:      "LLM",
		},
		{
			// operationNameValues is empty in this PR; dispatch still runs and
			// falls back to the input value.
			name:      "operation-name target with unmapped value passes through",
			targetKey: targetOperationName,
			value:     "anything",
			want:      "anything",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, transformValue(tt.targetKey, tt.value))
		})
	}
}
