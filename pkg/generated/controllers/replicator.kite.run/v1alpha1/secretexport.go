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

type SecretExportHandler func(string, *v1alpha1.SecretExport) (*v1alpha1.SecretExport, error)

type SecretExportController interface {
	generic.ControllerMeta
	SecretExportClient

	OnChange(ctx context.Context, name string, sync SecretExportHandler)
	OnRemove(ctx context.Context, name string, sync SecretExportHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() SecretExportCache
}

type SecretExportClient interface {
	Create(*v1alpha1.SecretExport) (*v1alpha1.SecretExport, error)
	Update(*v1alpha1.SecretExport) (*v1alpha1.SecretExport, error)
	UpdateStatus(*v1alpha1.SecretExport) (*v1alpha1.SecretExport, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.SecretExport, error)
	List(namespace string, opts metav1.ListOptions) (*v1alpha1.SecretExportList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SecretExport, err error)
}

type SecretExportCache interface {
	Get(namespace, name string) (*v1alpha1.SecretExport, error)
	List(namespace string, selector labels.Selector) ([]*v1alpha1.SecretExport, error)

	AddIndexer(indexName string, indexer SecretExportIndexer)
	GetByIndex(indexName, key string) ([]*v1alpha1.SecretExport, error)
}

type SecretExportIndexer func(obj *v1alpha1.SecretExport) ([]string, error)

type secretExportController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.SecretExportsGetter
	informer          informers.SecretExportInformer
	gvk               schema.GroupVersionKind
}

func NewSecretExportController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.SecretExportsGetter, informer informers.SecretExportInformer) SecretExportController {
	return &secretExportController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromSecretExportHandlerToHandler(sync SecretExportHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1alpha1.SecretExport
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1alpha1.SecretExport))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *secretExportController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1alpha1.SecretExport))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateSecretExportDeepCopyOnChange(client SecretExportClient, obj *v1alpha1.SecretExport, handler func(obj *v1alpha1.SecretExport) (*v1alpha1.SecretExport, error)) (*v1alpha1.SecretExport, error) {
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

func (c *secretExportController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *secretExportController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *secretExportController) OnChange(ctx context.Context, name string, sync SecretExportHandler) {
	c.AddGenericHandler(ctx, name, FromSecretExportHandlerToHandler(sync))
}

func (c *secretExportController) OnRemove(ctx context.Context, name string, sync SecretExportHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromSecretExportHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *secretExportController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, c.informer.Informer(), namespace, name)
}

func (c *secretExportController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controllerManager.EnqueueAfter(c.gvk, c.informer.Informer(), namespace, name, duration)
}

func (c *secretExportController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *secretExportController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *secretExportController) Cache() SecretExportCache {
	return &secretExportCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *secretExportController) Create(obj *v1alpha1.SecretExport) (*v1alpha1.SecretExport, error) {
	return c.clientGetter.SecretExports(obj.Namespace).Create(obj)
}

func (c *secretExportController) Update(obj *v1alpha1.SecretExport) (*v1alpha1.SecretExport, error) {
	return c.clientGetter.SecretExports(obj.Namespace).Update(obj)
}

func (c *secretExportController) UpdateStatus(obj *v1alpha1.SecretExport) (*v1alpha1.SecretExport, error) {
	return c.clientGetter.SecretExports(obj.Namespace).UpdateStatus(obj)
}

func (c *secretExportController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.SecretExports(namespace).Delete(name, options)
}

func (c *secretExportController) Get(namespace, name string, options metav1.GetOptions) (*v1alpha1.SecretExport, error) {
	return c.clientGetter.SecretExports(namespace).Get(name, options)
}

func (c *secretExportController) List(namespace string, opts metav1.ListOptions) (*v1alpha1.SecretExportList, error) {
	return c.clientGetter.SecretExports(namespace).List(opts)
}

func (c *secretExportController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.SecretExports(namespace).Watch(opts)
}

func (c *secretExportController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SecretExport, err error) {
	return c.clientGetter.SecretExports(namespace).Patch(name, pt, data, subresources...)
}

type secretExportCache struct {
	lister  listers.SecretExportLister
	indexer cache.Indexer
}

func (c *secretExportCache) Get(namespace, name string) (*v1alpha1.SecretExport, error) {
	return c.lister.SecretExports(namespace).Get(name)
}

func (c *secretExportCache) List(namespace string, selector labels.Selector) ([]*v1alpha1.SecretExport, error) {
	return c.lister.SecretExports(namespace).List(selector)
}

func (c *secretExportCache) AddIndexer(indexName string, indexer SecretExportIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1alpha1.SecretExport))
		},
	}))
}

func (c *secretExportCache) GetByIndex(indexName, key string) (result []*v1alpha1.SecretExport, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1alpha1.SecretExport))
	}
	return result, nil
}

type SecretExportStatusHandler func(obj *v1alpha1.SecretExport, status v1alpha1.SecretExportStatus) (v1alpha1.SecretExportStatus, error)

type SecretExportGeneratingHandler func(obj *v1alpha1.SecretExport, status v1alpha1.SecretExportStatus) ([]runtime.Object, v1alpha1.SecretExportStatus, error)

func RegisterSecretExportStatusHandler(ctx context.Context, controller SecretExportController, condition condition.Cond, name string, handler SecretExportStatusHandler) {
	statusHandler := &secretExportStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromSecretExportHandlerToHandler(statusHandler.sync))
}

func RegisterSecretExportGeneratingHandler(ctx context.Context, controller SecretExportController, apply apply.Apply,
	condition condition.Cond, name string, handler SecretExportGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &secretExportGeneratingHandler{
		SecretExportGeneratingHandler: handler,
		apply:                         apply,
		name:                          name,
		gvk:                           controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	RegisterSecretExportStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type secretExportStatusHandler struct {
	client    SecretExportClient
	condition condition.Cond
	handler   SecretExportStatusHandler
}

func (a *secretExportStatusHandler) sync(key string, obj *v1alpha1.SecretExport) (*v1alpha1.SecretExport, error) {
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

type secretExportGeneratingHandler struct {
	SecretExportGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *secretExportGeneratingHandler) Handle(obj *v1alpha1.SecretExport, status v1alpha1.SecretExportStatus) (v1alpha1.SecretExportStatus, error) {
	objs, newStatus, err := a.SecretExportGeneratingHandler(obj, status)
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
