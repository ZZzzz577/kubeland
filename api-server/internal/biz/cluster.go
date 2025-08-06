package biz

import (
	"api-server/api/v1/cluster"
	"context"
)

type ClusterBiz struct {}

func NewClusterBiz() *ClusterBiz {
	return &ClusterBiz{}
}

func (b *ClusterBiz) ListClusters(ctx context.Context, request *cluster.ListClustersRequest) (*cluster.ListClustersResponse, error) {
	return nil, nil
}