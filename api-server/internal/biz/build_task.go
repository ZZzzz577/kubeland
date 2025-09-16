package biz

import (
	"api-server/api/v1/application"
	task "api-server/api/v1/build_task"
	"api-server/internal/data"
	"api-server/internal/data/generated/gitrepo"
	appv1 "api-server/internal/kube/api/v1"
	"api-server/internal/kube/commons/job"
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	controllerruntime "sigs.k8s.io/controller-runtime"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
	"slices"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type BuildTaskBiz struct {
	cm   *ClusterManagers
	data *data.Data
}

func NewBuildTaskBiz(
	cm *ClusterManagers,
	data *data.Data,
) *BuildTaskBiz {
	return &BuildTaskBiz{
		cm:   cm,
		data: data,
	}
}

func (b *BuildTaskBiz) Create(ctx context.Context, request *task.CreateBuildTaskRequest) error {
	appName := request.GetAppName()
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

	buildTask := b.createBuildTask(ctx, request, &buildSettings)
	if err = controllerruntime.SetControllerReference(&buildSettings, buildTask, client.Scheme); err != nil {
		log.Error().Err(err).Msg("set controller reference error")
		return err
	}

	err = client.Create(ctx, buildTask)
	if err != nil {
		log.Error().Err(err).Msg("create build task error")
		return err
	}

	return nil
}

func (b *BuildTaskBiz) createBuildTask(ctx context.Context, request *task.CreateBuildTaskRequest, buildSettings *appv1.BuildSettings) *batchv1.Job {
	namespace := "default"
	taskName := fmt.Sprintf("%s-%s", buildSettings.Name, time.Now().Format("20060102150405"))
	gitSettings := buildSettings.Spec.Git
	gitRepo, err := b.data.GitRepo.Query().
		Where(gitrepo.Name(gitSettings.RepoName)).
		Only(ctx)
	if err != nil {
		log.Error().Err(err).Msg("get git repo error")
		return nil
	}
	gitUrl := fmt.Sprintf("%s/%s", gitRepo.URL, gitSettings.RepoPath)

	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      taskName,
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: lo.ToPtr(int32(0)),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "build-task",
							Image:           "crpi-mgl4ujhwwhrsi5e3.cn-hangzhou.personal.cr.aliyuncs.com/kubeland/build:v1",
							ImagePullPolicy: corev1.PullAlways,
							Env: []corev1.EnvVar{
								{
									Name:  "GIT_URL",
									Value: gitUrl,
								},
								{
									Name:  "GIT_COMMIT",
									Value: request.GetCommit(),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "git",
									MountPath: "/app/config/git",
								},
								{
									Name:      "image",
									MountPath: "/app/config/image",
								},
								{
									Name:      "dockerfile",
									MountPath: "/app/config/dockerfile",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: []corev1.Volume{
						{
							Name: "git",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: fmt.Sprintf("%s-git", buildSettings.Name),
								},
							},
						},
						{
							Name: "image",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: fmt.Sprintf("%s-image", buildSettings.Name),
								},
							},
						},
						{
							Name: "dockerfile",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: fmt.Sprintf("%s-dockerfile", buildSettings.Name),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (b *BuildTaskBiz) List(ctx context.Context, request *application.IdentityRequest) (*task.ListBuildTaskResponse, error) {
	appName := request.GetName()
	client, err := b.cm.GetClient(ctx, appName)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return nil, err
	}

	tasks, err := client.BatchV1().Jobs("default").
		List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Error().Err(err).Str("app", appName).Msg("list build tasks error")
		return nil, err
	}

	slices.SortFunc(tasks.Items, func(a, b batchv1.Job) int {
		if a.CreationTimestamp.Before(&b.CreationTimestamp) {
			return 1
		} else {
			return -1
		}
	})

	return &task.ListBuildTaskResponse{
		Items: lo.Map(tasks.Items, func(item batchv1.Job, index int) *task.BuildTask {
			return &task.BuildTask{
				Name:      item.Name,
				Status:    task.Status(job.GetJobStatus(item)),
				CreatedAt: timestamppb.New(item.CreationTimestamp.Time),
			}
		}),
	}, nil
}

func (b *BuildTaskBiz) Log(writer http.ResponseWriter, request *http.Request) {
	conn, err := upGrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Error().Err(err).Msg("upgrade websocket error")
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	ctx := request.Context()
	vars := mux.Vars(request)

	appName := vars["appName"]
	client, err := b.cm.GetClient(ctx, appName)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return
	}

	jobName := vars["jobName"]
	namespace := "default"

	labelSelector := fmt.Sprintf("job-name=%s", jobName)
	podList, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		log.Error().Err(err).Str("labelSelector", labelSelector).Msg("list pods error")
		return
	}

	if len(podList.Items) != 1 {
		err = errors.New("pods num is valid")
		log.Error().Err(err).Str("labelSelector", labelSelector).Msg("pods num is valid")
		return
	}

	pod := podList.Items[0]
	req := client.CoreV1().Pods(pod.GetNamespace()).GetLogs(pod.GetName(),
		&corev1.PodLogOptions{
			Follow:     true,
			Timestamps: true,
		})
	stream, err := req.Stream(ctx)
	if err != nil {
		log.Error().Err(err).Msg("get pod logs error")
		return
	}
	defer func() {
		_ = stream.Close()
	}()

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		if err = conn.WriteMessage(websocket.TextMessage, []byte(scanner.Text()+"\n")); err != nil {
			log.Error().Err(err).Msg("write to response error")
			return
		}
	}
	if err = scanner.Err(); err != nil {
		log.Error().Err(err).Msg("scan pod logs error")
		return
	}

}
