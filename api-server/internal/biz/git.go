package biz

import (
	"api-server/api/v1/application"
	"api-server/api/v1/common"
	"api-server/api/v1/git"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"api-server/internal/data/generated/gitrepo"
	appv1 "api-server/internal/kube/api/v1"
	"context"
	"fmt"
	"github.com/google/go-github/v74/github"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/url"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

type GitBiz struct {
	cm     *ClusterManagers
	client *github.Client
	data   *data.Data
}

func NewGitBiz(
	cm *ClusterManagers,
	data *data.Data,
) *GitBiz {
	client := github.NewClient(nil)
	return &GitBiz{
		cm:     cm,
		client: client,
		data:   data,
	}
}

func (g *GitBiz) GetGitRepo(ctx context.Context, request *git.IdentityRequest) (*git.GitRepo, error) {
	gitRepo, err := g.data.GitRepo.Query().
		Where(gitrepo.Name(request.GetName())).
		Only(ctx)
	if generated.IsNotFound(err) {
		return nil, status.Error(codes.NotFound, "git repo not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("get git repo error")
		return nil, err
	}
	return g.toProto(gitRepo, true), nil
}

func (g *GitBiz) ListGitRepos(ctx context.Context, request *git.ListGitReposRequest) (*git.ListGitReposResponse, error) {
	page, repos, err := data.Page[*generated.GitRepoQuery](ctx, g.data.GitRepo.Query(), request.Page)
	if err != nil {
		log.Error().Err(err).Msg("list git repos error")
		return nil, err
	}
	return &git.ListGitReposResponse{
		Pagination: page,
		Items: lo.Map(repos, func(item *generated.GitRepo, index int) *git.GitRepo {
			return g.toProto(item, false)
		}),
	}, nil
}

func (g *GitBiz) CreateGitRepo(ctx context.Context, request *git.GitRepo) error {
	return g.data.WithTx(ctx, func(tx *generated.Tx) error {
		exist, err := tx.GitRepo.Query().
			Where(gitrepo.Name(request.GetName())).
			Exist(ctx)
		if err != nil {
			log.Error().Err(err).Msg("check git repo exist error")
			return err
		}
		if exist {
			return status.Error(codes.AlreadyExists, "git repo already exists")
		}
		err = tx.GitRepo.Create().
			SetName(request.GetName()).
			SetDescription(request.GetDescription()).
			SetURL(request.GetUrl()).
			SetToken(request.GetToken()).
			Exec(ctx)
		if err != nil {
			log.Error().Err(err).Msg("create git repo error")
			return err
		}
		return nil
	})
}

func (g *GitBiz) UpdateGitRepo(ctx context.Context, request *git.GitRepo) error {
	err := g.data.GitRepo.Update().
		SetDescription(request.GetDescription()).
		SetURL(request.GetUrl()).
		SetToken(request.GetToken()).
		Where(gitrepo.Name(request.GetName())).
		Exec(ctx)
	if generated.IsNotFound(err) {
		return status.Error(codes.NotFound, "git repo not found")
	}
	if err != nil {
		log.Error().Err(err).Msg("update git repo error")
		return err
	}
	return nil
}

func (g *GitBiz) DeleteGitRepo(ctx context.Context, request *git.IdentityRequest) error {
	_, err := g.data.GitRepo.Delete().
		Where(gitrepo.Name(request.GetName())).
		Exec(ctx)
	if err != nil {
		log.Error().Err(err).Msg("delete git repo error")
		return err
	}
	return nil
}

func (g *GitBiz) toProto(source *generated.GitRepo, showSecret bool) *git.GitRepo {
	target := &git.GitRepo{
		Name:        source.Name,
		Description: source.Description,
		Url:         source.URL,
		CreatedAt:   timestamppb.New(source.CreatedAt),
		UpdatedAt:   timestamppb.New(source.UpdatedAt),
	}
	if showSecret {
		target.Token = source.Token
	}
	return target
}

func (g *GitBiz) GetGitSettings(ctx context.Context, request *application.IdentityRequest) (*git.GitSettings, error) {
	appName := request.GetName()

	client, err := g.cm.GetClient(ctx, appName)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return nil, err
	}

	namespace := "default"
	var buildSettings appv1.BuildSettings
	err = client.Get(ctx, cr.ObjectKey{
		Namespace: namespace,
		Name:      appName,
	}, &buildSettings)

	if errors.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		log.Error().Err(err).Msg("get build settings error")
		return nil, err
	}

	gitSettings := buildSettings.Spec.Git
	gitRepoName := gitSettings.RepoName
	gitRepo, err := g.data.GitRepo.Query().
		Where(gitrepo.Name(gitRepoName)).
		Only(ctx)
	if generated.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		log.Error().Err(err).Msg("get git repo error")
		return nil, err
	}

	return &git.GitSettings{
		Url:   fmt.Sprintf("%s/%s", gitRepo.URL, gitSettings.RepoPath),
		Token: gitRepo.Token,
	}, nil

}

