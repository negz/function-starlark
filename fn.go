package main

import (
	"context"
	"fmt"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"go.starlark.net/lib/json"
	starlarkproto "go.starlark.net/lib/proto"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"

	"github.com/crossplane/function-sdk-go/errors"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/response"

	"github.com/negz/function-starlark/input/v1beta1"
)

// Function returns whatever response you ask it to.
type Function struct {
	fnv1beta1.UnimplementedFunctionRunnerServiceServer

	log logging.Logger
}

// RunFunction runs the Function.
func (f *Function) RunFunction(_ context.Context, req *fnv1beta1.RunFunctionRequest) (*fnv1beta1.RunFunctionResponse, error) {
	log := f.log.WithValues("tag", req.GetMeta().GetTag())
	log.Info("Running Function", "tag")

	rsp := response.To(req, response.DefaultTTL)

	data, err := proto.Marshal(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot marshal %T to protocol buffer bytes", req))
		return rsp, nil
	}

	reqv, err := starlarkproto.Unmarshal(req.ProtoReflect().Descriptor(), data)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot create Starlark value from %T protocol buffer bytes", req))
		return rsp, nil
	}

	in := &v1beta1.Script{}
	if err := request.GetInput(req, in); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get Starlark script from %T", req))
		return rsp, nil
	}

	if in.Inline == nil {
		response.Fatal(rsp, errors.New("%T: inline script is required"))
		return rsp, nil
	}

	log.Debug("Running inline script", "script", string(*in.Inline))

	thread := &starlark.Thread{
		Name: req.GetMeta().GetTag(),
		Print: func(_ *starlark.Thread, msg string) {
			log.Debug("Starlark print() called", "starlark-msg", msg)
			fmt.Println(msg)
		},
	}

	// This allows protofile("v1beta1/run_function.proto") to work.
	starlarkproto.SetPool(thread, protoregistry.GlobalFiles)

	predeclared := starlark.StringDict{
		"proto": starlarkproto.Module,
		"json":  json.Module,
	}
	globals, err := starlark.ExecFile(thread, "main.star", []byte(*in.Inline), predeclared)
	if err != nil {
		response.Fatal(rsp, errors.Wrap(err, "cannot execute Starlark script"))
		return rsp, nil
	}

	mainv, ok := globals["main"]
	if !ok {
		response.Fatal(rsp, errors.New("Starlark script did not export a main() function"))
		return rsp, nil
	}

	main, ok := mainv.(*starlark.Function)
	if !ok {
		response.Fatal(rsp, errors.New("Starlark script exported a main global that is not a function"))
		return rsp, nil
	}

	v, err := starlark.Call(thread, main, starlark.Tuple{reqv}, nil)
	if err != nil {
		if eerr := (&starlark.EvalError{}); errors.As(err, &eerr) {
			response.Fatal(rsp, errors.Errorf("Starlark script error: %s", eerr.Backtrace()))
			return rsp, nil
		}
		response.Fatal(rsp, errors.Wrap(err, "cannot call Starlark main() function"))
		return rsp, nil
	}

	rspv, ok := v.(*starlarkproto.Message)
	if !ok {
		response.Fatal(rsp, errors.New("Starlark script main() function did not return a protocol buffer message"))
		return rsp, nil
	}

	data, err = proto.Marshal(rspv.Message())
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot marshal %T to protocol buffer bytes", rspv))
		return rsp, nil
	}

	if err := proto.Unmarshal(data, rsp); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot unmarshal protocol buffer bytes into %T", rsp))
		return rsp, nil
	}

	// Make sure we keep the tag from our request.
	if rsp.Meta == nil {
		rsp.Meta = &fnv1beta1.ResponseMeta{}
	}
	rsp.Meta.Tag = req.GetMeta().GetTag()

	return rsp, nil
}
