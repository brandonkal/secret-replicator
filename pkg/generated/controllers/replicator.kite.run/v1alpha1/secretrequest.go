/*
Copyright 2020 Brandon Kalinowski
Boilerplate is Copyright 2019 Wrangler Sample Controller Authors (Apache 2.0)
*/

// Code generated by wrangler. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/brandonkal/secret-replicator/pkg/apis/replicator.kite.run/v1alpha1"
	clientset "github.com/brandonkal/secret-replicator/pkg/generated/clientset/versioned/typed/replicator.kite.run/v1alpha1"
	informers "github.com/brandonkal/secret-replicator/pkg/generated/informers/externalversions/replicator.kite.run/v1alpha1"
	listers "github.com/brandonkal/secret-replicator/pkg/generated/listers/replicator.kite.run/v1alpha1"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type SecretRequestHandler func(string, *v1alpha1.SecretRequest) (*v1alpha1.SecretRequest, error)

type SecretRequestController interface {
	generic.ControllerMeta
	SecretRequestClient

	OnChange(ctx context.Context, name string, sync SecretRequestHandler)
	OnRemove(ctx context.Context, name string, sync SecretRequestHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() SecretRequestCache
}

type SecretRequestClient interface {
	Create(*v1alpha1.SecretRequest) (*v1alpha1.SecretRequest, error)
	Update(*v1alpha1.SecretRequest) (*v1alpha1.SecretRequest, error)
	UpdateStatus(*v1alpha1.SecretRequest) (*v1alpha1.SecretRequest, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.SecretRequest, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha1.SecretRequestList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SecretRequest, err error)
}

type SecretRequestCache interface {
	Get(namespace, name string) (*v1alpha1.SecretRequest, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.SecretRequest, error)

	AddIndexer(indexName string, indexer SecretRequestIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.SecretRequest, error)
}

type SecretRequestIndexer func(obj *v1alpha1.SecretRequest) ([]string, error)

type secretRequestController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.SecretRequestsGetter
	informer          informers.SecretRequestInformer
	gvk               schema.GroupVersionKind
}

func NewSecretRequestController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.SecretRequestsGetter, informer informers.SecretRequestInformer) SecretRequestController {
	return &secretRequestController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromSecretRequestHandlerToHandler(sync SecretRequestHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.SecretRequest
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.SecretRequest))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *secretRequestController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.SecretRequest))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateSecretRequestDeepCopyOnChange(client SecretRequestClient, obj *v1alpha1.SecretRequest, handler func(obj *v1alpha1.SecretRequest) (*v1alpha1.SecretRequest, error)) (*v1alpha1.SecretRequest, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *secretRequestController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *secretRequestController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *secretRequestController) OnChange(ctx context.Context, name string, sync SecretRequestHandler) {
	c.AddGenericHandler(ctx, name, FromSecretRequestHandlerToHandler(sync))
}

func (c *secretRequestController) OnRemove(ctx context.Context, name string, sync SecretRequestHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromSecretRequestHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *secretRequestController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *secretRequestController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), namespace, name, duration)
}

func (c *secretRequestController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *secretRequestController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *secretRequestController) Cache() SecretRequestCache {
	return &secretRequestCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *secretRequestController) Create(obj *v1alpha1.SecretRequest) (*v1alpha1.SecretRequest, error) {
	return c.clientGetter.SecretRequests(obj.Namespace).Create(obj)
}

func (c *secretRequestController) Update(obj *v1alpha1.SecretRequest) (*v1alpha1.SecretRequest, error) {
	return c.clientGetter.SecretRequests(obj.Namespace).Update(obj)
}

func (c *secretRequestController) UpdateStatus(obj *v1alpha1.SecretRequest) (*v1alpha1.SecretRequest, error) {
	return c.clientGetter.SecretRequests(obj.Namespace).UpdateStatus(obj)
}

func (c *secretRequestController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.SecretRequests(namespace).Delete(name, options)
}

func (c *secretRequestController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.SecretRequest, error) {
	return c.clientGetter.SecretRequests(namespace).Get(name, options)
}

func (c *secretRequestController) List(namespace string, opts metav1.ListOptions) (*v1alpha1.SecretRequestList, error) {
	return c.clientGetter.SecretRequests(namespace).List(opts)
}

func (c *secretRequestController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.SecretRequests(namespace).Watch(opts)
}

func (c *secretRequestController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SecretRequest, err error) {
	return c.clientGetter.SecretRequests(namespace).Patch(name, pt, data, subresources...)
}

type secretRequestCache struct {
	lister  listers.SecretRequestLister
	indexer cache.Indexer
}

func (c *secretRequestCache) Get(namespace, name string) (*v1alpha1.SecretRequest, error) {
	return c.lister.SecretRequests(namespace).Get(name)
}

func (c *secretRequestCache) List(namespace string, selector labels.Selector) ([]*v1alpha1.SecretRequest, error) {
	return c.lister.SecretRequests(namespace).List(selector)
}

func (c *secretRequestCache) AddIndexer(indexName string, indexer SecretRequestIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.SecretRequest))
		},
	}))
}

func (c *secretRequestCache) GetByIndex(indexName, key string) (result []*v1alpha1.SecretRequest, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.SecretRequest))
	}
	return result, nil
}

type SecretRequestStatusHandler func(obj *v1alpha1.SecretRequest, status v1alpha1.SecretRequestStatus) (v1alpha1.SecretRequestStatus, error)

type SecretRequestGeneratingHandler func(obj *v1alpha1.SecretRequest, status v1alpha1.SecretRequestStatus) ([]runtime.Object, v1alpha1.SecretRequestStatus, error)

func RegisterSecretRequestStatusHandler(ctx context.Context, controller SecretRequestController, condition condition.Cond, name string, handler SecretRequestStatusHandler) {
	statusHandler := &secretRequestStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromSecretRequestHandlerToHandler(statusHandler.sync))
}

func RegisterSecretRequestGeneratingHandler(ctx context.Context, controller SecretRequestController, apply apply.Apply,
	condition condition.Cond, name string, handler SecretRequestGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &secretRequestGeneratingHandler{
		SecretRequestGeneratingHandler: handler,
		apply:                          apply,
		name:                           name,
		gvk:                            controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	RegisterSecretRequestStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type secretRequestStatusHandler struct {
	client    SecretRequestClient
	condition condition.Cond
	handler   SecretRequestStatusHandler
}

func (a *secretRequestStatusHandler) sync(key string, obj *v1alpha1.SecretRequest) (*v1alpha1.SecretRequest, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	obj.Status = newStatus
	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(obj, "", nil)
		} else {
			a.condition.SetError(obj, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, obj.Status) {
		var newErr error
		obj, newErr = a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
	}
	return obj, err
}

type secretRequestGeneratingHandler struct {
	SecretRequestGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *secretRequestGeneratingHandler) Handle(obj *v1alpha1.SecretRequest, status v1alpha1.SecretRequestStatus) (v1alpha1.SecretRequestStatus, error) {
	objs, newStatus, err := a.SecretRequestGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	apply := a.apply

	if !a.opts.DynamicLookup {
		apply = apply.WithStrictCaching()
	}

	if !a.opts.AllowCrossNamespace && !a.opts.AllowClusterScoped {
		apply = apply.WithSetOwnerReference(true, false).
			WithDefaultNamespace(obj.GetNamespace()).
			WithListerNamespace(obj.GetNamespace())
	}

	if !a.opts.AllowClusterScoped {
		apply = apply.WithRestrictClusterScoped()
	}

	return newStatus, apply.
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
