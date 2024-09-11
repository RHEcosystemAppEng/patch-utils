// Copyright (c) 2024 Red Hat, Inc.

package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type (
	// JsonPatch encapsulates a JSON patch string
	JsonPatch struct {
		patch string
	}
	// PatchFunc is the function for you to execute to perform the patch
	PatchFunc func() error
	// NoPatchRequired is the error used to indicate your tool instructions have resulted in no patches being created
	NoPatchRequired struct {
		message string
	}
)

// Get will return the encapsulated patch string
func (p *JsonPatch) Get() string {
	return p.patch
}

// Error returns the encapsulated error message
func (e *NoPatchRequired) Error() string {
	return e.message
}

// SanitizeKeyForJsonPatch is used for preparing an annotation or a label key for patching using JSON patches. Further
// information can be found at https://jsonpatch.com/#json-pointer
func SanitizeKeyForJsonPatch(key string) string {
	return strings.ReplaceAll(strings.ReplaceAll(key, "~", "~0"), "/", "~1")
}

// UnSanitizeKeyForJsonPatch does the opposite of SanitizeKeyForJsonPatch
func UnSanitizeKeyForJsonPatch(key string) string {
	return strings.ReplaceAll(strings.ReplaceAll(key, "~1", "/"), "~0", "~")
}

// JsonPatchFinalizerInQ executes JsonPatchFinalizerInP, ignoring the patches. Use this function if you only need the
// PatchFunc
func JsonPatchFinalizerInQ(ctx context.Context, clt client.Client, obj client.Object, finalizer string) PatchFunc {
	_, f := JsonPatchFinalizerInP(ctx, clt, obj, finalizer)
	return f
}

// JsonPatchFinalizerInP executes JsonPatchFinalizerIn, panicking for errors. Use this function if you have no use of
// the returning error. The PatchFunc will still occur at runtime and might return an error; this function is only
// panicking for the errors that occurred while generating the patch
func JsonPatchFinalizerInP(ctx context.Context, clt client.Client, obj client.Object, finalizer string) (JsonPatch, PatchFunc) {
	p, f, e := JsonPatchFinalizerIn(ctx, clt, obj, finalizer)
	if e != nil {
		panic(e)
	}
	return p, f
}

// JsonPatchFinalizerIn uses JSON-type patches to add a finalizer to obj. It will return the JsonPatch for you to log,
// the PatchFunc for you to execute, and an error if it fails to generate the patch. This function isn't currently
// returning an error, but this might change in the future
func JsonPatchFinalizerIn(ctx context.Context, clt client.Client, obj client.Object, finalizer string) (JsonPatch, PatchFunc, error) {
	// create JSON patch adding finalizer
	var addFinalizerPatch JsonPatch
	if obj.GetFinalizers() == nil {
		// no finalizer found, add a list of the one finalizer
		addFinalizerPatch = JsonPatch{fmt.Sprintf("{\"op\": \"add\", \"path\": \"/metadata/finalizers\", \"value\": [\"%s\"]}", finalizer)}
	} else {
		// other finalizers exists, add the finalizer to the existing list
		addFinalizerPatch = JsonPatch{fmt.Sprintf("{\"op\": \"add\", \"path\": \"/metadata/finalizers/-\", \"value\": \"%s\"}", finalizer)}
	}

	return addFinalizerPatch, func() error {
		return clt.Patch(ctx, obj, client.RawPatch(types.JSONPatchType, []byte("["+addFinalizerPatch.Get()+"]")))
	}, nil
}

// JsonPatchFinalizerOutQ executes JsonPatchFinalizerOutP, ignoring the patches. Use this function if you only need the
// PatchFunc
func JsonPatchFinalizerOutQ(ctx context.Context, clt client.Client, obj client.Object, finalizer string) PatchFunc {
	_, f := JsonPatchFinalizerOutP(ctx, clt, obj, finalizer)
	return f
}

// JsonPatchFinalizerOutP executes JsonPatchFinalizerOut, panicking for errors. Use this function if you have no use of
// the returning error. The PatchFunc will still occur at runtime and might return an error; this function is only
// panicking for the errors that occurred while generating the patch
func JsonPatchFinalizerOutP(ctx context.Context, clt client.Client, obj client.Object, finalizer string) (JsonPatch, PatchFunc) {
	p, f, e := JsonPatchFinalizerOut(ctx, clt, obj, finalizer)
	if e != nil {
		panic(e)
	}
	return p, f
}

// JsonPatchFinalizerOut uses JSON-type patches to remove a finalizer from the obj. It will return the JsonPatch for you
// to log, the PatchFunc for you to execute, and an error if it fails to generate the patch
func JsonPatchFinalizerOut(ctx context.Context, clt client.Client, obj client.Object, finalizer string) (JsonPatch, PatchFunc, error) {
	// create JSON patch removing finalizer
	var removeFinalizerPatch JsonPatch
	if len(obj.GetFinalizers()) == 1 {
		// remove all finalizers if only one exists
		removeFinalizerPatch = JsonPatch{"{\"op\": \"remove\", \"path\": \"/metadata/finalizers\"}"}
	} else {
		// remove index-based specific finalizer if more than one exist
		for idx, fin := range obj.GetFinalizers() {
			if fin == finalizer {
				removeFinalizerPatch = JsonPatch{fmt.Sprintf("{\"op\": \"remove\", \"path\": \"/metadata/finalizers/%d\"}", idx)}
				break
			}
		}
		if len(removeFinalizerPatch.Get()) == 0 {
			return removeFinalizerPatch, nil, &NoPatchRequired{"finalizer index not found"}
		}
	}

	return removeFinalizerPatch, func() error {
		return clt.Patch(ctx, obj, client.RawPatch(types.JSONPatchType, []byte("["+removeFinalizerPatch.Get()+"]")))
	}, nil
}

