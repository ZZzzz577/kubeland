package biz

import (
	"api-server/api/v1/application"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	appv1 "api-server/internal/kube/api/v1"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
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

func (a *ApplicationBiz) ListApplications(ctx context.Context, request *application.ListApplicationsRequest) (*application.ListApplicationsResponse, error) {
	page, applications, err := data.Page[*generated.ApplicationQuery](ctx, a.db.Application.Query(), request.Page)
	if err != nil {
		log.Error().Err(err).Msg("list application error")
		return nil, err
	}
	return &application.ListApplicationsResponse{
		Pagination: page,
		Items: lo.Map(applications, func(item *generated.Application, index int) *application.Application {
			return a.toProto(item)
		}),
	}, nil
}

func (a *ApplicationBiz) GetApplication(ctx context.Context, request *application.IdRequest) (*application.Application, error) {
	app, err := a.db.Application.Get(ctx, request.Id)
	if err != nil {
		log.Error().Err(err).Msg("get application error")
		return nil, err
	}
	return a.toProto(app), nil
}

func (a *ApplicationBiz) CreateApplication(ctx context.Context, request *application.Application) error {
	client, err := a.cm.GetClient(request.ClusterId)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return err
	}
	return client.Create(ctx, &appv1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      request.Name,
		},
		Spec: appv1.ApplicationSpec{
			Description: request.Description,
		},
	})
}

func (a *ApplicationBiz) UpdateApplication(ctx context.Context, request *application.Application) error {
	source, err := a.db.Application.Get(ctx, request.Id)
	if err != nil {
		log.Error().Err(err).Msg("get application error")
		return err
	}

	client, err := a.cm.GetClient(source.ClusterID)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return err
	}

	var app appv1.Application
	err = client.Get(ctx, cr.ObjectKey{
		Name:      source.Name,
		Namespace: "default",
	}, &app)
	if err != nil {
		log.Error().Err(err).Msg("get application error")
		return err
	}
	app.Spec.Description = request.Description
	return client.Update(ctx, &app)
}

func (a *ApplicationBiz) DeleteApplication(ctx context.Context, request *application.IdRequest) error {
	app, err := a.db.Application.Get(ctx, request.Id)
	if err != nil {
		log.Error().Err(err).Msg("get application error")
		return err
	}

	client, err := a.cm.GetClient(app.ClusterID)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return err
	}
	return client.Delete(ctx, &appv1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      app.Name,
		},
	})
}

func (a *ApplicationBiz) toProto(source *generated.Application) *application.Application {
	return &application.Application{
		Id:          source.ID,
		ClusterId:   source.ClusterID,
		Name:        source.Name,
		Description: source.Description,
		CreatedAt:   timestamppb.New(source.CreatedAt),
	}
}
