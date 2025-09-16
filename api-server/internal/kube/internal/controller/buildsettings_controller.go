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
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"api-server/internal/data/generated/gitrepo"
	"api-server/internal/data/generated/imagerepo"
	appv1 "api-server/internal/kube/api/v1"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
)

// BuildSettingsReconciler reconciles a BuildSettings object
type BuildSettingsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Data   *data.Data
}

// SetupWithManager sets up the controller with the Manager.
func (r *BuildSettingsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.BuildSettings{}).
		Owns(&corev1.Secret{}).
		Owns(&corev1.ConfigMap{}).
		Named("buildsettings").
		WithOptions(controller.TypedOptions[ctrl.Request]{
			SkipNameValidation: lo.ToPtr(true),
		}).
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

	if err := r.ApplyGit(ctx, &buildSettings); err != nil {
		log.Error().Err(err).Msg("apply git error")
		return ctrl.Result{}, err
	}

	if err := r.ApplyImage(ctx, &buildSettings); err != nil {
		log.Error().Err(err).Msg("apply image error")
		return ctrl.Result{}, err
	}

	if err := r.ApplyDockerfile(ctx, &buildSettings); err != nil {
		log.Error().Err(err).Msg("apply dockerfile error")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil

}

func (r *BuildSettingsReconciler) ApplyGit(ctx context.Context, buildSettings *appv1.BuildSettings) error {
	gitSettings := buildSettings.Spec.Git
	gitRepo, err := r.Data.GitRepo.Query().
		Where(gitrepo.Name(gitSettings.RepoName)).
		Only(ctx)
	if generated.IsNotFound(err) {
		return nil
	}
	if err != nil {
		log.Error().Err(err).Msg("get git repo error")
		return err
	}

	var gitSecret corev1.Secret
	secretName := fmt.Sprintf("%s-git", buildSettings.Name)
	secretObjKey := client.ObjectKey{
		Namespace: buildSettings.Namespace,
		Name:      secretName,
	}
	err = r.Get(ctx, secretObjKey, &gitSecret)
	if err != nil && !errors.IsNotFound(err) {
		log.Error().Err(err).Msg("get git secret error")
		return err
	}

	const (
		GitTokenKey = "GIT_TOKEN"
	)
	if errors.IsNotFound(err) {
		gitSecret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretObjKey.Name,
				Namespace: secretObjKey.Namespace,
			},
			Data: map[string][]byte{
				GitTokenKey: []byte(gitRepo.Token),
			},
		}

		if err = ctrl.SetControllerReference(buildSettings, &gitSecret, r.Scheme); err != nil {
			log.Error().Err(err).Msg("set owner reference error")
			return err
		}

		err = r.Create(ctx, &gitSecret)
		if err != nil {
			log.Error().Err(err).Msg("create git secret error")
			return err
		}
		return nil

	} else {
		shouldUpdate := false
		if string(gitSecret.Data[GitTokenKey]) != gitRepo.Token {
			gitSecret.Data[GitTokenKey] = []byte(gitRepo.Token)
			shouldUpdate = true
		}

		if shouldUpdate {
			if err = r.Update(ctx, &gitSecret); err != nil {
				log.Error().Err(err).Msg("update git secret error")
				return err
			}

			if err = ctrl.SetControllerReference(buildSettings, &gitSecret, r.Scheme); err != nil {
				log.Error().Err(err).Msg("set owner reference error")
				return err
			}
		}
		return nil
	}
}

