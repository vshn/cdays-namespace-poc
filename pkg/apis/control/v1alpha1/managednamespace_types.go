//
// Copyright (c) 2019, VSHN AG, info@vshn.ch
// Licensed under "BSD 3-Clause". See LICENSE file.
//
//

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// ManagedNamespaceSpec defines the desired state of ManagedNamespace
// +k8s:openapi-gen=true
type ManagedNamespaceSpec struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// Description may show further (human readable) information
	Description string `json:"description,omitempty"`
}

// ManagedNamespaceStatus defines the observed state of ManagedNamespace
// +k8s:openapi-gen=true
type ManagedNamespaceStatus struct {
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file

	// CreatedNamespace references the UID of the created namespace object
	CreatedNamespace types.UID `json:"createdNamespace,omitempty"`
	// Phase is the current lifecycle phase of the ManagedNamespace.
	// +optional
	Phase corev1.NamespacePhase `json:"phase,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedNamespace is the Schema for the managednamespaces API
// +k8s:openapi-gen=true
type ManagedNamespace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ManagedNamespaceSpec   `json:"spec,omitempty"`
	Status ManagedNamespaceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedNamespaceList contains a list of ManagedNamespace
type ManagedNamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ManagedNamespace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ManagedNamespace{}, &ManagedNamespaceList{})
}
