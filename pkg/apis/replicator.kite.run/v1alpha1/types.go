/*
Copyright 2010 Brandon Kalinowski
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretExport is a specification for a SecretExport resource
type SecretExport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretExportSpec   `json:"spec"`
	Status SecretExportStatus `json:"status"`
}

// SecretExportSpec is the spec for a SecretExport resource
type SecretExportSpec struct {
	DeploymentName string   `json:"deploymentName"`
	Replicas       *int32   `json:"replicas"`
	ToNamespaces   []string `json:"toNamespaces,omitempty"`
}

// SecretExportStatus is the status for a SecretExport resource
type SecretExportStatus struct {
	AvailableReplicas             int32  `json:"availableReplicas"`
	ObservedSecretResourceVersion string `json:"observedSecretResourceVersion,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretRequest specifies a request to import a secret into a namespace.
type SecretRequest struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretRequestSpec   `json:"spec"`
	Status SecretRequestStatus `json:"status"`
}

type SecretRequestSpec struct {
	FromNamespace string `json:"fromNamespace,omitempty"`
}

type SecretRequestStatus struct {
	// NOTE: added string type here.
	GenericStatus string `json:",inline"`
}
