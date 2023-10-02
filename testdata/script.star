
v1beta1 = proto.file("v1beta1/run_function.proto")
structpb = proto.file("google/protobuf/struct.proto")

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

