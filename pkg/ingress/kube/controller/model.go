package controller

import (
	"errors"
	"time"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/kube/controllers"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/alibaba/higress/pkg/ingress/kube/util"
	. "github.com/alibaba/higress/pkg/ingress/log"
)

type Controller[lister any] interface {
	AddEventHandler(addOrUpdate func(util.ClusterNamespacedName), delete ...func(util.ClusterNamespacedName))

	Run(stop <-chan struct{})

	HasSynced() bool

	Lister() lister

	Informer() cache.SharedIndexInformer
}

type GetObjectFunc[lister any] func(lister, types.NamespacedName) (controllers.Object, error)

type CommonController[lister any] struct {
	typeName      string
	queue         workqueue.RateLimitingInterface
	informer      cache.SharedIndexInformer
	lister        lister
	updateHandler func(util.ClusterNamespacedName)
	removeHandler func(util.ClusterNamespacedName)
	getFunc       GetObjectFunc[lister]
	clusterId     string
}

func NewCommonController[lister any](typeName string, listerObj lister, informer cache.SharedIndexInformer,
	getFunc GetObjectFunc[lister], clusterId string) Controller[lister] {
	q := workqueue.NewRateLimitingQueue(workqueue.DefaultItemBasedRateLimiter())
	handler := controllers.LatestVersionHandlerFuncs(controllers.EnqueueForSelf(q))
	informer.AddEventHandler(handler)
	return &CommonController[lister]{
		typeName:  typeName,
		queue:     q,
		lister:    listerObj,
		informer:  informer,
		clusterId: clusterId,
		getFunc:   getFunc,
	}
}

func (c *CommonController[lister]) Lister() lister {
	return c.lister
}

func (c *CommonController[lister]) Informer() cache.SharedIndexInformer {
	return c.informer
}

func (c *CommonController[lister]) AddEventHandler(addOrUpdate func(util.ClusterNamespacedName), delete ...func(util.ClusterNamespacedName)) {
	c.updateHandler = addOrUpdate
	if len(delete) > 0 {
		c.removeHandler = delete[0]
	}
}

func (c *CommonController[lister]) Run(stop <-chan struct{}) {
	defer utilruntime.HandleCrash()
	defer c.queue.ShutDown()

	if !cache.WaitForCacheSync(stop, c.HasSynced) {
		IngressLog.Errorf("Failed to sync %s controller cache", c.typeName)
		return
	}
	go wait.Until(c.worker, time.Second, stop)
	<-stop
}

func (c *CommonController[lister]) worker() {
	for c.processNextWorkItem() {
	}
}

func (c *CommonController[lister]) processNextWorkItem() bool {
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)
	ingressNamespacedName := key.(types.NamespacedName)
	IngressLog.Debugf("%s %s push to queue", c.typeName, ingressNamespacedName)
	if err := c.onEvent(ingressNamespacedName); err != nil {
		IngressLog.Errorf("error processing %s item (%v) (retrying): %v", c.typeName, key, err)
		c.queue.AddRateLimited(key)
	} else {
		c.queue.Forget(key)
	}
	return true
}

func (c *CommonController[lister]) onEvent(namespacedName types.NamespacedName) error {
	if c.getFunc == nil {
		return errors.New("getFunc is nil")
	}
	obj := util.ClusterNamespacedName{
		NamespacedName: model.NamespacedName{
			Namespace: namespacedName.Namespace,
			Name:      namespacedName.Name,
		},
		ClusterId: c.clusterId,
	}
	_, err := c.getFunc(c.lister, namespacedName)
	if err != nil {
		if kerrors.IsNotFound(err) {
			if c.removeHandler == nil {
				return nil
			}
			c.removeHandler(obj)
		} else {
			return err
		}
	}

	c.updateHandler(obj)
	return nil
}

func (c *CommonController[lister]) HasSynced() bool {
	return c.informer.HasSynced()
}
