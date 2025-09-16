package kube

import (
	"api-server/internal/data"
	apiv1 "api-server/internal/kube/api/v1"
	"api-server/internal/kube/internal/controller"
	ctrl "sigs.k8s.io/controller-runtime"
)

func RegisterControllers(
	mgr ctrl.Manager,
	data *data.Data,
) error {
	if err := apiv1.AddToScheme(mgr.GetScheme()); err != nil {
		return err
	}
	return (&controller.BuildSettingsReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Data:   data,
	}).SetupWithManager(mgr)
}
