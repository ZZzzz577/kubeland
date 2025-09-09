package service

import (
	"api-server/api/v1/application"
	"api-server/api/v1/git"
	"api-server/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GitService struct {
	git.UnimplementedGitServiceServer

	biz *biz.GitBiz
}

func NewGitService(biz *biz.GitBiz) *GitService {
	return &GitService{
		biz: biz,
	}
}

func (g *GitService) Register(gs *grpc.Server, hs *http.Server) {
	git.RegisterGitServiceServer(gs, g)
	git.RegisterGitServiceHTTPServer(hs, g)
}

func (g *GitService) ApplyGitSettings(ctx context.Context, request *git.ApplyGitSettingsRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, g.biz.ApplyGitSettings(ctx, request)
}
func (g *GitService) GetGitSettings(ctx context.Context, request *application.IdentityRequest) (*git.GitSettings, error) {
	return g.biz.GetGitSettings(ctx, request)
}
func (g *GitService) ListBranches(ctx context.Context, request *application.IdentityRequest) (*git.ListBranchesResponse, error) {
	return g.biz.ListBranches(ctx, request)
}

func (g *GitService) ListCommits(ctx context.Context, request *git.ListCommitsRequest) (*git.ListCommitsResponse, error) {
	return g.biz.ListCommits(ctx, request)
}
