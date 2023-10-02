package main

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/negz/function-starlark/input/v1beta1"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/utils/pointer"

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/resource"
)

func TestRunFunction(t *testing.T) {

	type args struct {
		ctx context.Context
		req *fnv1beta1.RunFunctionRequest
	}
	type want struct {
		rsp *fnv1beta1.RunFunctionResponse
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"RunTheScript": {
			reason: "The Function should run a starlark script...",
			args: args{
				req: &fnv1beta1.RunFunctionRequest{
					Input: MustScript("testdata/script.star"),
					Desired: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							ConnectionDetails: map[string][]byte{
								"very": []byte("secret"),
							},
						},
						Resources: map[string]*fnv1beta1.Resource{
							"existing-resource": {
								Ready: fnv1beta1.Ready_READY_TRUE,
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1beta1.RunFunctionResponse{
					Desired: &fnv1beta1.State{
						Composite: &fnv1beta1.Resource{
							ConnectionDetails: map[string][]byte{
								"very": []byte("secret"),
							},
						},
						Resources: map[string]*fnv1beta1.Resource{
							"existing-resource": {
								Ready: fnv1beta1.Ready_READY_TRUE,
							},
							"my-desired-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "example.org/v1",
									"kind": "ComposedResource",
									"spec": {
										"widgets": 42
									}
								}`),
								Ready: fnv1beta1.Ready_READY_TRUE,
							},
						},
					},
					Results: []*fnv1beta1.Result{{
						Severity: fnv1beta1.Severity_SEVERITY_NORMAL,
						Message:  "Hi from Starlark!",
					}},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := &Function{log: logging.NewNopLogger()}
			rsp, err := f.RunFunction(tc.args.ctx, tc.args.req)

			if diff := cmp.Diff(tc.want.rsp, rsp, protocmp.Transform(), protocmp.IgnoreEmptyMessages()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want rsp, +got rsp:\n%s", tc.reason, diff)
			}

			if diff := cmp.Diff(tc.want.err, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want err, +got err:\n%s", tc.reason, diff)
			}
		})
	}
}

func MustScript(file string) *structpb.Struct {
	b, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	s := &v1beta1.Script{
		Source: v1beta1.ScriptSourceInline,
		Inline: pointer.String(string(b)),
	}
	st, err := resource.AsStruct(s)
	if err != nil {
		panic(err)
	}
	return st
}
