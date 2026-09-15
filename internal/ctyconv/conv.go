// Package ctyconv converts Go runtime values to cty.Value and evaluates HCL expressions.
package ctyconv

import (
	"fmt"
	"math/big"
	"strings"
	"time"
	"uuid"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

// ToCty converts any Go primitive, slice, or map into a cty.Value.
func ToCty(val any) cty.Value {
	if val == nil {
		return cty.NilVal
	}

	switch v := val.(type) {
	case string:
		return cty.StringVal(v)
	case bool:
		return cty.BoolVal(v)
	case int:
		return cty.NumberIntVal(int64(v))
	case int64:
		return cty.NumberIntVal(v)
	case float64:
		return cty.NumberFloatVal(v)
	case map[string]any:
		if len(v) == 0 {
			return cty.EmptyObjectVal
		}
		attrs := make(map[string]cty.Value, len(v))
		for k, item := range v {
			attrs[k] = ToCty(item)
		}
		return cty.ObjectVal(attrs)
	case map[string]string:
		if len(v) == 0 {
			return cty.EmptyObjectVal
		}
		attrs := make(map[string]cty.Value, len(v))
		for k, item := range v {
			attrs[k] = cty.StringVal(item)
		}
		return cty.ObjectVal(attrs)
	case []any:
		if len(v) == 0 {
			return cty.EmptyTupleVal
		}
		elems := make([]cty.Value, len(v))
		for i, item := range v {
			elems[i] = ToCty(item)
		}
		return cty.TupleVal(elems)
	case []map[string]any:
		if len(v) == 0 {
			return cty.EmptyTupleVal
		}
		elems := make([]cty.Value, len(v))
		for i, item := range v {
			elems[i] = ToCty(item)
		}
		return cty.TupleVal(elems)
	default:
		return cty.StringVal(fmt.Sprintf("%v", v))
	}
}

// ToNative converts a cty.Value back into standard Go primitives, slices, or maps.
func ToNative(val cty.Value) any {
	if val.IsNull() || !val.IsKnown() {
		return nil
	}

	ty := val.Type()
	switch {
	case ty == cty.String:
		return val.AsString()
	case ty == cty.Bool:
		return val.True()
	case ty == cty.Number:
		bf := val.AsBigFloat()
		if i, acc := bf.Int64(); acc == big.Exact {
			return i
		}
		f, _ := bf.Float64()
		return f
	case ty.IsTupleType() || ty.IsListType() || ty.IsSetType():
		out := make([]any, 0, val.LengthInt())
		for it := val.ElementIterator(); it.Next(); {
			_, el := it.Element()
			out = append(out, ToNative(el))
		}
		return out
	case ty.IsObjectType() || ty.IsMapType():
		out := make(map[string]any)
		for it := val.ElementIterator(); it.Next(); {
			k, el := it.Element()
			out[k.AsString()] = ToNative(el)
		}
		return out
	default:
		return nil
	}
}

// EvalBool evaluates an HCL expression expecting a boolean result.
func EvalBool(expr hcl.Expression, ctx *hcl.EvalContext) (bool, error) {
	if expr == nil {
		return true, nil
	}
	val, diags := expr.Value(ctx)
	if diags.HasErrors() {
		return false, diags
	}
	if val.Type() != cty.Bool {
		return false, fmt.Errorf("expected boolean expression, got %s", val.Type().FriendlyName())
	}
	return val.True(), nil
}

// EvalAny evaluates an HCL expression into a native Go value.
func EvalAny(expr hcl.Expression, ctx *hcl.EvalContext) (any, error) {
	if expr == nil {
		return nil, nil
	}
	val, diags := expr.Value(ctx)
	if diags.HasErrors() {
		return nil, diags
	}
	return ToNative(val), nil
}

// BuiltinFunctions returns standard helper functions available in HCL expressions.
func BuiltinFunctions() map[string]function.Function {
	return map[string]function.Function{
		"now": function.New(&function.Spec{
			Params: []function.Parameter{},
			Type:   function.StaticReturnType(cty.String),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				return cty.StringVal(time.Now().UTC().Format(time.RFC3339)), nil
			},
		}),
		"uuid": function.New(&function.Spec{
			Params: []function.Parameter{},
			Type:   function.StaticReturnType(cty.String),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				return cty.StringVal(uuid.New().String()), nil
			},
		}),
		"problem": function.New(&function.Spec{
			Params: []function.Parameter{
				{Name: "status", Type: cty.Number},
				{Name: "detail", Type: cty.String},
			},
			Type: function.StaticReturnType(cty.Object(map[string]cty.Type{
				"status": cty.Number,
				"title":  cty.String,
				"detail": cty.String,
			})),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				status, _ := args[0].AsBigFloat().Int64()
				detail := args[1].AsString()
				return cty.ObjectVal(map[string]cty.Value{
					"status": cty.NumberIntVal(status),
					"title":  cty.StringVal("Error"),
					"detail": cty.StringVal(detail),
				}), nil
			},
		}),
		"upper": function.New(&function.Spec{
			Params: []function.Parameter{{Name: "str", Type: cty.String}},
			Type:   function.StaticReturnType(cty.String),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				return cty.StringVal(strings.ToUpper(args[0].AsString())), nil
			},
		}),
		"lower": function.New(&function.Spec{
			Params: []function.Parameter{{Name: "str", Type: cty.String}},
			Type:   function.StaticReturnType(cty.String),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				return cty.StringVal(strings.ToLower(args[0].AsString())), nil
			},
		}),
		"trim": function.New(&function.Spec{
			Params: []function.Parameter{{Name: "str", Type: cty.String}},
			Type:   function.StaticReturnType(cty.String),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				return cty.StringVal(strings.TrimSpace(args[0].AsString())), nil
			},
		}),
		"trimprefix": function.New(&function.Spec{
			Params: []function.Parameter{
				{Name: "str", Type: cty.String},
				{Name: "prefix", Type: cty.String},
			},
			Type: function.StaticReturnType(cty.String),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				return cty.StringVal(strings.TrimPrefix(args[0].AsString(), args[1].AsString())), nil
			},
		}),
		"trimsuffix": function.New(&function.Spec{
			Params: []function.Parameter{
				{Name: "str", Type: cty.String},
				{Name: "suffix", Type: cty.String},
			},
			Type: function.StaticReturnType(cty.String),
			Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
				return cty.StringVal(strings.TrimSuffix(args[0].AsString(), args[1].AsString())), nil
			},
		}),
	}
}
