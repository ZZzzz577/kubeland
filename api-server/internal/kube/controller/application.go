package controller

import (
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"api-server/internal/data/generated/application"
	appv1 "api-server/internal/kube/api/v1"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ApplicationReconciler struct {
	clusterId uint64
	client    client.Client
	db        *data.Data
}

func (a *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.Application{}).
		Named("application").
		WithOptions(controller.TypedOptions[ctrl.Request]{
			SkipNameValidation: lo.ToPtr(true),
		}).
		Complete(a)
}

// +kubebuilder:rbac:groups=app.kubeland.com,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.kubeland.com,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.kubeland.com,resources=applications/finalizers,verbs=update

func (a *ApplicationReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	app := &appv1.Application{}
	err := a.client.Get(ctx, request.NamespacedName, app)
	if errors.IsNotFound(err) {
		_, err = a.db.Application.Delete().Where(
			application.ClusterID(a.clusterId),
			application.Name(app.Name),
		).Exec(ctx)
		if err != nil {
			log.Error().Err(err).Msg("delete application error")
			return ctrl.Result{}, err
		}
	}
	if err != nil {
		log.Error().Err(err).Msg("get application error")
		return ctrl.Result{}, err
	}

	err = a.db.WithTx(ctx, func(tx *generated.Tx) error {
		exist, err := tx.Application.Query().Where(
			application.ClusterID(a.clusterId),
			application.Name(app.Name),
		).Exist(ctx)
		if err != nil {
			return err
		}
		if !exist {
			if err = tx.Application.Create().
				SetName(app.Name).
				SetDescription(app.Spec.Description).
				SetClusterID(a.clusterId).
				Exec(ctx); err != nil {
				return err
			}
		} else {
			if err = tx.Application.Update().
				SetDescription(app.Spec.Description).
				Where(
					application.ClusterID(a.clusterId),
					application.Name(app.Name)).
				Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("reconcile application error")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
