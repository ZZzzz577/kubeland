package biz

import (
	settings "api-server/api/v1/build_settings"
	appv1 "api-server/internal/kube/api/v1"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
)

type BuildSettingsBiz struct {
	cm *ClusterManagers
}

func NewBuildSettingsBiz(
	cm *ClusterManagers,
) *BuildSettingsBiz {
	return &BuildSettingsBiz{
		cm: cm,
	}
}

func (b *BuildSettingsBiz) GetBuildSettings(ctx context.Context, request *settings.IdRequest) (*settings.BuildSettings, error) {
	appId := request.GetApplicationId()
	client, err := b.cm.GetClient(ctx, request.GetApplicationId())
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return nil, err
	}

	name := fmt.Sprintf("build-settings-%d", appId)
	namespace := "default"
	var buildSettings appv1.BuildSettings
	err = client.Get(ctx, cr.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, &buildSettings)

	if errors.IsNotFound(err) {
		return &settings.BuildSettings{
			ApplicationId: appId,
			Dockerfile:    "",
		}, nil
	}
	if err != nil {
		log.Error().Err(err).Msg("get build settings error")
		return nil, err
	}

	return &settings.BuildSettings{
		ApplicationId: appId,
		Dockerfile:    buildSettings.Spec.Dockerfile.Data["dockerfile"],
	}, nil
}

func (b *BuildSettingsBiz) ApplyBuildSettings(ctx context.Context, request *settings.BuildSettings) error {
	appId := request.GetApplicationId()
	client, err := b.cm.GetClient(ctx, request.GetApplicationId())
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return err
	}

	name := fmt.Sprintf("build-settings-%d", appId)
	namespace := "default"
	var buildSettings appv1.BuildSettings
	err = client.Get(ctx, cr.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, &buildSettings, &cr.GetOptions{
		Raw: &metav1.GetOptions{},
	})
	if err != nil && !errors.IsNotFound(err) {
		log.Error().Err(err).Msg("get build settings error")
		return err
	}

	if errors.IsNotFound(err) {
		err = client.Create(ctx, &appv1.BuildSettings{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
			Spec: appv1.BuildSettingsSpec{
				Dockerfile: v1.ConfigMap{
					Data: map[string]string{
						"dockerfile": request.GetDockerfile(),
					},
				},
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("create build settings error")
			return err
		}
	} else {
		buildSettings.Spec.Dockerfile.Data["dockerfile"] = request.GetDockerfile()
		err = client.Update(ctx, &buildSettings)
		if err != nil {
			log.Error().Err(err).Msg("update build settings error")
			return err
		}
	}

	return nil

}
