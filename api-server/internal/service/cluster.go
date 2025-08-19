package service

import (
	"api-server/api/v1/cluster"
	"api-server/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterService struct {
	cluster.UnimplementedClusterServiceServer
	clusterBiz *biz.ClusterBiz
}

func NewClusterService(
	clusterBiz *biz.ClusterBiz,
) *ClusterService {
	return &ClusterService{
		clusterBiz: clusterBiz,
	}
}

func (c *ClusterService) Register(gs *grpc.Server, hs *http.Server) {
	cluster.RegisterClusterServiceServer(gs, c)
	cluster.RegisterClusterServiceHTTPServer(hs, c)
}

func (c *ClusterService) ListClusters(ctx context.Context, request *cluster.ListClustersRequest) (*cluster.ListClustersResponse, error) {
	return c.clusterBiz.ListClusters(ctx, request)
}

func (c *ClusterService) GetCluster(ctx context.Context, request *cluster.IdRequest) (*cluster.Cluster, error) {
	return c.clusterBiz.GetCluster(ctx, request)
}
func (c *ClusterService) CreateCluster(ctx context.Context, request *cluster.Cluster) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.clusterBiz.CreateCluster(ctx, request)
}
func (c *ClusterService) UpdateCluster(ctx context.Context, request *cluster.Cluster) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.clusterBiz.UpdateCluster(ctx, request)
}
func (c *ClusterService) DeleteCluster(ctx context.Context, request *cluster.IdRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, c.clusterBiz.DeleteCluster(ctx, request)
}

func (c *ClusterService) ResolveKubeConfig(ctx context.Context, request *cluster.ResolveKubeConfigRequest) (*cluster.ResolveKubeConfigResponse, error) {
	return c.clusterBiz.ResolveKubeConfig(ctx, request)
}

func (c *ClusterService) TestConnection(ctx context.Context, request *cluster.Connection) (*cluster.TestConnectionResponse, error) {
	return c.clusterBiz.TestConnection(ctx, request)
}

func (c *ClusterService) TestOperator(ctx context.Context, request *cluster.Connection) (*cluster.TestOperatorResponse, error) {
	return c.clusterBiz.TestOperator(ctx, request)
}
