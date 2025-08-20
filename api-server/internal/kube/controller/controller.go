package controller

import (
	"api-server/internal/data"
	apiv1 "api-server/internal/kube/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func NewControllers(
	clusterId uint64,
	mgr ctrl.Manager,
	db *data.Data,
) error {
	if err := apiv1.AddToScheme(mgr.GetScheme()); err != nil {
		return err
	}
	return (&ApplicationReconciler{
		clusterId: clusterId,
		client:    mgr.GetClient(),
		db:        db,
	}).SetupWithManager(mgr)
}
