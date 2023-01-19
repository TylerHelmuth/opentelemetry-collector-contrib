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

package ottl

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottltest"
)

func Test_NewFunctionCall_invalid(t *testing.T) {
	functions := make(map[string]FunctionFactory[any])
	functions["testing_error"] = functionThatHasAnErrorFactory[any]{}
	functions["testing_getsetter"] = functionWithGetSetterFactory[any]{}
	functions["testing_getter"] = functionWithGetterFactory[any]{}
	functions["testing_multiple_args"] = functionWithMultipleArgsFactory[any]{}
	functions["testing_string"] = functionWithStringFactory[any]{}
	functions["testing_string_slice"] = functionWithStringSliceFactory[any]{}
	functions["testing_byte_slice"] = functionWithByteSliceFactory[any]{}
	functions["testing_enum"] = functionWithEnumFactory[any]{}
	functions["testing_telemetry_settings_first"] = functionWithTelemetrySettingsFirstFactory[any]{}

	p := NewParser(
		functions,
		testParsePath,
		testParseEnum,
		component.TelemetrySettings{},
	)

	tests := []struct {
		name string
		inv  invocation
	}{
		{
			name: "unknown function",
			inv: invocation{
				Function:  "unknownfunc",
				Arguments: []value{},
			},
		},
		{
			name: "not accessor",
			inv: invocation{
				Function: "testing_getsetter",
				Arguments: []value{
					{
						String: ottltest.Strp("not path"),
					},
				},
			},
		},
		{
			name: "not reader (invalid function)",
			inv: invocation{
				Function: "testing_getter",
				Arguments: []value{
					{
						Literal: &mathExprLiteral{
							Converter: &converter{
								Function: "Unknownfunc",
							},
						},
					},
				},
			},
		},
		{
			name: "not enough args",
			inv: invocation{
				Function: "testing_multiple_args",
				Arguments: []value{
					{
						Literal: &mathExprLiteral{
							Path: &Path{
								Fields: []Field{
									{
										Name: "name",
									},
								},
							},
						},
					},
					{
						String: ottltest.Strp("test"),
					},
				},
			},
		},
		{
			name: "too many args",
			inv: invocation{
				Function: "testing_multiple_args",
				Arguments: []value{
					{
						Literal: &mathExprLiteral{
							Path: &Path{
								Fields: []Field{
									{
										Name: "name",
									},
								},
							},
						},
					},
					{
						String: ottltest.Strp("test"),
					},
					{
						String: ottltest.Strp("test"),
					},
				},
			},
		},
		{
			name: "not enough args with telemetrySettings",
			inv: invocation{
				Function: "testing_telemetry_settings_first",
				Arguments: []value{
					{
						String: ottltest.Strp("test"),
					},
					{
						List: &list{
							Values: []value{
								{
									String: ottltest.Strp("test"),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "too many args with telemetrySettings",
			inv: invocation{
				Function: "testing_telemetry_settings_first",
				Arguments: []value{
					{
						String: ottltest.Strp("test"),
					},
					{
						List: &list{
							Values: []value{
								{
									String: ottltest.Strp("test"),
								},
							},
						},
					},
					{
						Literal: &mathExprLiteral{
							Int: ottltest.Intp(10),
						},
					},
					{
						Literal: &mathExprLiteral{
							Int: ottltest.Intp(10),
						},
					},
				},
			},
		},
		{
			name: "not matching arg type",
			inv: invocation{
				Function: "testing_string",
				Arguments: []value{
					{
						Literal: &mathExprLiteral{
							Int: ottltest.Intp(10),
						},
					},
				},
			},
		},
		{
			name: "not matching arg type when byte slice",
			inv: invocation{
				Function: "testing_byte_slice",
				Arguments: []value{
					{
						String: ottltest.Strp("test"),
					},
					{
						String: ottltest.Strp("test"),
					},
					{
						String: ottltest.Strp("test"),
					},
				},
			},
		},
		{
			name: "mismatching slice element type",
			inv: invocation{
				Function: "testing_string_slice",
				Arguments: []value{
					{
						List: &list{
							Values: []value{
								{
									String: ottltest.Strp("test"),
								},
								{
									Literal: &mathExprLiteral{
										Int: ottltest.Intp(10),
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "mismatching slice argument type",
			inv: invocation{
				Function: "testing_string_slice",
				Arguments: []value{
					{
						String: ottltest.Strp("test"),
					},
				},
			},
		},
		{
			name: "function call returns error",
			inv: invocation{
				Function: "testing_error",
			},
		},
		{
			name: "Enum not found",
			inv: invocation{
				Function: "testing_enum",
				Arguments: []value{
					{
						Enum: (*EnumSymbol)(ottltest.Strp("SYMBOL_NOT_FOUND")),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := p.newFunctionCall(tt.inv)
			assert.Error(t, err)
		})
	}
}

func Test_NewFunctionCall(t *testing.T) {
	p := NewParser(
		defaultFunctionsForTests(),
		testParsePath,
		testParseEnum,
		component.TelemetrySettings{},
	)

	tests := []struct {
		name string
		inv  invocation
		want any
	}{
		{
			name: "empty slice arg",
			inv: invocation{
				Function: "testing_string_slice",
				Arguments: []value{
					{
						List: &list{
							Values: []value{},
						},
					},
				},
			},
			want: 0,
		},
		{
			name: "string slice arg",
			inv: invocation{
				Function: "testing_string_slice",
				Arguments: []value{
					{
						List: &list{
							Values: []value{
								{
									String: ottltest.Strp("test"),
								},
								{
									String: ottltest.Strp("test"),
								},
								{
									String: ottltest.Strp("test"),
								},
							},
						},
					},
				},
			},
			want: 3,
		},
		{
			name: "float slice arg",
			inv: invocation{
				Function: "testing_float_slice",
				Arguments: []value{
					{
						List: &list{
							Values: []value{
								{
									Literal: &mathExprLiteral{
										Float: ottltest.Floatp(1.1),
									},
								},
								{
									Literal: &mathExprLiteral{
										Float: ottltest.Floatp(1.2),
									},
								},
								{
									Literal: &mathExprLiteral{
										Float: ottltest.Floatp(1.3),
									},
								},
							},
						},
					},
				},
			},
			want: 3,
		},
		{
			name: "int slice arg",
			inv: invocation{
				Function: "testing_int_slice",
				Arguments: []value{
					{
						List: &list{
							Values: []value{
								{
									Literal: &mathExprLiteral{
										Int: ottltest.Intp(1),
									},
								},
								{
									Literal: &mathExprLiteral{
										Int: ottltest.Intp(1),
									},
								},
								{
									Literal: &mathExprLiteral{
										Int: ottltest.Intp(1),
									},
								},
							},
						},
					},
				},
			},
			want: 3,
		},
		{
			name: "getter slice arg",
			inv: invocation{
				Function: "testing_getter_slice",
				Arguments: []value{
					{
						List: &list{
							Values: []value{
								{
									Literal: &mathExprLiteral{
										Path: &Path{
											Fields: []Field{
												{
													Name: "name",
												},
											},
										},
									},
								},
								{
									String: ottltest.Strp("test"),
								},
								{
									Literal: &mathExprLiteral{
										Int: ottltest.Intp(1),
									},
								},
								{
									Literal: &mathExprLiteral{
										Float: ottltest.Floatp(1.1),
									},
								},
								{
									Bool: (*boolean)(ottltest.Boolp(true)),
								},
								{
									Enum: (*EnumSymbol)(ottltest.Strp("TEST_ENUM")),
								},
								{
									List: &list{
										Values: []value{
											{
												String: ottltest.Strp("test"),
											},
											{
												String: ottltest.Strp("test"),
											},
										},
									},
								},
								{
									List: &list{
										Values: []value{
											{
												String: ottltest.Strp("test"),
											},
											{
												List: &list{
													Values: []value{
														{
															String: ottltest.Strp("test"),
														},
														{
															List: &list{
																Values: []value{
																	{
																		String: ottltest.Strp("test"),
																	},
																	{
																		String: ottltest.Strp("test"),
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
								{
									Literal: &mathExprLiteral{
										Converter: &converter{
											Function: "Testing_getter",
											Arguments: []value{
												{
													Literal: &mathExprLiteral{
														Path: &Path{
															Fields: []Field{
																{
																	Name: "name",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: 9,
		},
		{
			name: "setter arg",
			inv: invocation{
				Function: "testing_setter",
				Arguments: []value{
					{
						Literal: &mathExprLiteral{
							Path: &Path{
								Fields: []Field{
									{
										Name: "name",
									},
								},
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "getsetter arg",
			inv: invocation{
				Function: "testing_getsetter",
				Arguments: []value{
					{
						Literal: &mathExprLiteral{
							Path: &Path{
								Fields: []Field{
									{
										Name: "name",
									},
								},
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "getter arg",
			inv: invocation{
				Function: "testing_getter",
				Arguments: []value{
					{
						Literal: &mathExprLiteral{
							Path: &Path{
								Fields: []Field{
									{
										Name: "name",
									},
								},
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "getter arg with nil literal",
			inv: invocation{
				Function: "testing_getter",
				Arguments: []value{
					{
						IsNil: (*isNil)(ottltest.Boolp(true)),
					},
				},
			},
			want: nil,
		},
		{
			name: "getter arg with list",
			inv: invocation{
				Function: "testing_getter",
				Arguments: []value{
					{
						List: &list{
							Values: []value{
								{
									String: ottltest.Strp("test"),
								},
								{
									Literal: &mathExprLiteral{
										Int: ottltest.Intp(1),
									},
								},
								{
									Literal: &mathExprLiteral{
										Float: ottltest.Floatp(1.1),
									},
								},
								{
									Bool: (*boolean)(ottltest.Boolp(true)),
								},
								{
									Bytes: (*byteSlice)(&[]byte{1, 2, 3, 4, 5, 6, 7, 8}),
								},
								{
									Literal: &mathExprLiteral{
										Path: &Path{
											Fields: []Field{
												{
													Name: "name",
												},
											},
										},
									},
								},
								{
									Literal: &mathExprLiteral{
										Converter: &converter{
											Function: "Testing_getter",
											Arguments: []value{
												{
													Literal: &mathExprLiteral{
														Path: &Path{
															Fields: []Field{
																{
																	Name: "name",
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "string arg",
			inv: invocation{
				Function: "testing_string",
				Arguments: []value{
					{
						String: ottltest.Strp("test"),
					},
				},
			},
			want: nil,
		},
		{
			name: "float arg",
			inv: invocation{
				Function: "testing_float",
				Arguments: []value{
					{
						Literal: &mathExprLiteral{
							Float: ottltest.Floatp(1.1),
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "int arg",
			inv: invocation{
				Function: "testing_int",
				Arguments: []value{
					{
						Literal: &mathExprLiteral{
							Int: ottltest.Intp(1),
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "bool arg",
			inv: invocation{
				Function: "testing_bool",
				Arguments: []value{
					{
						Bool: (*boolean)(ottltest.Boolp(true)),
					},
				},
			},
			want: nil,
		},
		{
			name: "byteSlice arg",
			inv: invocation{
				Function: "testing_byte_slice",
				Arguments: []value{
					{
						Bytes: (*byteSlice)(&[]byte{1, 2, 3, 4, 5, 6, 7, 8}),
					},
				},
			},
			want: nil,
		},
		{
			name: "multiple args",
			inv: invocation{
				Function: "testing_multiple_args",
				Arguments: []value{
					{
						Literal: &mathExprLiteral{
							Path: &Path{
								Fields: []Field{
									{
										Name: "name",
									},
								},
							},
						},
					},
					{
						String: ottltest.Strp("test"),
					},
					{
						Literal: &mathExprLiteral{
							Float: ottltest.Floatp(1.1),
						},
					},
					{
						Literal: &mathExprLiteral{
							Int: ottltest.Intp(1),
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "Enum arg",
			inv: invocation{
				Function: "testing_enum",
				Arguments: []value{
					{
						Enum: (*EnumSymbol)(ottltest.Strp("TEST_ENUM")),
					},
				},
			},
			want: nil,
		},
		{
			name: "telemetrySettings first",
			inv: invocation{
				Function: "testing_telemetry_settings_first",
				Arguments: []value{
					{
						String: ottltest.Strp("test0"),
					},
					{
						List: &list{
							Values: []value{
								{
									String: ottltest.Strp("test"),
								},
							},
						},
					},
					{
						Literal: &mathExprLiteral{
							Int: ottltest.Intp(1),
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "telemetrySettings middle",
			inv: invocation{
				Function: "testing_telemetry_settings_middle",
				Arguments: []value{
					{
						String: ottltest.Strp("test0"),
					},
					{
						List: &list{
							Values: []value{
								{
									String: ottltest.Strp("test"),
								},
							},
						},
					},
					{
						Literal: &mathExprLiteral{
							Int: ottltest.Intp(1),
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "telemetrySettings last",
			inv: invocation{
				Function: "testing_telemetry_settings_last",
				Arguments: []value{
					{
						String: ottltest.Strp("test0"),
					},
					{
						List: &list{
							Values: []value{
								{
									String: ottltest.Strp("test"),
								},
							},
						},
					},
					{
						Literal: &mathExprLiteral{
							Int: ottltest.Intp(1),
						},
					},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, err := p.newFunctionCall(tt.inv)
			assert.NoError(t, err)

			if tt.want != nil {
				result, _ := fn.Eval(context.Background(), nil)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

type functionWithStringSliceFactory[K any] struct{}

func (f functionWithStringSliceFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithStringSlice)
}

func (f functionWithStringSliceFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return Convert[K](reflect.ValueOf(f.functionWithStringSlice).Call(args))
}

func (f functionWithStringSliceFactory[K]) functionWithStringSlice(strs []string) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return len(strs), nil
	}, nil
}

type functionWithFloatSliceFactory[K any] struct{}

func (f functionWithFloatSliceFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithFloatSlice)
}

func (f functionWithFloatSliceFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return Convert[K](reflect.ValueOf(f.functionWithFloatSlice).Call(args))
}

func (f functionWithFloatSliceFactory[K]) functionWithFloatSlice(floats []float64) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return len(floats), nil
	}, nil
}

type functionWithIntSliceFactory[K any] struct{}

func (f functionWithIntSliceFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithIntSlice)
}

func (f functionWithIntSliceFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return Convert[K](reflect.ValueOf(f.functionWithIntSlice).Call(args))
}

func (f functionWithIntSliceFactory[K]) functionWithIntSlice(ints []int64) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return len(ints), nil
	}, nil
}

type functionWithByteSliceFactory[K any] struct{}

func (f functionWithByteSliceFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithByteSlice)
}

func (f functionWithByteSliceFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return f.functionWithByteSlice(args[0].Bytes())
}

func (f functionWithByteSliceFactory[K]) functionWithByteSlice(bytes []byte) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return len(bytes), nil
	}, nil
}

type functionWithGetterSliceFactory[K any] struct{}

func (f functionWithGetterSliceFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithGetterSlice)
}

func (f functionWithGetterSliceFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return f.functionWithGetterSlice(args[0].Interface().([]Getter[K]))
}

func (f functionWithGetterSliceFactory[K]) functionWithGetterSlice(getters []Getter[K]) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return len(getters), nil
	}, nil
}

type functionWithSetterFactory[K any] struct{}

func (f functionWithSetterFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithSetter)
}

func (f functionWithSetterFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return f.functionWithSetter(args[0].Interface().(Setter[K]))
}

func (f functionWithSetterFactory[K]) functionWithSetter(_ Setter[K]) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

type functionWithGetSetterFactory[K any] struct{}

func (f functionWithGetSetterFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithGetSetter)
}

func (f functionWithGetSetterFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return f.functionWithGetSetter(args[0].Interface().(GetSetter[K]))
}

func (f functionWithGetSetterFactory[K]) functionWithGetSetter(_ GetSetter[K]) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

type functionWithGetterFactory[K any] struct{}

func (f functionWithGetterFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithGetter)
}

func (f functionWithGetterFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return f.functionWithGetter(args[0].Interface().(Getter[K]))
}

func (f functionWithGetterFactory[K]) functionWithGetter(_ Getter[K]) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

type functionWithStringFactory[K any] struct{}

func (f functionWithStringFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithString)
}

func (f functionWithStringFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return f.functionWithString(args[0].String())
}

func (f functionWithStringFactory[K]) functionWithString(_ string) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

type functionWithFloatFactory[K any] struct{}

func (f functionWithFloatFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithFloat)
}

func (f functionWithFloatFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return f.functionWithFloat(args[0].Float())
}

func (f functionWithFloatFactory[K]) functionWithFloat(_ float64) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

type functionWithIntFactory[K any] struct{}

func (f functionWithIntFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithInt)
}

func (f functionWithIntFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return Convert[K](reflect.ValueOf(f.functionWithInt).Call(args))
}

func (f functionWithIntFactory[K]) functionWithInt(_ int64) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

type functionWithBoolFactory[K any] struct{}

func (f functionWithBoolFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithBool)
}

func (f functionWithBoolFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return f.functionWithBool(args[0].Bool())
}

func (f functionWithBoolFactory[K]) functionWithBool(_ bool) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

type functionWithMultipleArgsFactory[K any] struct{}

func (f functionWithMultipleArgsFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithMultipleArgs)
}

func (f functionWithMultipleArgsFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return f.functionWithMultipleArgs(args[0].Interface().(GetSetter[K]), args[1].String(), args[2].Float(), args[3].Int())
}

func (f functionWithMultipleArgsFactory[K]) functionWithMultipleArgs(_ GetSetter[K], _ string, _ float64, _ int64) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

type functionThatHasAnErrorFactory[K any] struct{}

func (f functionThatHasAnErrorFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionThatHasAnError)
}

func (f functionThatHasAnErrorFactory[K]) CreateFunction(_ []reflect.Value) (ExprFunc[K], error) {
	return f.functionThatHasAnError()
}

func (f functionThatHasAnErrorFactory[K]) functionThatHasAnError() (ExprFunc[K], error) {
	err := errors.New("testing")
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, err
}

type functionWithEnumFactory[K any] struct{}

func (f functionWithEnumFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithEnum)
}

func (f functionWithEnumFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return f.functionWithEnum(args[0].Interface().(Enum))
}

func (f functionWithEnumFactory[K]) functionWithEnum(_ Enum) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

type functionWithTelemetrySettingsFirstFactory[K any] struct{}

func (f functionWithTelemetrySettingsFirstFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithTelemetrySettingsFirst)
}

func (f functionWithTelemetrySettingsFirstFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return Convert[K](reflect.ValueOf(f.functionWithTelemetrySettingsFirst).Call(args))
}

func (f functionWithTelemetrySettingsFirstFactory[K]) functionWithTelemetrySettingsFirst(_ component.TelemetrySettings, _ string, _ []string, _ int64) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

type functionWithTelemetrySettingsMiddleFactory[K any] struct{}

func (f functionWithTelemetrySettingsMiddleFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithTelemetrySettingsMiddle)
}

func (f functionWithTelemetrySettingsMiddleFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return Convert[K](reflect.ValueOf(f.functionWithTelemetrySettingsMiddle).Call(args))
}

func (f functionWithTelemetrySettingsMiddleFactory[K]) functionWithTelemetrySettingsMiddle(_ string, _ []string, _ component.TelemetrySettings, _ int64) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

type functionWithTelemetrySettingsLastFactory[K any] struct{}

func (f functionWithTelemetrySettingsLastFactory[K]) GetFunctionType() reflect.Type {
	return reflect.TypeOf(f.functionWithTelemetrySettingsLast)
}

func (f functionWithTelemetrySettingsLastFactory[K]) CreateFunction(args []reflect.Value) (ExprFunc[K], error) {
	return Convert[K](reflect.ValueOf(f.functionWithTelemetrySettingsLast).Call(args))
}

func (f functionWithTelemetrySettingsLastFactory[K]) functionWithTelemetrySettingsLast(_ string, _ []string, _ int64, _ component.TelemetrySettings) (ExprFunc[K], error) {
	return func(context.Context, K) (interface{}, error) {
		return "anything", nil
	}, nil
}

func defaultFunctionsForTests() map[string]FunctionFactory[any] {
	functions := make(map[string]FunctionFactory[any])
	functions["testing_string_slice"] = functionWithStringSliceFactory[any]{}
	functions["testing_float_slice"] = functionWithFloatSliceFactory[any]{}
	functions["testing_int_slice"] = functionWithIntSliceFactory[any]{}
	functions["testing_byte_slice"] = functionWithByteSliceFactory[any]{}
	functions["testing_getter_slice"] = functionWithGetterSliceFactory[any]{}
	functions["testing_setter"] = functionWithSetterFactory[any]{}
	functions["testing_getsetter"] = functionWithGetSetterFactory[any]{}
	functions["testing_getter"] = functionWithGetterFactory[any]{}
	functions["Testing_getter"] = functionWithGetterFactory[any]{}
	functions["testing_string"] = functionWithStringFactory[any]{}
	functions["testing_float"] = functionWithFloatFactory[any]{}
	functions["testing_int"] = functionWithIntFactory[any]{}
	functions["testing_bool"] = functionWithBoolFactory[any]{}
	functions["testing_multiple_args"] = functionWithMultipleArgsFactory[any]{}
	functions["testing_enum"] = functionWithEnumFactory[any]{}
	functions["testing_telemetry_settings_first"] = functionWithTelemetrySettingsFirstFactory[any]{}
	functions["testing_telemetry_settings_middle"] = functionWithTelemetrySettingsMiddleFactory[any]{}
	functions["testing_telemetry_settings_last"] = functionWithTelemetrySettingsLastFactory[any]{}
	return functions
}
