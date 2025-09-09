package biz

import (
	"api-server/api/v1/application"
	"api-server/api/v1/common"
	"api-server/api/v1/git"
	appv1 "api-server/internal/kube/api/v1"
	"context"
	"fmt"
	"github.com/google/go-github/v74/github"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/url"
	controllerruntime "sigs.k8s.io/controller-runtime"
	cr "sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	GitUrlKey   = "GIT_URL"
	GitTokenKey = "GIT_TOKEN"
)

type GitBiz struct {
	cm     *ClusterManagers
	client *github.Client
}

func NewGitBiz(
	cm *ClusterManagers,
) *GitBiz {
	client := github.NewClient(nil)
	return &GitBiz{
		cm:     cm,
		client: client,
	}
}

func (g *GitBiz) ApplyGitSettings(ctx context.Context, request *git.ApplyGitSettingsRequest) error {
	appName := request.GetName()
	namespace := "default"

	client, err := g.cm.GetClient(ctx, appName)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return err
	}

	var buildSettings appv1.BuildSettings
	err = client.Get(ctx, cr.ObjectKey{
		Namespace: namespace,
		Name:      appName,
	}, &buildSettings)

	gitSettingsName := fmt.Sprintf("%s-git", appName)
	gitSettings, err := client.CoreV1().Secrets(namespace).Get(ctx, gitSettingsName, metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		log.Error().Err(err).Msg("get git settings error")
		return err
	}

	if errors.IsNotFound(err) {
		gitSettings, err = client.CoreV1().Secrets(namespace).Create(ctx, &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      gitSettingsName,
				Namespace: namespace,
			},
			StringData: map[string]string{
				GitUrlKey:   request.GetGitSettings().GetUrl(),
				GitTokenKey: request.GetGitSettings().GetToken(),
			},
		}, metav1.CreateOptions{})
		if err != nil {
			log.Error().Err(err).Msg("create git secret error")
			return err
		}
	} else {
		gitSettings.Data[GitUrlKey] = []byte(request.GetGitSettings().GetUrl())
		gitSettings.Data[GitTokenKey] = []byte(request.GetGitSettings().GetToken())
		gitSettings, err = client.CoreV1().Secrets(namespace).Update(ctx, gitSettings, metav1.UpdateOptions{})
		if err != nil {
			log.Error().Err(err).Msg("update git secret error")
			return err
		}
	}

	if err = controllerruntime.SetControllerReference(&buildSettings, gitSettings, client.Scheme); err != nil {
		log.Error().Err(err).Msg("set controller reference error")
		return err
	}

	return nil

}

func (g *GitBiz) GetGitSettings(ctx context.Context, request *application.IdentityRequest) (*git.GitSettings, error) {
	appName := request.GetName()
	namespace := "default"

	client, err := g.cm.GetClient(ctx, appName)
	if err != nil {
		log.Error().Err(err).Msg("get cluster client error")
		return nil, err
	}

	gitSettings, err := client.CoreV1().Secrets(namespace).
		Get(ctx, fmt.Sprintf("%s-git", appName), metav1.GetOptions{})
	if err != nil {
		log.Error().Err(err).Msg("get git settings error")
		return nil, err
	}
	return &git.GitSettings{
		Url:   string(gitSettings.Data[GitUrlKey]),
		Token: string(gitSettings.Data[GitTokenKey]),
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
