// Package v1beta1 contains the input type for this Function
// +kubebuilder:object:generate=true
// +groupName=starlark.fn.crossplane.io
// +versionName=v1beta1
package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// This isn't a custom resource, in the sense that we never install its CRD.
// It is a KRM-like object, so we generate a CRD to describe its schema.

// A ScriptSource is a source from which a script can be loaded.
type ScriptSource string

// Supported script sources.
const (
	// ScriptSourceInline specifies a script inline.
	ScriptSourceInline ScriptSource = "Inline"
)

// Script can be used to provide input to this Function.
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:categories=crossplane
type Script struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Source of this script. Currently only Inline is supported.
	// +kubebuilder:validation:Enum=Inline
	// +kubebuilder:default=Inline
	Source ScriptSource `json:"source"`

	// Inline specifies a script inline
	// +optional
	Inline *string `json:"inline,omitempty"`
}
