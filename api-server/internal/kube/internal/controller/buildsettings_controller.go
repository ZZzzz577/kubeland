/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	appv1 "api-server/internal/kube/api/v1"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BuildSettingsReconciler reconciles a BuildSettings object
type BuildSettingsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// SetupWithManager sets up the controller with the Manager.
func (r *BuildSettingsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.BuildSettings{}).
		Owns(&v1.ConfigMap{}).
		Named("buildsettings").
		Complete(r)
}

// +kubebuilder:rbac:groups=app.kubeland.com,resources=buildsettings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=app.kubeland.com,resources=buildsettings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=app.kubeland.com,resources=buildsettings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the BuildSettings object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *BuildSettingsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var buildSettings appv1.BuildSettings
	if err := r.Get(ctx, req.NamespacedName, &buildSettings); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		} else {
			log.Error().Err(err).Msg("get build settings error")
			return ctrl.Result{}, err
		}
	}
	// 检查是否创建了dockerfile configmap
	var dockerfileConfigMap v1.ConfigMap
	cmName := fmt.Sprintf("%s-dockerfile-cm", buildSettings.Name)
	cmObjKey := client.ObjectKey{
		Namespace: req.Namespace,
		Name:      cmName,
	}
	err := r.Get(ctx, cmObjKey, &dockerfileConfigMap)
	if err != nil && !errors.IsNotFound(err) {
		log.Error().Err(err).Msg("get build settings error")
		return ctrl.Result{}, err
	}

	if errors.IsNotFound(err) {
		dockerfileConfigMap = v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cmObjKey.Name,
				Namespace: cmObjKey.Namespace,
			},
			Data: map[string]string{
				"Dockerfile": buildSettings.Spec.Dockerfile,
			},
		}

		if err = ctrl.SetControllerReference(&buildSettings, &dockerfileConfigMap, r.Scheme); err != nil {
			log.Error().Err(err).Msg("set owner reference error")
			return ctrl.Result{}, err
		}

		err = r.Create(ctx, &dockerfileConfigMap)
		if err != nil {
			log.Error().Err(err).Msg("create build settings dockerfile configmap error")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil

	} else {
		shouldUpdate := false

		if dockerfileConfigMap.Data["Dockerfile"] != buildSettings.Spec.Dockerfile {
			dockerfileConfigMap.Data["Dockerfile"] = buildSettings.Spec.Dockerfile
			shouldUpdate = true
		}

		if shouldUpdate {
			if err = ctrl.SetControllerReference(&buildSettings, &dockerfileConfigMap, r.Scheme); err != nil {
				log.Error().Err(err).Msg("set owner reference error")
				return ctrl.Result{}, err
			}

			if err = r.Update(ctx, &dockerfileConfigMap); err != nil {
				log.Error().Err(err).Msg("update build settings dockerfile configmap error")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
}
