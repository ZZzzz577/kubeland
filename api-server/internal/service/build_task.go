package service

import (
	"api-server/api/v1/application"
	"api-server/api/v1/build_task"
	"api-server/internal/biz"
	"context"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BuildTaskService struct {
	task.UnimplementedBuildTaskServiceServer

	biz *biz.BuildTaskBiz
}

func NewBuildTaskService(
	biz *biz.BuildTaskBiz,
) *BuildTaskService {
	return &BuildTaskService{
		biz: biz,
	}
}

func (b *BuildTaskService) Register(gs *kgrpc.Server, hs *http.Server) {
	task.RegisterBuildTaskServiceServer(gs, b)
	task.RegisterBuildTaskServiceHTTPServer(hs, b)
}

func (b *BuildTaskService) Create(ctx context.Context, request *application.IdentityRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, b.biz.Create(ctx, request)
}

func (b *BuildTaskService) List(ctx context.Context, request *application.IdentityRequest) (*task.ListBuildTaskResponse, error) {
	return b.biz.List(ctx, request)
}