// JsonPatchMapQ executes JsonPatchMapP, ignoring the patches. Use this function if you only need the PatchFunc
func JsonPatchMapQ(ctx context.Context, clt client.Client, obj client.Object, path string, origMap, newMap map[string]string) PatchFunc {
	_, f := JsonPatchMapP(ctx, clt, obj, path, origMap, newMap)
	return f
}

// JsonPatchMapP executes JsonPatchMap, panicking for errors. Use this function if you have no use of the returning
// error. The PatchFunc will still occur at runtime and might return an error; this function is only panicking for the
// errors that occurred while generating the patch
func JsonPatchMapP(ctx context.Context, clt client.Client, obj client.Object, path string, origMap, newMap map[string]string) ([]JsonPatch, PatchFunc) {
	p, f, e := JsonPatchMap(ctx, clt, obj, path, origMap, newMap)
	if e != nil {
		panic(e)
	}
	return p, f
}

// JsonPatchMap uses JSON-type patches to add or replace all members in the given path for the given obj with the member
// in newMap; we use origMap to determine whether we need to add or replace. It will return JsonPatches for you to log,
// the PatchFunc for you to execute, and an error if it fails to generate the patch
func JsonPatchMap(ctx context.Context, clt client.Client, obj client.Object, path string, origMap, newMap map[string]string) ([]JsonPatch, PatchFunc, error) {
	patchNewMapTemplate := "{\"op\": \"add\", \"path\": \"%s\", \"value\": {\"%s\": \"%s\"}}"
	patchExistingMapTemplate := "{\"op\": \"%s\", \"path\": \"%s/%s\", \"value\": \"%s\"}"

	var patches []string

	if len(origMap) == 0 {
		// no previous member exists - load add all given member to the patch
		first := true
		for k, v := range newMap {
			if first {
				// if the original map doesn't exist, the key should be un-sanitized and part of the value
				patches = append(patches, fmt.Sprintf(patchNewMapTemplate, path, UnSanitizeKeyForJsonPatch(k), v))
				first = false
			} else {
				patches = append(patches, fmt.Sprintf(patchExistingMapTemplate, "add", path, k, v))
			}
		}
	} else {
		// found previous member - verify add/replace/exists before loading the given member to the patch
		for k, v := range newMap {
			if value, found := origMap[k]; found {
				if v != value {
					// found existing member with the key and a different value - replace
					patches = append(patches, fmt.Sprintf(patchExistingMapTemplate, "replace", path, k, v))
				}
			} else {
				// existing member with the key not found - add
				patches = append(patches, fmt.Sprintf(patchExistingMapTemplate, "add", path, k, v))
			}
		}
	}

	var patchObjs []JsonPatch
	for _, p := range patches {
		patchObjs = append(patchObjs, JsonPatch{p})
	}

	if len(patches) < 1 {
		return patchObjs, nil, &NoPatchRequired{"nothing to patch in map"}
	}
	return patchObjs, func() error {
		return clt.Patch(ctx, obj, client.RawPatch(types.JSONPatchType, []byte("["+strings.Join(patches, ",")+"]")))
	}, nil
}

// JsonPatchSpecQ executes JsonPatchSpecP, ignoring the patches. Use this function if you only need the PatchFunc
func JsonPatchSpecQ[T interface{}](ctx context.Context, clt client.Client, obj client.Object, spec *T) PatchFunc {
	_, f := JsonPatchSpecP(ctx, clt, obj, spec)
	return f
}

// JsonPatchSpecP executes JsonPatchSpec, panicking for errors. Use this function if you have no use of the returning
// error. The PatchFunc will still occur at runtime and might return an error; this function is only panicking for the
// errors that occurred while generating the patch
func JsonPatchSpecP[T interface{}](ctx context.Context, clt client.Client, obj client.Object, spec *T) (JsonPatch, PatchFunc) {
	p, f, e := JsonPatchSpec(ctx, clt, obj, spec)
	if e != nil {
		panic(e)
	}
	return p, f
}

// JsonPatchSpec uses JSON-type patches to replace a Spec in obj. It will return a JsonPatch for you to log, the
// PatchFunc for you to execute, and an error if it fails to generate the patch
func JsonPatchSpec[T interface{}](ctx context.Context, clt client.Client, obj client.Object, spec *T) (JsonPatch, PatchFunc, error) {
	marshaled, err := json.Marshal(spec)
	if err != nil {
		return JsonPatch{""}, nil, err
	}

	patch := JsonPatch{"{\"op\": \"replace\", \"path\": \"/spec\", \"value\": " + string(marshaled) + "}"}
	return patch, func() error {
		return clt.Patch(ctx, obj, client.RawPatch(types.JSONPatchType, []byte("["+patch.Get()+"]")))
	}, nil
}
