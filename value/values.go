package value

import (
	"go.starlark.net/lib/time"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/crossplane/function-sdk-go/errors"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
)

// TODO(negz): Should this be a module?

// RunFunctionRequest is a Starlark representation of a RunFunctionRequest.
type RunFunctionRequest struct {
	wrapped *fnv1beta1.RunFunctionRequest
}

func NewRunFunctionRequest(req *fnv1beta1.RunFunctionRequest) *RunFunctionRequest {
	return &RunFunctionRequest{wrapped: req}
}

var _ starlark.Value = &RunFunctionRequest{}
var _ starlark.HasAttrs = &RunFunctionRequest{}

func (r *RunFunctionRequest) String() string       { return r.wrapped.String() }
func (r *RunFunctionRequest) Type() string         { return "RunFunctionRequest" }
func (r *RunFunctionRequest) Freeze()              {}
func (r *RunFunctionRequest) Truth() starlark.Bool { return starlark.True }
func (r *RunFunctionRequest) Hash() (uint32, error) {
	return 0, errors.Errorf("unhashable: %s", r.Type())
}

func (r *RunFunctionRequest) AttrNames() []string {
	return []string{"meta", "observed", "desired", "input"}
}

func (r *RunFunctionRequest) Attr(name string) (starlark.Value, error) {
	switch name {
	case "meta":
		return &RequestMeta{wrapped: r.wrapped.GetMeta()}, nil
	case "observed":
		return &State{wrapped: r.wrapped.GetObserved()}, nil
	case "desired":
		return &State{wrapped: r.wrapped.GetDesired()}, nil
	case "input":
		return FromProtoStruct(r.wrapped.GetInput()), nil
	default:
		return nil, nil
	}
}

type RequestMeta struct {
	wrapped *fnv1beta1.RequestMeta
}

var _ starlark.Value = &RequestMeta{}
var _ starlark.HasAttrs = &RequestMeta{}

func (r *RequestMeta) String() string       { return r.wrapped.String() }
func (r *RequestMeta) Type() string         { return "RequestMeta" }
func (r *RequestMeta) Freeze()              {}
func (r *RequestMeta) Truth() starlark.Bool { return starlark.True }
func (r *RequestMeta) Hash() (uint32, error) {
	return 0, errors.Errorf("unhashable: %s", r.Type())
}

func (r *RequestMeta) AttrNames() []string {
	return []string{"tag"}
}

func (r *RequestMeta) Attr(name string) (starlark.Value, error) {
	switch name {
	case "tag":
		return starlark.String(r.wrapped.GetTag()), nil
	default:
		return nil, nil
	}
}

type State struct {
	wrapped *fnv1beta1.State
}

var _ starlark.Value = &State{}
var _ starlark.HasAttrs = &State{}

func (r *State) String() string       { return r.wrapped.String() }
func (r *State) Type() string         { return "State" }
func (r *State) Freeze()              {}
func (r *State) Truth() starlark.Bool { return starlark.True }
func (r *State) Hash() (uint32, error) {
	return 0, errors.Errorf("unhashable: %s", r.Type())
}

func (r *State) AttrNames() []string {
	return []string{"composite", "resources"}
}

func (r *State) Attr(name string) (starlark.Value, error) {
	switch name {
	case "composite":
		return &Resource{wrapped: r.wrapped.GetComposite()}, nil
	case "resources":
		d := starlark.NewDict(len(r.wrapped.GetResources()))
		for name, r := range r.wrapped.GetResources() {
			_ = d.SetKey(starlark.String(name), &Resource{wrapped: r})
		}
		return d, nil
	default:
		return nil, nil
	}
}

type Resource struct {
	wrapped *fnv1beta1.Resource
}

var _ starlark.Value = &Resource{}
var _ starlark.HasAttrs = &Resource{}

func (r *Resource) String() string       { return r.wrapped.String() }
func (r *Resource) Type() string         { return "Resource" }
func (r *Resource) Freeze()              {}
func (r *Resource) Truth() starlark.Bool { return starlark.True }
func (r *Resource) Hash() (uint32, error) {
	return 0, errors.Errorf("unhashable: %s", r.Type())
}

func (r *Resource) AttrNames() []string {
	return []string{"resource", "connection_details", "ready"}
}

func (r *Resource) Attr(name string) (starlark.Value, error) {
	switch name {
	case "resource":
		return FromProtoStruct(r.wrapped.GetResource()), nil
	case "connection_details":
		d := starlark.NewDict(len(r.wrapped.GetConnectionDetails()))
		for k, v := range r.wrapped.GetConnectionDetails() {
			_ = d.SetKey(starlark.String(k), starlark.Bytes(v))
		}
		return nil, nil
	case "ready":
		return starlark.String(r.wrapped.GetReady().String()), nil
	default:
		return nil, nil
	}
}

type RunFunctionResponse struct {
	wrapped *fnv1beta1.RunFunctionResponse
}

