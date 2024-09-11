// Copyright (c) 2024 Red Hat, Inc.

package pkg

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	testdata "github.com/rhecosystemappeng/patch-utils/pkg/testdata/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Context("Sanitizing and Un-Sanitizing map keys for JSON patch", func() {
	It("should sanitize key correctly", func() {
		Expect(SanitizeKeyForJsonPatch("my~map/key")).To(Equal("my~0map~1key"))
	})

	It("should un-sanitize key correctly", func() {
		Expect(UnSanitizeKeyForJsonPatch("my~0map~1key")).To(Equal("my~map/key"))
	})
})

var _ = Context("JSON patch", func() {
	Context("a map", func() {
		It("should work when patching one member into an empty map", func(ctx SpecContext) {
			obj := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dummy-map-obj-1",
				},
			}

			Expect(clt.Create(ctx, obj)).To(Succeed())

			Expect(JsonPatchMapQ(ctx, clt, obj, "/metadata/annotations", obj.Annotations, map[string]string{
				"annotation_key1": "annotation_value1",
			})()).To(Succeed())

			Expect(clt.Get(ctx, types.NamespacedName{Name: obj.Name}, obj)).To(Succeed())
			Expect(obj.Annotations["annotation_key1"]).To(Equal("annotation_value1"))
		})

		It("should work when patching multiple members into an empty map", func(ctx SpecContext) {
			obj := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dummy-map-obj-2",
				},
			}

			Expect(clt.Create(ctx, obj)).To(Succeed())

			Expect(JsonPatchMapQ(ctx, clt, obj, "/metadata/annotations", obj.Annotations, map[string]string{
				"annotation_key1": "annotation_value1",
				"annotation_key2": "annotation_value2",
			})()).To(Succeed())

			Expect(clt.Get(ctx, types.NamespacedName{Name: obj.Name}, obj)).To(Succeed())
			Expect(obj.Annotations["annotation_key1"]).To(Equal("annotation_value1"))
			Expect(obj.Annotations["annotation_key2"]).To(Equal("annotation_value2"))
		})

		It("should replace existing members, ignore others, and add new ones", func(ctx SpecContext) {
			obj := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dummy-map-obj-3",
					Labels: map[string]string{
						"ignore_me_key":  "ignore_me_value",
						"dont_ignore_me": "replace",
					},
				},
			}

			Expect(clt.Create(ctx, obj)).To(Succeed())

			Expect(JsonPatchMapQ(ctx, clt, obj, "/metadata/labels", obj.Labels, map[string]string{
				"dont_ignore_me": "i_got_you",
				"add_me":         "added",
			})()).To(Succeed())

			Expect(clt.Get(ctx, types.NamespacedName{Name: obj.Name}, obj)).To(Succeed())
			Expect(obj.Labels["ignore_me_key"]).To(Equal("ignore_me_value"))
			Expect(obj.Labels["dont_ignore_me"]).To(Equal("i_got_you"))
			Expect(obj.Labels["add_me"]).To(Equal("added"))
		})
	})

	Context("a finalizer in", func() {
		It("should work when no other finalizers exist", func(ctx SpecContext) {
			obj := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dummy-finin-obj-1",
				},
			}

			Expect(clt.Create(ctx, obj)).To(Succeed())

			Expect(JsonPatchFinalizerInQ(ctx, clt, obj, "add-my/finalizer")()).To(Succeed())

			Expect(clt.Get(ctx, types.NamespacedName{Name: obj.Name}, obj)).To(Succeed())
			Expect(obj.Finalizers).To(ContainElement("add-my/finalizer"))
		})

		It("should work with a non-empty finalizer list", func(ctx SpecContext) {
			obj := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "dummy-finin-obj-2",
					Finalizers: []string{"already-existing/finalizer"},
				},
			}

			Expect(clt.Create(ctx, obj)).To(Succeed())

			Expect(JsonPatchFinalizerInQ(ctx, clt, obj, "add-my/finalizer")()).To(Succeed())

			Expect(clt.Get(ctx, types.NamespacedName{Name: obj.Name}, obj)).To(Succeed())
			Expect(obj.Finalizers).To(ContainElement("already-existing/finalizer"))
			Expect(obj.Finalizers).To(ContainElement("add-my/finalizer"))
		})

		It("should work with an existing finalizer", func(ctx SpecContext) {
			obj := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "dummy-finin-obj-3",
					Finalizers: []string{"already-existing/finalizer"},
				},
			}

			Expect(clt.Create(ctx, obj)).To(Succeed())

			Expect(JsonPatchFinalizerInQ(ctx, clt, obj, "already-existing/finalizer")()).To(Succeed())

			Expect(clt.Get(ctx, types.NamespacedName{Name: obj.Name}, obj)).To(Succeed())
			Expect(obj.Finalizers).To(ContainElement("already-existing/finalizer"))
		})
	})

	Context("a finalizer out", func() {
		It("should work when only one finalizer exists", func(ctx SpecContext) {
			obj := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "dummy-finout-obj-1",
					Finalizers: []string{"remove-this/finalizer"},
				},
			}

			Expect(clt.Create(ctx, obj)).To(Succeed())

			Expect(JsonPatchFinalizerOutQ(ctx, clt, obj, "remove-this/finalizer")()).To(Succeed())

			Expect(clt.Get(ctx, types.NamespacedName{Name: obj.Name}, obj)).To(Succeed())
			Expect(obj.Finalizers).To(BeEmpty())
		})

		It("should only remove the required finalizer when multiple exist", func(ctx SpecContext) {
			obj := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "dummy-finout-obj-2",
					Finalizers: []string{"remove-this/finalizer", "do-not-remove-this/finalizer"},
				},
			}

			Expect(clt.Create(ctx, obj)).To(Succeed())

			Expect(JsonPatchFinalizerOutQ(ctx, clt, obj, "remove-this/finalizer")()).To(Succeed())

			Expect(clt.Get(ctx, types.NamespacedName{Name: obj.Name}, obj)).To(Succeed())
			Expect(obj.Finalizers).To(ContainElement("do-not-remove-this/finalizer"))
			Expect(obj.Finalizers).ToNot(ContainElement("remove-this/finalizer"))
		})

		It("should return an error when attempting to remove a non-existing finalizer", func(ctx SpecContext) {
			obj := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "dummy-finout-obj-3",
					Finalizers: []string{"do-not-remove-this/finalizer"},
				},
			}

			Expect(clt.Create(ctx, obj)).To(Succeed())

			_, _, err := JsonPatchFinalizerOut(ctx, clt, obj, "remove-this/finalizer")
			Expect(err).To(MatchError("finalizer index not found"))
		})
	})

	Context("a spec", func() {
		It("should work with any spec", func(ctx SpecContext) {
			obj := &testdata.DummyCRD{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dummy-spec-obj-1",
				},
				Spec: testdata.DummyCRDSpec{
					FirstDummyValue:  "replace-me",
					SecondDummyValue: "leave-me-alone",
				},
			}

			Expect(clt.Create(ctx, obj)).To(Succeed())

			newSpec := &obj.Spec
			newSpec.FirstDummyValue = "replaced"

			Expect(JsonPatchSpecQ(ctx, clt, obj, newSpec)()).To(Succeed())

			Expect(clt.Get(ctx, types.NamespacedName{Name: obj.Name}, obj)).To(Succeed())
			Expect(obj.Spec.FirstDummyValue).To(Equal("replaced"))
			Expect(obj.Spec.SecondDummyValue).To(Equal("leave-me-alone"))
		})
	})
})
