package biz

import (
	"api-server/api/v1/application"
	appv1 "api-server/internal/kube/api/v1"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type BuildTaskBiz struct {
	cm *ClusterManagers
}

func NewBuildTaskBiz(
	cm *ClusterManagers,
) *BuildTaskBiz {
	return &BuildTaskBiz{
		cm: cm,
	}
}

func (b *BuildTaskBiz) Create(ctx context.Context, request *application.IdentityRequest) error {
	appName := request.GetName()
	client, err := b.cm.GetClient(ctx, appName)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return err
	}

	namespace := "default"
	var buildSettings appv1.BuildSettings
	err = client.Get(ctx, cr.ObjectKey{
		Namespace: namespace,
		Name:      appName,
	}, &buildSettings)
	if err != nil {
		log.Error().Err(err).Msg("get build settings error")
		return err
	}

	taskName := fmt.Sprintf("%s-%s", appName, time.Now().Format("20060102150405"))
	labels := map[string]string{"app": appName}
	buildTask := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      taskName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: lo.ToPtr(int32(0)),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "build-task",
							Image: "crpi-mgl4ujhwwhrsi5e3.cn-hangzhou.personal.cr.aliyuncs.com/kubeland/build:v1",
							Env: []corev1.EnvVar{
								{
									Name:  "GIT_URL",
									Value: buildSettings.Spec.Git.Url,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "dockerfile-config",
									MountPath: "/app/config",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: []corev1.Volume{
						{
							Name: "dockerfile-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: fmt.Sprintf("%s-dockerfile-cm", buildSettings.Name),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	err = client.Create(ctx, buildTask)
	if err != nil {
		log.Error().Err(err).Msg("create build task error")
		return err
	}

	return nil
}
