package service

import (
	"api-server/api/v1/application"
	settings "api-server/api/v1/build_settings"
	"api-server/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BuildSettingsService struct {
	settings.UnimplementedBuildSettingsServiceServer

	biz *biz.BuildSettingsBiz
}

func NewBuildSettingsService(
	biz *biz.BuildSettingsBiz,
) *BuildSettingsService {
	return &BuildSettingsService{
		biz: biz,
	}
}

func (b *BuildSettingsService) Register(gs *grpc.Server, hs *http.Server) {
	settings.RegisterBuildSettingsServiceServer(gs, b)
	settings.RegisterBuildSettingsServiceHTTPServer(hs, b)
}

func (b *BuildSettingsService) GetBuildSettings(ctx context.Context, request *application.IdentityRequest) (*settings.BuildSettings, error) {
	return b.biz.GetBuildSettings(ctx, request)
}

func (b *BuildSettingsService) ApplyBuildSettings(ctx context.Context, request *settings.ApplyBuildSettingsRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, b.biz.ApplyBuildSettings(ctx, request)
}
