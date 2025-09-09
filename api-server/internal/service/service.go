package service

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/google/wire"
)

type Service interface {
	Register(gs *grpc.Server, hs *http.Server)
}

func NewServices(
	cluster *ClusterService,
	application *ApplicationService,
	buildSettings *BuildSettingsService,
	buildTask *BuildTaskService,
	git *GitService,
) []Service {
	return []Service{
		cluster,
		application,
		buildSettings,
		buildTask,
		git,
	}
}

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewServices,
	NewClusterService,
	NewApplicationService,
	NewBuildSettingsService,
	NewBuildTaskService,
	NewGitService,
)
