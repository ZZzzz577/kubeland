package controller

import (
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"api-server/internal/data/generated/application"
	apiv1 "api-server/internal/kube/api/v1"
	"context"
	"github.com/rs/zerolog/log"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ApplicationReconciler struct {
	clusterId uint64
	client    client.Client
	db        *data.Data
}

func (a *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1.Application{}).
		Named("application").
		Complete(a)
}

func (a *ApplicationReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	app := &apiv1.Application{}
	if err := a.client.Get(ctx, request.NamespacedName, app); err != nil {
		log.Error().Err(err).Msg("get application error")
		return ctrl.Result{}, err
	}
	err := a.db.WithTx(ctx, func(tx *generated.Tx) error {
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
