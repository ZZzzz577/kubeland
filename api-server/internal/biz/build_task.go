package biz

import (
	"api-server/api/v1/application"
	task "api-server/api/v1/build_task"
	appv1 "api-server/internal/kube/api/v1"
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	controllerruntime "sigs.k8s.io/controller-runtime"
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

	buildTask := b.createBuildTask(&buildSettings)
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

func (b *BuildTaskBiz) createBuildTask(buildSettings *appv1.BuildSettings) *batchv1.Job {
	namespace := "default"
	taskName := fmt.Sprintf("%s-%s", buildSettings.Name, time.Now().Format("20060102150405"))
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      taskName,
			Namespace: namespace,
			// Labels:    labels,
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: lo.ToPtr(int32(0)),
			Template: corev1.PodTemplateSpec{
				//ObjectMeta: metav1.ObjectMeta{
				//	Labels: labels,
				//},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "build-task",
							Image:           "crpi-mgl4ujhwwhrsi5e3.cn-hangzhou.personal.cr.aliyuncs.com/kubeland/build:v1",
							ImagePullPolicy: corev1.PullAlways,
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
}

func (b *BuildTaskBiz) List(ctx context.Context, request *application.IdentityRequest) (*task.ListBuildTaskResponse, error) {
	appName := request.GetName()
	client, err := b.cm.GetClient(ctx, appName)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return nil, err
	}

	var tasks batchv1.JobList
	err = client.List(ctx, &tasks,
		cr.InNamespace("default"),
		// todo change matching labels to matching fields
		// cr.MatchingLabels(map[string]string{"app": appName}),
	)
	if err != nil {
		log.Error().Err(err).Str("app", appName).Msg("list build tasks error")
		return nil, err
	}

	return &task.ListBuildTaskResponse{
		Items: lo.Map(tasks.Items, func(item batchv1.Job, index int) *task.BuildTask {
			return &task.BuildTask{
				Name:      item.Name,
				Status:    item.Status.String(),
				CreatedAt: timestamppb.New(item.CreationTimestamp.Time),
			}
		}),
	}, nil
}

func (b *BuildTaskBiz) Log(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/event-stream")

	flusher, ok := writer.(http.Flusher)
	if !ok {
		http.Error(writer, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	ctx := request.Context()
	vars := mux.Vars(request)

	appName := vars["appName"]
	client, err := b.cm.GetClient(ctx, appName)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
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
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(podList.Items) != 1 {
		err = errors.New("pods num is valid")
		log.Error().Err(err).Str("labelSelector", labelSelector).Msg("pods num is valid")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
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
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = stream.Close()
	}()

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		event := "data: " + scanner.Text() + "\n\n"
		if _, err = writer.Write([]byte(event)); err != nil {
			log.Error().Err(err).Msg("write to response error")
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		flusher.Flush()
	}
	if err = scanner.Err(); err != nil {
		log.Error().Err(err).Msg("scan pod logs error")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}
