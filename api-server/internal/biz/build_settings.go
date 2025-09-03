package biz

import (
	"api-server/api/v1/application"
	settings "api-server/api/v1/build_settings"
	appv1 "api-server/internal/kube/api/v1"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
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

func (b *BuildSettingsBiz) GetBuildSettings(ctx context.Context, request *application.IdentityRequest) (*settings.BuildSettings, error) {
	appName := request.GetName()
	client, err := b.cm.GetClient(ctx, appName)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return nil, err
	}

	namespace := "default"
	var buildSettings appv1.BuildSettings
	err = client.Get(ctx, cr.ObjectKey{
		Namespace: namespace,
		Name:      appName,
	}, &buildSettings)

	if errors.IsNotFound(err) {
		return &settings.BuildSettings{
			Dockerfile: "",
		}, nil
	}
	if err != nil {
		log.Error().Err(err).Msg("get build settings error")
		return nil, err
	}

	return &settings.BuildSettings{
		Git: &settings.BuildSettings_GitSettings{
			Url: buildSettings.Spec.Git.Url,
		},
		Image: &settings.BuildSettings_ImageSettings{
			Url: buildSettings.Spec.Image.Url,
		},
		Dockerfile: buildSettings.Spec.Dockerfile,
	}, nil
}

func (b *BuildSettingsBiz) ApplyBuildSettings(ctx context.Context, request *settings.ApplyBuildSettingsRequest) error {
	appName := request.GetName()
	buildSettings := request.GetBuildSettings()
	gitSettings := buildSettings.GetGit()
	imageSettings := buildSettings.GetImage()
	dockerfile := buildSettings.GetDockerfile()

	client, err := b.cm.GetClient(ctx, appName)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return err
	}

	name := fmt.Sprintf("%s", appName)
	namespace := "default"
	var kubeBuildSettings appv1.BuildSettings
	err = client.Get(ctx, cr.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, &kubeBuildSettings)
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
				Git: appv1.GitSettings{
					Url: gitSettings.GetUrl(),
				},
				Image: appv1.ImageSettings{
					Url: imageSettings.GetUrl(),
				},
				Dockerfile: dockerfile,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("create build settings error")
			return err
		}
	} else {
		kubeBuildSettings.Spec.Git.Url = gitSettings.GetUrl()
		kubeBuildSettings.Spec.Image.Url = imageSettings.GetUrl()
		kubeBuildSettings.Spec.Dockerfile = dockerfile
		err = client.Update(ctx, &kubeBuildSettings)
		if err != nil {
			log.Error().Err(err).Msg("update build settings error")
			return err
		}
	}

	return nil

}
