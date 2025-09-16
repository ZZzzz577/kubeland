package biz

import (
	"api-server/api/v1/application"
	settings "api-server/api/v1/build_settings"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"api-server/internal/data/generated/gitrepo"
	"api-server/internal/data/generated/imagerepo"
	appv1 "api-server/internal/kube/api/v1"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
)

type BuildSettingsBiz struct {
	cm   *ClusterManagers
	data *data.Data
}

func NewBuildSettingsBiz(
	cm *ClusterManagers,
	data *data.Data,
) *BuildSettingsBiz {
	return &BuildSettingsBiz{
		cm:   cm,
		data: data,
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
	resp := &settings.BuildSettings{}
	var buildSettings appv1.BuildSettings
	err = client.Get(ctx, cr.ObjectKey{
		Namespace: namespace,
		Name:      appName,
	}, &buildSettings)

	if errors.IsNotFound(err) {
		return resp, nil
	}
	if err != nil {
		log.Error().Err(err).Msg("get build settings error")
		return nil, err
	}
	resp.Dockerfile = buildSettings.Spec.Dockerfile

	gitSettings := buildSettings.Spec.Git
	gitRepo, err := b.data.GitRepo.Query().
		Where(gitrepo.Name(gitSettings.RepoName)).
		Only(ctx)
	if generated.IsNotFound(err) {
		return resp, nil
	}
	if err != nil {
		log.Error().Err(err).Msg("get git repo error")
		return nil, err
	}
	resp.Git = &settings.BuildSettings_GitSettings{
		RepoName: gitSettings.RepoName,
		RepoPath: gitSettings.RepoPath,
		Url:      fmt.Sprintf("%s/%s", gitRepo.URL, gitSettings.RepoPath),
	}

	imageSettings := buildSettings.Spec.Image
	imageRepo, err := b.data.ImageRepo.Query().
		Where(imagerepo.Name(imageSettings.RepoName)).
		Only(ctx)
	if generated.IsNotFound(err) {
		return resp, nil
	}
	if err != nil {
		log.Error().Err(err).Msg("get image repo error")
		return nil, err
	}
	resp.Image = &settings.BuildSettings_ImageSettings{
		RepoName: imageSettings.RepoName,
		Url:      imageRepo.URL,
	}

	return resp, nil
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
					RepoName: gitSettings.GetRepoName(),
					RepoPath: gitSettings.GetRepoPath(),
				},
				Image: appv1.ImageSettings{
					RepoName: imageSettings.GetRepoName(),
				},
				Dockerfile: dockerfile,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("create build settings error")
			return err
		}
	} else {
		kubeBuildSettings.Spec.Git.RepoName = gitSettings.GetRepoName()
		kubeBuildSettings.Spec.Git.RepoPath = gitSettings.GetRepoPath()
		kubeBuildSettings.Spec.Image.RepoName = imageSettings.GetRepoName()
		kubeBuildSettings.Spec.Dockerfile = dockerfile
		err = client.Update(ctx, &kubeBuildSettings)
		if err != nil {
			log.Error().Err(err).Msg("update build settings error")
			return err
		}
	}

	return nil

}
