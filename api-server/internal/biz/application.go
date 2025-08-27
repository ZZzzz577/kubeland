package biz

import (
	"api-server/api/v1/application"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ApplicationBiz struct {
	db *data.Data
}

func NewApplicationBiz(
	db *data.Data,
) *ApplicationBiz {
	return &ApplicationBiz{
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
	err := a.db.Application.Create().
		SetName(request.Name).
		SetClusterID(request.ClusterId).
		SetDescription(request.Description).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("create application error")
		return err
	}
	return nil
}

func (a *ApplicationBiz) UpdateApplication(ctx context.Context, request *application.Application) error {
	err := a.db.Application.UpdateOneID(request.Id).
		SetDescription(request.Description).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("update application error")
		return err
	}
	return nil
}

func (a *ApplicationBiz) DeleteApplication(ctx context.Context, request *application.IdRequest) error {
	err := a.db.Application.DeleteOneID(request.Id).Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("delete application error")
		return err
	}
	return nil
}

func (a *ApplicationBiz) toProto(source *generated.Application) *application.Application {
	return &application.Application{
		Id:          source.ID,
		ClusterId:   source.ClusterID,
		Name:        source.Name,
		Description: source.Description,
		CreatedAt:   timestamppb.New(source.CreatedAt),
		UpdateAt:    timestamppb.New(source.UpdatedAt),
	}
}
