package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// +groupName=test.dummy.cr

type (
	DummyCRDSpec struct {
		FirstDummyValue  string `json:"firstDummyValue,omitempty"`
		SecondDummyValue string `json:"secondDummyValue,omitempty"`
	}

	// +kubebuilder:object:root=true
	// +kubebuilder:resource:scope=Cluster
	DummyCRD struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`
		Spec              DummyCRDSpec `json:"spec,omitempty"`
	}

	// +kubebuilder:object:root=true
	DummyCRDList struct {
		metav1.TypeMeta `json:",inline"`
		metav1.ListMeta `json:"metadata,omitempty"`
		Items           []DummyCRD `json:"items"`
	}
)

func InitTestApi(s *runtime.Scheme) error {
	gv := schema.GroupVersion{Group: "test.dummy.cr", Version: "v1"}
	builder := &scheme.Builder{GroupVersion: gv}

	builder.Register(&DummyCRD{}, &DummyCRDList{})

	return builder.AddToScheme(s)
}
