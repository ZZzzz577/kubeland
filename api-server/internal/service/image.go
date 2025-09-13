package service

import (
	"api-server/api/v1/image"
	"api-server/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ImageService struct {
	image.UnimplementedImageServiceServer

	biz *biz.ImageBiz
}

func NewImageService(biz *biz.ImageBiz) *ImageService {
	return &ImageService{
		biz: biz,
	}
}

func (i *ImageService) Register(gs *grpc.Server, hs *http.Server) {
	image.RegisterImageServiceServer(gs, i)
	image.RegisterImageServiceHTTPServer(hs, i)
}

func (i *ImageService) GetImageRepo(ctx context.Context, request *image.IdentityRequest) (*image.ImageRepo, error) {
	return i.biz.GetImageRepo(ctx, request)
}
func (i *ImageService) ListImageRepos(ctx context.Context, request *image.ListImageReposRequest) (*image.ListImageReposResponse, error) {
	return i.biz.ListImageRepos(ctx, request)
}
func (i *ImageService) CreateImageRepo(ctx context.Context, request *image.ImageRepo) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, i.biz.CreateImageRepo(ctx, request)
}
func (i *ImageService) UpdateImageRepo(ctx context.Context, request *image.ImageRepo) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, i.biz.UpdateImageRepo(ctx, request)
}
func (i *ImageService) DeleteImageRepo(ctx context.Context, request *image.IdentityRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, i.biz.DeleteImageRepo(ctx, request)
}
