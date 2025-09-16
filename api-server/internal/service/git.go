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

func (g *GitService) GetGitRepo(ctx context.Context, request *git.IdentityRequest) (*git.GitRepo, error) {
	return g.biz.GetGitRepo(ctx, request)
}
func (g *GitService) ListGitRepos(ctx context.Context, request *git.ListGitReposRequest) (*git.ListGitReposResponse, error) {
	return g.biz.ListGitRepos(ctx, request)
}
func (g *GitService) CreateGitRepo(ctx context.Context, request *git.GitRepo) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, g.biz.CreateGitRepo(ctx, request)
}
func (g *GitService) UpdateGitRepo(ctx context.Context, request *git.GitRepo) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, g.biz.UpdateGitRepo(ctx, request)
}

func (g *GitService) DeleteGitRepo(ctx context.Context, request *git.IdentityRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, g.biz.DeleteGitRepo(ctx, request)
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
