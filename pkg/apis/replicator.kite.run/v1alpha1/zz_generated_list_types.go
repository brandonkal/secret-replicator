/*
Copyright 2020 Brandon Kalinowski
Boilerplate is Copyright 2019 Wrangler Sample Controller Authors (Apache 2.0)
*/

// Code generated by wrangler. DO NOT EDIT.

// +k8s:deepcopy-gen=package
// +groupName=replicator.kite.run
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretExportList is a list of SecretExport resources
type SecretExportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SecretExport `json:"items"`
}

func NewSecretExport(namespace, name string, obj SecretExport) *SecretExport {
	obj.APIVersion, obj.Kind = SchemeGroupVersion.WithKind("SecretExport").ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretRequestList is a list of SecretRequest resources
type SecretRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SecretRequest `json:"items"`
}

func NewSecretRequest(namespace, name string, obj SecretRequest) *SecretRequest {
	obj.APIVersion, obj.Kind = SchemeGroupVersion.WithKind("SecretRequest").ToAPIVersionAndKind()
	obj.Name = name
	obj.Namespace = namespace
	return &obj
}
