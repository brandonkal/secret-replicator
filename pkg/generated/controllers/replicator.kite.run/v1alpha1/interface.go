/*
Copyright 2020 Brandon Kalinowski
Boilerplate is Copyright 2019 Wrangler Sample Controller Authors (Apache 2.0)
*/

// Code generated by wrangler. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/brandonkal/secret-replicator/pkg/apis/replicator.kite.run/v1alpha1"
	clientset "github.com/brandonkal/secret-replicator/pkg/generated/clientset/versioned/typed/replicator.kite.run/v1alpha1"
	informers "github.com/brandonkal/secret-replicator/pkg/generated/informers/externalversions/replicator.kite.run/v1alpha1"
	"github.com/rancher/wrangler/pkg/generic"
)

type Interface interface {
	SecretExport() SecretExportController
	SecretRequest() SecretRequestController
}

func New(controllerManager *generic.ControllerManager, client clientset.ReplicatorV1alpha1Interface,
	informers informers.Interface) Interface {
	return &version{
		controllerManager: controllerManager,
		client:            client,
		informers:         informers,
	}
}

type version struct {
	controllerManager *generic.ControllerManager
	informers         informers.Interface
	client            clientset.ReplicatorV1alpha1Interface
}

func (c *version) SecretExport() SecretExportController {
	return NewSecretExportController(v1alpha1.SchemeGroupVersion.WithKind("SecretExport"), c.controllerManager, c.client, c.informers.SecretExports())
}
func (c *version) SecretRequest() SecretRequestController {
	return NewSecretRequestController(v1alpha1.SchemeGroupVersion.WithKind("SecretRequest"), c.controllerManager, c.client, c.informers.SecretRequests())
}
