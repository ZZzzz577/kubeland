package biz

import (
	"api-server/api/v1/cluster"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"context"

	clusterdb "api-server/internal/data/generated/cluster"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (c *ClusterBiz) ListClusters(ctx context.Context, request *cluster.ListClustersRequest) (*cluster.ListClustersResponse, error) {
	query := c.db.Cluster.Query()
	page, list, err := data.Page[*generated.ClusterQuery, *generated.Cluster](ctx, query, request.Page)
	if err != nil {
		log.Error().Err(err).Msg("ListClusters")
		return nil, err
	}
	return &cluster.ListClustersResponse{
		Pagination: page,
		Items: lo.Map(list, func(item *generated.Cluster, index int) *cluster.Cluster {
			return &cluster.Cluster{
				Id:          item.ID,
				Name:        item.Name,
				Description: item.Description,
				Address:     item.Address,
				CreatedAt:   timestamppb.New(item.CreatedAt),
				UpdatedAt:   timestamppb.New(item.UpdatedAt),
			}
		}),
	}, err
}

func (c *ClusterBiz) GetCluster(ctx context.Context, request *cluster.IdRequest) (*cluster.Cluster, error) {
	cr, err := c.db.Cluster.Query().
		Where(clusterdb.ID(request.Id)).
		WithSecurity().
		Only(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get cluster")
		return nil, err
	}
	var security *cluster.Cluster_Security
	csr := cr.Edges.Security
	if csr != nil {
		security = &cluster.Cluster_Security{
			Type:  cluster.Cluster_Security_Type(csr.Type),
			Ca:    csr.Ca,
			Cert:  csr.Cert,
			Key:   csr.Key,
			Token: csr.Token,
		}
	}
	return &cluster.Cluster{
		Id:          cr.ID,
		Name:        cr.Name,
		Description: cr.Description,
		Address:     cr.Address,
		Security:    security,
		CreatedAt:   timestamppb.New(cr.CreatedAt),
		UpdatedAt:   timestamppb.New(cr.UpdatedAt),
	}, nil
}
