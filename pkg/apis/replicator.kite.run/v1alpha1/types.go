/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
	ToNamespace    string   `json:"toNamespace,omitempty"`
	ToNamespaces   []string `json:"toNamespaces,omitempty"`
}

// SecretExportStatus is the status for a SecretExport resource
type SecretExportStatus struct {
	AvailableReplicas             int32  `json:"availableReplicas"`
	ObservedSecretResourceVersion string `json:"observedSecretResourceVersion,omitempty"`
}