func (r *BuildSettingsReconciler) ApplyImage(ctx context.Context, buildSettings *appv1.BuildSettings) error {
	imageSettings := buildSettings.Spec.Image
	imageRepo, err := r.Data.ImageRepo.Query().
		Where(imagerepo.Name(imageSettings.RepoName)).
		Only(ctx)
	if generated.IsNotFound(err) {
		return nil
	}
	if err != nil {
		log.Error().Err(err).Msg("get image repo error")
		return err
	}

	var imageSecret corev1.Secret
	secretName := fmt.Sprintf("%s-image", buildSettings.Name)
	secretObjKey := client.ObjectKey{
		Namespace: buildSettings.Namespace,
		Name:      secretName,
	}
	err = r.Get(ctx, secretObjKey, &imageSecret)
	if err != nil && !errors.IsNotFound(err) {
		log.Error().Err(err).Msg("get image secret error")
		return err
	}

	const (
		ImageUrlKey      = "IMAGE_URL"
		ImageUsernameKey = "IMAGE_USERNAME"
		ImagePasswordKey = "IMAGE_PASSWORD"
	)
	if errors.IsNotFound(err) {
		imageSecret = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretObjKey.Name,
				Namespace: secretObjKey.Namespace,
			},
			Data: map[string][]byte{
				ImageUrlKey:      []byte(imageRepo.URL),
				ImageUsernameKey: []byte(imageRepo.Username),
				ImagePasswordKey: []byte(imageRepo.Password),
			},
		}

		if err = ctrl.SetControllerReference(buildSettings, &imageSecret, r.Scheme); err != nil {
			log.Error().Err(err).Msg("set owner reference error")
			return err
		}

		err = r.Create(ctx, &imageSecret)
		if err != nil {
			log.Error().Err(err).Msg("create image secret error")
			return err
		}
		return nil

	} else {
		shouldUpdate := false
		if string(imageSecret.Data[ImageUrlKey]) != imageRepo.URL {
			imageSecret.Data[ImageUrlKey] = []byte(imageRepo.URL)
			shouldUpdate = true
		}
		if string(imageSecret.Data[ImageUsernameKey]) != imageRepo.Username {
			imageSecret.Data[ImageUsernameKey] = []byte(imageRepo.Username)
			shouldUpdate = true
		}
		if string(imageSecret.Data[ImagePasswordKey]) != imageRepo.Password {
			imageSecret.Data[ImagePasswordKey] = []byte(imageRepo.Password)
			shouldUpdate = true
		}

		if shouldUpdate {
			if err = r.Update(ctx, &imageSecret); err != nil {
				log.Error().Err(err).Msg("update image secret error")
				return err
			}

			if err = ctrl.SetControllerReference(buildSettings, &imageSecret, r.Scheme); err != nil {
				log.Error().Err(err).Msg("set owner reference error")
				return err
			}
		}
		return nil
	}
}

func (r *BuildSettingsReconciler) ApplyDockerfile(ctx context.Context, buildSettings *appv1.BuildSettings) error {
	var dockerfileConfigMap corev1.ConfigMap
	cmName := fmt.Sprintf("%s-dockerfile", buildSettings.Name)
	cmObjKey := client.ObjectKey{
		Namespace: buildSettings.Namespace,
		Name:      cmName,
	}
	err := r.Get(ctx, cmObjKey, &dockerfileConfigMap)
	if err != nil && !errors.IsNotFound(err) {
		log.Error().Err(err).Msg("get dockerfile configmap error")
		return err
	}

	if errors.IsNotFound(err) {
		dockerfileConfigMap = corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cmObjKey.Name,
				Namespace: cmObjKey.Namespace,
			},
			Data: map[string]string{
				"Dockerfile": buildSettings.Spec.Dockerfile,
			},
		}

		if err = ctrl.SetControllerReference(buildSettings, &dockerfileConfigMap, r.Scheme); err != nil {
			log.Error().Err(err).Msg("set owner reference error")
			return err
		}

		err = r.Create(ctx, &dockerfileConfigMap)
		if err != nil {
			log.Error().Err(err).Msg("create build settings dockerfile configmap error")
			return err
		}
		return nil

	} else {
		shouldUpdate := false
		if dockerfileConfigMap.Data["Dockerfile"] != buildSettings.Spec.Dockerfile {
			dockerfileConfigMap.Data["Dockerfile"] = buildSettings.Spec.Dockerfile
			shouldUpdate = true
		}

		if shouldUpdate {
			if err = r.Update(ctx, &dockerfileConfigMap); err != nil {
				log.Error().Err(err).Msg("update build settings dockerfile configmap error")
				return err
			}

			if err = ctrl.SetControllerReference(buildSettings, &dockerfileConfigMap, r.Scheme); err != nil {
				log.Error().Err(err).Msg("set owner reference error")
				return err
			}
		}
		return nil
	}
}
