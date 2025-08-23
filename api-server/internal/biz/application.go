package biz

import (
	"api-server/api/v1/application"
	"api-server/internal/data"
	applicationv1 "api-server/internal/kube/apis/application/v1"
	"context"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApplicationBiz struct {
	cm *ClusterManagers
	db *data.Data
}

func NewApplicationBiz(
	cm *ClusterManagers,
	db *data.Data,
) *ApplicationBiz {
	return &ApplicationBiz{
		cm: cm,
		db: db,
	}
}

func (a *ApplicationBiz) CreateApplication(ctx context.Context, request *application.Application) error {
	client, err := a.cm.GetClient(request.ClusterId)
	if err != nil {
		log.Error().Err(err).Msg("get cluster error")
		return err
	}
	_, err = client.KubelandV1().Applications("default").Create(ctx, &applicationv1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name: request.Name,
		},
		Spec: applicationv1.ApplicationSpec{
			Description: request.Description,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Error().Err(err).Msg("create application error")
		return err
	}
	return nil
}
