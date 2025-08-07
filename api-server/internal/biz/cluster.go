package biz

import (
	"api-server/api/v1/cluster"
	"api-server/internal/data"
	"context"
)

type ClusterBiz struct {
	db *data.Data
}

func NewClusterBiz(
	db *data.Data,
) *ClusterBiz {
	return &ClusterBiz{
		db: db,
	}
}

func (b *ClusterBiz) ListClusters(ctx context.Context, request *cluster.ListClustersRequest) (*cluster.ListClustersResponse, error) {
	return nil, nil
}
