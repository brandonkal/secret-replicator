/*
Copyright 2020 Brandon Kalinowski
Boilerplate is Copyright 2019 Wrangler Sample Controller Authors (Apache 2.0)
*/

// Code generated by wrangler. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/brandonkal/secret-replicator/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// SecretExports returns a SecretExportInformer.
	SecretExports() SecretExportInformer
	// SecretRequests returns a SecretRequestInformer.
	SecretRequests() SecretRequestInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// SecretExports returns a SecretExportInformer.
func (v *version) SecretExports() SecretExportInformer {
	return &secretExportInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// SecretRequests returns a SecretRequestInformer.
func (v *version) SecretRequests() SecretRequestInformer {
	return &secretRequestInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