func (g *GitBiz) ListBranches(ctx context.Context, request *application.IdentityRequest) (*git.ListBranchesResponse, error) {
	gitSettings, err := g.GetGitSettings(ctx, request)
	if err != nil {
		log.Error().Err(err).Msg("get git settings error")
		return nil, err
	}

	owner, repo, err := g.getDetailFromURL(gitSettings.GetUrl())
	if err != nil {
		log.Error().Err(err).Msg("get detail from url error")
		return nil, err
	}

	client := g.client.WithAuthToken(gitSettings.GetToken())
	branches := make([]string, 0)
	nextPage := 1
	for nextPage != 0 {
		curBranches, response, err := client.Repositories.ListBranches(ctx, owner, repo, &github.BranchListOptions{
			ListOptions: github.ListOptions{
				Page:    nextPage,
				PerPage: 100,
			},
		})
		if err != nil {
			log.Error().Err(err).Msg("list branches error")
			return nil, err
		}
		branches = append(branches, lo.Map(curBranches, func(item *github.Branch, _ int) string {
			return item.GetName()
		})...)
		nextPage = response.NextPage
	}
	return &git.ListBranchesResponse{Items: branches}, nil
}

func (g *GitBiz) ListCommits(ctx context.Context, request *git.ListCommitsRequest) (*git.ListCommitsResponse, error) {
	page, size := 1, 20
	if request.GetPage().GetCurrent() > 0 {
		page = int(request.GetPage().GetCurrent())
	}
	if request.GetPage().GetSize() > 0 {
		size = int(request.GetPage().GetSize())
	}

	gitSettings, err := g.GetGitSettings(ctx, &application.IdentityRequest{Name: request.GetName()})
	if err != nil {
		log.Error().Err(err).Msg("get git settings error")
		return nil, err
	}

	owner, repo, err := g.getDetailFromURL(gitSettings.GetUrl())
	if err != nil {
		log.Error().Err(err).Msg("get detail from url error")
		return nil, err
	}

	client := g.client.WithAuthToken(gitSettings.GetToken())
	commits, response, err := client.Repositories.ListCommits(ctx, owner, repo, &github.CommitsListOptions{
		SHA: request.BranchName,
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: size,
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("list commits error")
		return nil, err

	}
	return &git.ListCommitsResponse{
		Pagination: &common.Pagination{
			Current:   uint32(page),
			Size:      uint32(size),
			TotalPage: uint32(response.LastPage),
		},
		Items: lo.Map(commits, func(item *github.RepositoryCommit, _ int) *git.ListCommitsResponse_Commit {
			createAt := item.GetCommit().GetCommitter().GetDate().Time
			return &git.ListCommitsResponse_Commit{
				Sha:       item.GetSHA(),
				Message:   item.GetCommit().GetMessage(),
				CreatedAt: timestamppb.New(createAt),
			}
		}),
	}, nil
}

func (g *GitBiz) getDetailFromURL(gitUrl string) (string, string, error) {
	gitUrl = strings.TrimSuffix(gitUrl, ".git")
	var path string
	if strings.HasPrefix(gitUrl, "http://") || strings.HasPrefix(gitUrl, "https://") {
		parse, err := url.Parse(gitUrl)
		if err != nil {
			log.Error().Err(err).Msg("parse git url error")
			return "", "", err
		}
		path = parse.Path
	} else if strings.HasPrefix(gitUrl, "git@") {
		split := strings.Split(gitUrl, ":")
		if len(split) >= 2 {
			path = split[1]
		}
	}
	path = strings.TrimPrefix(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		err := fmt.Errorf("invalid git url")
		log.Error().Err(err).Msg("invalid git url")
		return "", "", err
	}
	return parts[0], parts[1], nil
}
