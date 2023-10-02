# function-starlark

A [Crossplane] Composition Function that lets you compose resources using
[Starlark], a limited dialect of Python.

## What is this?

This Function lets you express composition logic in Starlark. Currently it only
supports loading a script specified inline in your Composition. This makes it
best for smaller tasks, where you want the expressiveness of a programming
language but don't want to have to actually _build_ your own Function.

Note that this Function's Starlark API is not yet stable. It currently requires
a branch of the Starlark proto package, until Starlark PR [#511] is merged.

Here's an example of a Composition that uses this Composition Function.

```yaml
apiVersion: apiextensions.crossplane.io/v1
kind: Composition
metadata:
  name: test-crossplane
spec:
  compositeTypeRef:
    apiVersion: database.example.com/v1alpha1
    kind: NoSQL
  mode: Pipeline
  pipeline:
  - step: compose-a-resource-with-starlark
    functionRef:
      name: function-starlark
    input:
      apiVersion: starlark.fn.crossplane.io/v1beta1
      kind: Script
      source: Inline
      inline: |
        v1beta1 = proto.file("v1beta1/run_function.proto")
        structpb = proto.file("google/protobuf/struct.proto")

        # Your Starlark script must define a Function named main.
        # The main Function is passed a RunFunctionRequest in the form of a
        # Pythonic object. It must return a RunFunctionResponse.
        def main(req):
            # This structpb API is pretty awkward. We should probably improve it per
            # https://github.com/google/starlark-go/issues/513.
            resources = {
                "my-desired-resource": v1beta1.Resource(
                    resource=structpb.Struct(
                        fields={
                            "apiVersion": structpb.Value(string_value="example.org/v1"),
                            "kind": structpb.Value(string_value="ComposedResource"),
                            "spec": structpb.Value(struct_value=structpb.Struct(
                                fields={
                                    "widgets": structpb.Value(number_value=42)
                                }
                            ))
                        }
                    ),
                    ready=v1beta1.Ready.READY_TRUE
                )
            }

            # Note that you can't append to rsp.results, or even set it to an empty list
            # then append to to. See https://github.com/google/starlark-go/issues/512.
            results = [v1beta1.Result(
                severity=v1beta1.Severity.SEVERITY_NORMAL,
                message="Hi from Starlark!"
            )]

            # Keep any existing desired resources.
            resources.update(req.desired.resources)

            return v1beta1.RunFunctionResponse(
                desired=v1beta1.State(
                    composite=req.desired.composite,
                    resources=resources
                ),
                results=results,
            )
```

Notice that it has a `pipeline` (of Composition Functions) instead of an array
of `resources`.

## Developing this Function

This Function doesn't use the typical Crossplane build submodule and Makefile,
since we'd like Functions to have a less heavyweight developer experience.
It mostly relies on regular old Go tools:

```shell
# Run code generation - see input/generate.go
$ go generate ./...

# Run tests
$ go test -cover ./...
?       github.com/crossplane/function-template-go/input/v1beta1      [no test files]
ok      github.com/crossplane/function-template-go    0.006s  coverage: 25.8% of statements

# Lint the code
$ docker run --rm -v $(pwd):/app -v ~/.cache/golangci-lint/v1.54.2:/root/.cache -w /app golangci/golangci-lint:v1.54.2 golangci-lint run

# Build a Docker image - see Dockerfile
$ docker build .
```

This Function can be pushed to any Docker registry. To push to xpkg.upbound.io
use `docker push` and `docker-credential-up` from
https://github.com/upbound/up/.

[Crossplane]: https://crossplane.io
[function-design]: https://github.com/crossplane/crossplane/blob/3996f20/design/design-doc-composition-functions.md
[function-pr]: https://github.com/crossplane/crossplane/pull/4500
[Starlark]: https://github.com/google/starlark-go
[#511]: https://github.com/google/starlark-go/pull/511