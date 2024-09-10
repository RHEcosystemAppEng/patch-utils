# Kubernetes Operator Patch Utils

When developing Kubernetes Operators for production use, patching is often the most robust approach for updating
existing objects. This module hosts utilities for creating standard patching functions.

> [!WARNING]
>
> This is a WIP, use at your own risk.

## Usage

### JSON Patching

```go
package yourpkg

import (
    "context"
    patchutils "github.com/rhecosystemappeng/patch-utils/pkg"
    "sigs.k8s.io/controller-runtime/pkg/client"
)

func yourFunc() {
    // assign with your own
    var ctx context.Context
    var clt client.Client
    var obj client.Object

    // #####################################################################
    // ##### Patch Labels, Annotations, or any other map[string]string #####
    // #####################################################################
    // patches is a list of patchutils.JsonPatch that will get executed once you execute the patchFunc, patchutils.PatchFunc.
    // err is an error from the patch generation process; executing patchFunc might return an error arising from the API calls.
    // use the wrapper function patchutils.JsonPatchMapP if you don't need the generation error.
    // use the wrapper function patchutils.JsonPatchMapQ if you only need the patchFunc.
    patches, patchFunc, err := patchutils.JsonPatchMap(ctx, clt, obj, "/metadata/labels", obj.GetLabels(), map[string]string{
        patchutils.SanitizeKeyForJsonPatch("your.label/first.key"): "label-value",
        patchutils.SanitizeKeyForJsonPatch("your.label/second.key"): "another-label-value",
    })

    // #################################
    // ##### Patch any Spec object #####
    // #################################
    // assign with your own spec
    var spec *interface{}
    // patch is the patchutils.JsonPatch that will get executed once you execute the patchFunc, patchutils.PatchFunc.
    // err is an error from the patch generation process; executing patchFunc might return an error arising from the API calls.
    // use the wrapper function patchutils.JsonPatchSpecP if you don't need the generation error.
    // use the wrapper function patchutils.JsonPatchSpecQ if you only need the patchFunc.
    patch, patchFunc, err := patchutils.JsonPatchSpec(ctx, clt, obj, spec)

    // ############################################
    // ##### Patch a finalizer into an object #####
    // ############################################
    // patch is the patchutils.JsonPatch that will get executed once you execute the patchFunc, patchutils.PatchFunc.
    // err is an error from the patch generation process; executing patchFunc might return an error arising from the API calls.
    // use the wrapper function patchutils.JsonPatchFinalizerInP if you don't need the generation error.
    // use the wrapper function patchutils.JsonPatchFinalizerInQ if you only need the patchFunc.
    patch, patchFunc, err := patchutils.JsonPatchFinalizerIn(ctx, clt, obj, "my.custom/cleanup-finalizer")

    // ##############################################
    // ##### Patch a finalizer out of an object #####
    // ##############################################
    // patch is the patchutils.JsonPatch that will get executed once you execute the patchFunc, patchutils.PatchFunc.
    // err is an error from the patch generation process; executing patchFunc might return an error arising from the API calls.
    // use the wrapper function patchutils.JsonPatchFinalizerOutP if you don't need the generation error.
    // use the wrapper function patchutils.JsonPatchFinalizerOutQ if you only need the patchFunc.
    patch, patchFunc, err := patchutils.JsonPatchFinalizerOut(ctx, clt, obj, "my.custom/cleanup-finalizer")
}

```
