package main

import (
	"istio.io/api/networking/v1alpha3"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VirtualServiceTemplateSpec
type VirtualServiceTemplateSpec struct {

	// Metadata of the virtual services created from this template
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec indicates the behavior of a virtual service.
	// +optional
	Spec v1alpha3.VirtualService `json:"spec,omitempty"`
}
