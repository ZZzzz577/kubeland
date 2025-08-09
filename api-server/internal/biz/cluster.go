package biz

import (
	"api-server/api/v1/cluster"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	clusterdb "api-server/internal/data/generated/cluster"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ClusterNotFoundError         = status.Error(codes.NotFound, "cluster not found")
	ClusterSecurityNotFoundError = status.Error(codes.NotFound, "cluster security not found")
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
			return c.toProto(item)
		}),
	}, err
}

func (c *ClusterBiz) GetCluster(ctx context.Context, request *cluster.IdRequest) (*cluster.Cluster, error) {
	cr, err := c.db.Cluster.Query().
		Where(clusterdb.ID(request.Id)).
		WithSecurity().
		Only(ctx)
	if generated.IsNotFound(err) {
		return nil, ClusterNotFoundError
	}
	if err != nil {
		log.Error().Err(err).Msg("failed to get cluster")
		return nil, err
	}
	return c.toProto(cr), nil
}

func (c *ClusterBiz) CreateCluster(ctx context.Context, request *cluster.Cluster) error {
	return c.db.WithTx(ctx, func(tx *generated.Tx) error {
		cr, err := tx.Cluster.Create().
			SetName(request.Name).
			SetDescription(request.Description).
			SetAddress(request.Address).
			Save(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to create cluster")
			return err
		}
		security := request.Security
		css := tx.ClusterSecurity.Create().
			SetCluster(cr).
			SetType(uint8(security.Type))
		if security.Type == cluster.Cluster_Security_TLS_CERT {
			css.SetCert(security.Cert).
				SetKey(security.Key)
		} else if security.Type == cluster.Cluster_Security_TLS_TOKEN {
			css.SetToken(security.Token)
		}
		err = css.Exec(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to create cluster security")
			return err
		}
		return nil
	})
}

func (c *ClusterBiz) UpdateCluster(ctx context.Context, request *cluster.Cluster) error {
	return c.db.WithTx(ctx, func(tx *generated.Tx) error {
		cr, err := tx.Cluster.Query().
			WithSecurity().
			Where().
			Only(ctx)
		if generated.IsNotFound(err) {
			return ClusterNotFoundError
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to query cluster")
			return err
		}

		err = tx.Cluster.UpdateOne(cr).
			SetName(request.Name).
			SetDescription(request.Description).
			SetAddress(request.Address).
			Exec(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to update cluster")
			return err
		}
		if request.Security == nil {
			return nil
		}
		// 更新安全配置
		security := request.Security
		csr := cr.Edges.Security
		if csr == nil {
			return ClusterSecurityNotFoundError
		}
		css := tx.ClusterSecurity.UpdateOne(csr).
			SetType(uint8(security.Type))
		if security.Type == cluster.Cluster_Security_TLS_CERT {
			css.SetCert(security.Cert).
				SetKey(security.Key)
		} else if security.Type == cluster.Cluster_Security_TLS_TOKEN {
			css.SetToken(security.Token)
		}
		return css.Exec(ctx)
	})

}

func (c *ClusterBiz) DeleteCluster(ctx context.Context, request *cluster.IdRequest) error {
	return c.db.WithTx(ctx, func(tx *generated.Tx) error {
		cr, err := tx.Cluster.Query().
			WithSecurity().
			Where().
			Only(ctx)
		if generated.IsNotFound(err) {
			return ClusterNotFoundError
		}
		if err != nil {
			log.Error().Err(err).Msg("failed to query cluster")
			return err
		}
		if err = tx.Cluster.DeleteOne(cr).Exec(ctx); err != nil {
			log.Error().Err(err).Msg("failed to delete cluster")
			return err
		}
		if err = tx.ClusterSecurity.DeleteOne(cr.Edges.Security).Exec(ctx); err != nil {
			log.Error().Err(err).Msg("failed to delete cluster security")
			return err
		}
		return nil
	})
}

func (c *ClusterBiz) toProto(cr *generated.Cluster) *cluster.Cluster {
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
	}
}
