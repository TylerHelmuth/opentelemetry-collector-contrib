// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ottl // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"

import (
	"github.com/alecthomas/participle/v2"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"testing"
)

func Test_evaluateMathExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "simple subtraction",
			input:    "(1000 - 600)",
			expected: 400,
		},
		{
			name:     "simple division",
			input:    "(1 / 1)",
			expected: 1,
		},
		{
			name:     "subtraction and addition",
			input:    "(1000 - 600 + 1)",
			expected: 401,
		},
		{
			name:     "order of operations",
			input:    "(10 - 6 * 2 + 2)",
			expected: 0,
		},
		{
			name:     "parentheses",
			input:    "(30 - 6 * (2 + 2))",
			expected: 6,
		},
		{
			name:     "complex",
			input:    "((4 * 2) + 1 + 1 - 3 / 3 + ( 2 + 1 - (6 / 3)))",
			expected: 10,
		},
		{
			name:     "floats",
			input:    "(.5 + 2.6)",
			expected: 3.1,
		},
		{
			name:     "complex floats",
			input:    "((.5 * 4.0) / .1 + 3.9)",
			expected: 23.9,
		},
	}

	functions := map[string]interface{}{"hello": hello[interface{}]}

	p := NewParser[any](
		functions,
		testParsePath,
		testParseEnum,
		component.TelemetrySettings{},
	)

	lex := buildLexer()
	mathParser, _ := participle.Build[value](
		participle.Lexer(lex),
		participle.Unquote("String"),
		participle.Elide("whitespace"),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := mathParser.ParseString("", tt.input)
			assert.NoError(t, err)

			getter, err := p.evaluateMathExpression(parsed.MathExpression)
			assert.NoError(t, err)

			result, err := getter.Get(nil)
			assert.NoError(t, err)

			assert.EqualValues(t, tt.expected, result)
		})
	}
}

func Test_evaluateMathExpression_error(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "mixing int and float",
			input: "(1 - 1.0)",
		},
	}

	functions := map[string]interface{}{"hello": hello[interface{}]}

	p := NewParser[any](
		functions,
		testParsePath,
		testParseEnum,
		component.TelemetrySettings{},
	)

	lex := buildLexer()
	mathParser, _ := participle.Build[value](
		participle.Lexer(lex),
		participle.Unquote("String"),
		participle.Elide("whitespace"),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := mathParser.ParseString("", tt.input)
			assert.NoError(t, err)

			getter, err := p.evaluateMathExpression(parsed.MathExpression)
			assert.NoError(t, err)

			result, err := getter.Get(nil)
			assert.Nil(t, result)
			assert.Error(t, err)
		})
	}
}
