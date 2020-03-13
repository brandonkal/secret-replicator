package main

import (
	"context"
	"fmt"

	samplev1alpha1 "github.com/brandonkal/secret-replicator/pkg/apis/replicator.kite.run/v1alpha1"
	samplescheme "github.com/brandonkal/secret-replicator/pkg/generated/clientset/versioned/scheme"
	"github.com/brandonkal/secret-replicator/pkg/generated/controllers/replicator.kite.run/v1alpha1"
	v1 "github.com/rancher/wrangler-api/pkg/generated/controllers/apps/v1"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

const controllerAgentName = "sample-controller"

const (
	// ErrResourceExists is used as part of the Event 'reason' when a SecretExport fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by SecretExport"
)

// Handler is the controller implementation for SecretExport resources
type Handler struct {
	deployments        v1.DeploymentClient
	deploymentsCache   v1.DeploymentCache
	secretExports      v1alpha1.SecretExportController
	secretExportsCache v1alpha1.SecretExportCache
	recorder           record.EventRecorder
}

// NewController returns a new sample controller
func Register(
	ctx context.Context,
	events typedcorev1.EventInterface,
	deployments v1.DeploymentController,
	secretExports v1alpha1.SecretExportController) {

	controller := &Handler{
		deployments:        deployments,
		deploymentsCache:   deployments.Cache(),
		secretExports:      secretExports,
		secretExportsCache: secretExports.Cache(),
		recorder:           buildEventRecorder(events),
	}

	// Register handlers
	deployments.OnChange(ctx, "secretExport-handler", controller.OnDeploymentChanged)
	secretExports.OnChange(ctx, "secretExport-handler", controller.OnSecretExportChanged)
}

func buildEventRecorder(events typedcorev1.EventInterface) record.EventRecorder {
	// Create event broadcaster
	// Add sample-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-controller types.
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	logrus.Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logrus.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: events})
	return eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})
}

func (h *Handler) OnSecretExportChanged(key string, secretExport *samplev1alpha1.SecretExport) (*samplev1alpha1.SecretExport, error) {
	// secretExport will be nil if key is deleted from cache
	if secretExport == nil {
		return nil, nil
	}

	deploymentName := secretExport.Spec.DeploymentName
	if deploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
		return nil, nil
	}

	// Get the deployment with the name specified in SecretExport.spec
	deployment, err := h.deploymentsCache.Get(secretExport.Namespace, deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = h.deployments.Create(newDeployment(secretExport))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return nil, err
	}

	// If the Deployment is not controlled by this SecretExport resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(deployment, secretExport) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		h.recorder.Event(secretExport, corev1.EventTypeWarning, ErrResourceExists, msg)
		// Notice we don't return an error here, this is intentional because an
		// error means we should retry to reconcile.  In this situation we've done all
		// we could, which was log an error.
		return nil, nil
	}

	// If this number of the replicas on the SecretExport resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if secretExport.Spec.Replicas != nil && *secretExport.Spec.Replicas != *deployment.Spec.Replicas {
		logrus.Infof("SecretExport %s replicas: %d, deployment replicas: %d", secretExport.Name, *secretExport.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = h.deployments.Update(newDeployment(secretExport))
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. THis could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return nil, err
	}

	// Finally, we update the status block of the SecretExport resource to reflect the
	// current state of the world
	err = h.updateSecretExportStatus(secretExport, deployment)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) updateSecretExportStatus(secretExport *samplev1alpha1.SecretExport, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	secretExportCopy := secretExport.DeepCopy()
	secretExportCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the SecretExport resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := h.secretExports.Update(secretExportCopy)
	return err
}

func (h *Handler) OnDeploymentChanged(key string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	// When an item is deleted the deployment is nil, just ignore
	if deployment == nil {
		return nil, nil
	}

	if ownerRef := metav1.GetControllerOf(deployment); ownerRef != nil {
		// If this object is not owned by a SecretExport, we should not do anything more
		// with it.
		if ownerRef.Kind != "SecretExport" {
			return nil, nil
		}

		secretExport, err := h.secretExportsCache.Get(deployment.Namespace, ownerRef.Name)
		if err != nil {
			logrus.Infof("ignoring orphaned object '%s' of secretExport '%s'", deployment.GetSelfLink(), ownerRef.Name)
			return nil, nil
		}

		h.secretExports.Enqueue(secretExport.Namespace, secretExport.Name)
		return nil, nil
	}

	return nil, nil
}

// newDeployment creates a new Deployment for a SecretExport resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the SecretExport resource that 'owns' it.
func newDeployment(secretExport *samplev1alpha1.SecretExport) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "nginx",
		"controller": secretExport.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretExport.Spec.DeploymentName,
			Namespace: secretExport.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(secretExport, schema.GroupVersionKind{
					Group:   samplev1alpha1.SchemeGroupVersion.Group,
					Version: samplev1alpha1.SchemeGroupVersion.Version,
					Kind:    "SecretExport",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: secretExport.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	}
}
