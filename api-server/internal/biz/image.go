package biz

import (
	"api-server/api/v1/image"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"api-server/internal/data/generated/imagerepo"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ImageBiz struct {
	data *data.Data
}

func NewImageBiz(
	data *data.Data,
) *ImageBiz {
	return &ImageBiz{
		data: data,
	}
}

func (i *ImageBiz) GetImageRepo(ctx context.Context, request *image.IdentityRequest) (*image.ImageRepo, error) {
	imageRepo, err := i.data.ImageRepo.Query().
		Where(imagerepo.Name(request.GetName())).
		Only(ctx)
	if generated.IsNotFound(err) {
		return nil, status.Error(codes.NotFound, "image repo not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("get image repo error")
		return nil, err
	}
	return i.toProto(imageRepo, true), nil
}

func (i *ImageBiz) ListImageRepos(ctx context.Context, request *image.ListImageReposRequest) (*image.ListImageReposResponse, error) {
	page, repos, err := data.Page[*generated.ImageRepoQuery](ctx, i.data.ImageRepo.Query(), request.Page)
	if err != nil {
		log.Error().Err(err).Msg("list image repos error")
		return nil, err
	}
	return &image.ListImageReposResponse{
		Pagination: page,
		Items: lo.Map(repos, func(item *generated.ImageRepo, index int) *image.ImageRepo {
			return i.toProto(item, false)
		}),
	}, nil
}

func (i *ImageBiz) CreateImageRepo(ctx context.Context, request *image.ImageRepo) error {
	return i.data.WithTx(ctx, func(tx *generated.Tx) error {
		exist, err := tx.ImageRepo.Query().
			Where(imagerepo.Name(request.GetName())).
			Exist(ctx)
		if err != nil {
			log.Error().Err(err).Msg("check image repo exist error")
			return err
		}
		if exist {
			return status.Error(codes.AlreadyExists, "image repo already exists")
		}
		err = tx.ImageRepo.Create().
			SetName(request.GetName()).
			SetDescription(request.GetDescription()).
			SetURL(request.GetUrl()).
			SetUsername(request.GetUsername()).
			SetPassword(request.GetPassword()).
			Exec(ctx)
		if err != nil {
			log.Error().Err(err).Msg("create image repo error")
			return err
		}
		return nil
	})
}

func (i *ImageBiz) UpdateImageRepo(ctx context.Context, request *image.ImageRepo) error {
	err := i.data.ImageRepo.Update().
		SetDescription(request.GetDescription()).
		SetURL(request.GetUrl()).
		SetUsername(request.GetUsername()).
		SetPassword(request.GetPassword()).
		Where(imagerepo.Name(request.GetName())).
		Exec(ctx)
	if generated.IsNotFound(err) {
		return status.Error(codes.NotFound, "image repo not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("update image repo error")
		return err
	}
	return nil
}

func (i *ImageBiz) DeleteImageRepo(ctx context.Context, request *image.IdentityRequest) error {
	_, err := i.data.ImageRepo.Delete().
		Where(imagerepo.Name(request.GetName())).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("delete image repo error")
		return err
	}
	return nil
}

func (i *ImageBiz) toProto(source *generated.ImageRepo, showSecret bool) *image.ImageRepo {
	target := &image.ImageRepo{
		Name:        source.Name,
		Description: source.Description,
		Url:         source.URL,
		Username:    source.Username,
		CreatedAt:   timestamppb.New(source.CreatedAt),
		UpdatedAt:   timestamppb.New(source.UpdatedAt),
	}
	if showSecret {
		target.Password = source.Password
	}
	return target
}
