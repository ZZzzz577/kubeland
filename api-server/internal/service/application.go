package service

import (
	"api-server/api/v1/application"
	"api-server/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ApplicationService struct {
	application.UnimplementedApplicationServiceServer

	biz *biz.ApplicationBiz
}

func NewApplicationService(
	biz *biz.ApplicationBiz,
) *ApplicationService {
	return &ApplicationService{
		biz: biz,
	}
}

func (a *ApplicationService) Register(gs *grpc.Server, hs *http.Server) {
	application.RegisterApplicationServiceServer(gs, a)
	application.RegisterApplicationServiceHTTPServer(hs, a)
}

func (a *ApplicationService) CreateApplication(ctx context.Context, request *application.Application) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, a.biz.CreateApplication(ctx, request)
}