var _ starlark.Value = &RunFunctionResponse{}
var _ starlark.HasAttrs = &RunFunctionResponse{}
var _ starlark.HasSetField = &RunFunctionResponse{}

func FromStarlarkResponse(srsp *RunFunctionResponse) *fnv1beta1.RunFunctionResponse {
	return srsp.wrapped
}

func (r *RunFunctionResponse) String() string       { return "RunFunctionResponse(...)" }
func (r *RunFunctionResponse) Type() string         { return "RunFunctionResponse" }
func (r *RunFunctionResponse) Freeze()              {}
func (r *RunFunctionResponse) Truth() starlark.Bool { return starlark.True }
func (r *RunFunctionResponse) Hash() (uint32, error) {
	return 0, errors.Errorf("unhashable: %s", r.Type())
}

func (r *RunFunctionResponse) AttrNames() []string {
	return []string{"meta", "desired", "results"}
}

func (r *RunFunctionResponse) Attr(name string) (starlark.Value, error) {
	switch name {
	case "meta":
		return &ResponseMeta{r.wrapped.GetMeta()}, nil
	case "desired":
		return &State{r.wrapped.GetDesired()}, nil
	case "results":
		rs := r.wrapped.GetResults()
		elems := make([]starlark.Value, len(rs))
		for i := range rs {
			elems[i] = &Result{wrapped: rs[i]}
		}
		return starlark.NewList(elems), nil
	}
	return nil, nil
}

func (r *RunFunctionResponse) SetField(name string, v starlark.Value) error {
	return nil
}

type ResponseMeta struct {
	wrapped *fnv1beta1.ResponseMeta
}

var _ starlark.Value = &ResponseMeta{}
var _ starlark.HasAttrs = &ResponseMeta{}

func (r *ResponseMeta) String() string       { return r.wrapped.String() }
func (r *ResponseMeta) Type() string         { return "ResponseMeta" }
func (r *ResponseMeta) Freeze()              {}
func (r *ResponseMeta) Truth() starlark.Bool { return starlark.True }
func (r *ResponseMeta) Hash() (uint32, error) {
	return 0, errors.Errorf("unhashable: %s", r.Type())
}

func (r *ResponseMeta) AttrNames() []string {
	return []string{"tag", "ttl"}
}

func (r *ResponseMeta) Attr(name string) (starlark.Value, error) {
	switch name {
	case "tag":
		return starlark.String(r.wrapped.GetTag()), nil
	case "ttl":
		return time.Duration(r.wrapped.Ttl.AsDuration()), nil
	default:
		return nil, nil
	}
}

type Result struct {
	wrapped *fnv1beta1.Result
}

var _ starlark.Value = &Result{}
var _ starlark.HasAttrs = &Result{}

func (r *Result) String() string       { return r.wrapped.String() }
func (r *Result) Type() string         { return "Result" }
func (r *Result) Freeze()              {}
func (r *Result) Truth() starlark.Bool { return starlark.True }
func (r *Result) Hash() (uint32, error) {
	return 0, errors.Errorf("unhashable: %s", r.Type())
}

func (r *Result) AttrNames() []string {
	return []string{"severity", "message"}
}

func (r *Result) Attr(name string) (starlark.Value, error) {
	switch name {
	case "severity":
		return starlark.String(r.wrapped.GetSeverity().String()), nil
	case "message":
		return starlark.String(r.wrapped.GetMessage()), nil
	default:
		return nil, nil
	}
}

// Copied from https://github.com/FuzzyMonkeyCo/monkey/blob/06652b9/pkg/starlarkvalue/from_protovalue.go#L59

// FromProtoValue converts a Google Well-Known-Type Value to a Starlark value.
// Panics on unexpected proto value.
func FromProtoValue(x *structpb.Value) starlark.Value {
	switch x.GetKind().(type) {

	case *structpb.Value_NullValue:
		return starlark.None

	case *structpb.Value_BoolValue:
		return starlark.Bool(x.GetBoolValue())

	case *structpb.Value_NumberValue:
		return starlark.Float(x.GetNumberValue())

	case *structpb.Value_StringValue:
		return starlark.String(x.GetStringValue())

	case *structpb.Value_ListValue:
		xs := x.GetListValue().GetValues()
		values := make([]starlark.Value, 0, len(xs))
		for _, x := range xs {
			value := FromProtoValue(x)
			values = append(values, value)
		}
		return starlark.NewList(values)

	case *structpb.Value_StructValue:
		kvs := x.GetStructValue().GetFields()
		values := starlark.NewDict(len(kvs))
		for k, v := range kvs {
			_ = values.SetKey(starlark.String(k), FromProtoValue(v)) // unreachable: hashable key, not iterating, not frozen.
		}
		return values

	default:
		return starlark.None // unreachable: only proto values.
	}
}

func FromProtoStruct(s *structpb.Struct) starlark.Value {
	if s == nil {
		return starlark.NewDict(0)
	}
	return FromProtoValue(structpb.NewStructValue(s))
}
