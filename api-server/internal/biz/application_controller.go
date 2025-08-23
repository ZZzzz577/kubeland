package biz

import (
	"api-server/internal/kube/informers/externalversions"
	applicationv1 "api-server/internal/kube/listers/application/v1"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
	"k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type ApplicationController struct {
	clusterId uint64
	synced    cache.InformerSynced
	listener  applicationv1.ApplicationLister
	workQueue workqueue.TypedRateLimitingInterface[cache.ObjectName]
}

func NewApplicationController(
	clusterId uint64,
	informerFactory externalversions.SharedInformerFactory,
) *ApplicationController {

	informer := informerFactory.Kubeland().V1().Applications()

	rateLimiter := workqueue.NewTypedMaxOfRateLimiter(
		workqueue.NewTypedItemExponentialFailureRateLimiter[cache.ObjectName](5*time.Millisecond, 1000*time.Second),
		&workqueue.TypedBucketRateLimiter[cache.ObjectName]{Limiter: rate.NewLimiter(rate.Limit(50), 300)},
	)
	workQueue := workqueue.NewTypedRateLimitingQueue(rateLimiter)

	controller := &ApplicationController{
		clusterId: clusterId,
		synced:    informer.Informer().HasSynced,
		listener:  informer.Lister(),
		workQueue: workQueue,
	}

	_, _ = informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueue,
		UpdateFunc: func(oldObj, newObj interface{}) {
			controller.enqueue(newObj)
		},
	})
	return controller
}

func (c *ApplicationController) enqueue(obj interface{}) {
	if objectRef, err := cache.ObjectToName(obj); err != nil {
		utilruntime.HandleError(err)
		return
	} else {
		c.workQueue.Add(objectRef)
	}
}

func (c *ApplicationController) Run(ctx context.Context, workers int) error {
	defer utilruntime.HandleCrash()
	defer c.workQueue.ShutDown()

	log.Info().Uint64("clusterId", c.clusterId).Msg("waiting for informer caches to sync")

	if !cache.WaitForCacheSync(ctx.Done(), c.synced) {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	log.Info().Uint64("clusterId", c.clusterId).Int("count", workers).Msg("starting workers")
	// Launch workers to process resources
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runWorker, time.Second)
	}

	log.Info().Uint64("clusterId", c.clusterId).Msg("started workers")
	<-ctx.Done()
	log.Info().Uint64("clusterId", c.clusterId).Msg("shutting down workers")

	return nil
}

func (c *ApplicationController) runWorker(ctx context.Context) {
	for c.processNextWorkItem(ctx) {
	}
}

func (c *ApplicationController) processNextWorkItem(ctx context.Context) bool {
	objRef, shutdown := c.workQueue.Get()
	if shutdown {
		return false
	}

	defer c.workQueue.Done(objRef)
	err := c.syncHandler(ctx, objRef)
	if err == nil {
		c.workQueue.Forget(objRef)
		return true
	}

	utilruntime.HandleErrorWithContext(ctx, err, "Error syncing; requeuing for later retry", "objectReference", objRef)
	c.workQueue.AddRateLimited(objRef)
	return true
}

func (c *ApplicationController) syncHandler(ctx context.Context, objectRef cache.ObjectName) error {

	app, err := c.listener.Applications(objectRef.Namespace).Get(objectRef.Name)
	if err != nil {
		if errors.IsNotFound(err) {
			utilruntime.HandleErrorWithContext(ctx, err, "Foo referenced by item in work queue no longer exists", "objectReference", objectRef)
			return nil
		}
		return err
	}
	fmt.Println(app)

	return nil
}
