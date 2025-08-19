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
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/url"
)

var (
	ClusterNotFoundError           = status.Error(codes.NotFound, "cluster not found")
	ClusterConnectionRequiredError = status.Error(codes.InvalidArgument, "cluster connection is required")
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
	query := c.db.Cluster.Query().WithConnection()
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
		WithConnection().
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
	if request.Connection == nil {
		return ClusterConnectionRequiredError
	}
	return c.db.WithTx(ctx, func(tx *generated.Tx) error {
		cr, err := tx.Cluster.Create().
			SetName(request.Name).
			SetDescription(request.Description).
			Save(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to create cluster")
			return err
		}
		connection := request.Connection
		css := tx.ClusterConnection.Create().
			SetCluster(cr).
			SetAddress(connection.Address).
			SetCa(connection.Ca).
			SetType(uint8(connection.Type))
		if connection.Type == cluster.Connection_TLS_CERT {
			css.SetCert(connection.Cert).
				SetKey(connection.Key)
		} else if connection.Type == cluster.Connection_TLS_TOKEN {
			css.SetToken(connection.Token)
		}
		err = css.Exec(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to create cluster connection")
			return err
		}
		return nil
	})
}

func (c *ClusterBiz) UpdateCluster(ctx context.Context, request *cluster.Cluster) error {
	return c.db.WithTx(ctx, func(tx *generated.Tx) error {
		cr, err := tx.Cluster.Query().
			WithConnection().
			Where(clusterdb.ID(request.Id)).
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
			Exec(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to update cluster")
			return err
		}
		if request.Connection == nil {
			return nil
		}
		// 更新安全配置
		connection := request.Connection
		ccr := cr.Edges.Connection
		ccs := tx.ClusterConnection.UpdateOne(ccr).
			SetAddress(connection.Address).
			SetCa(connection.Ca).
			SetType(uint8(connection.Type))
		if connection.Type == cluster.Connection_TLS_CERT {
			ccs.SetCert(connection.Cert).
				SetKey(connection.Key)
		} else if connection.Type == cluster.Connection_TLS_TOKEN {
			ccs.SetToken(connection.Token)
		}
		return ccs.Exec(ctx)
	})

}

func (c *ClusterBiz) DeleteCluster(ctx context.Context, request *cluster.IdRequest) error {
	return c.db.WithTx(ctx, func(tx *generated.Tx) error {
		cr, err := tx.Cluster.Query().
			WithConnection().
			Where(clusterdb.ID(request.Id)).
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
		if err = tx.ClusterConnection.DeleteOne(cr.Edges.Connection).Exec(ctx); err != nil {
			log.Error().Err(err).Msg("failed to delete cluster connection")
			return err
		}
		return nil
	})
}

func (c *ClusterBiz) toProto(cr *generated.Cluster) *cluster.Cluster {
	var connection *cluster.Connection
	ccr := cr.Edges.Connection
	if ccr != nil {
		connection = &cluster.Connection{
			Address: ccr.Address,
			Type:    cluster.Connection_Type(ccr.Type),
			Ca:      ccr.Ca,
			Cert:    ccr.Cert,
			Key:     ccr.Key,
			Token:   ccr.Token,
		}
	}
	return &cluster.Cluster{
		Id:          cr.ID,
		Name:        cr.Name,
		Description: cr.Description,
		Connection:  connection,
		CreatedAt:   timestamppb.New(cr.CreatedAt),
		UpdatedAt:   timestamppb.New(cr.UpdatedAt),
	}
}

// ResolveKubeConfig 提取kubeConfig中有效的上下文配置信息。
// 具体实现参考了clientcmd.BuildConfigFromFlags
// 由于是在server端解析，不处理文件类型的认证参数（没有data结尾的参数）
func (c *ClusterBiz) ResolveKubeConfig(_ context.Context, request *cluster.ResolveKubeConfigRequest) (*cluster.ResolveKubeConfigResponse, error) {
	config, err := clientcmd.Load([]byte(request.Content))
	if err != nil {
		log.Error().
			Err(err).
			Str("content", request.Content).
			Msg("failed to load kube config")
		return nil, err
	}
	for name, obj := range config.AuthInfos {
		config.AuthInfos[name] = obj
	}
	for name, obj := range config.Clusters {
		config.Clusters[name] = obj
	}
	res := make([]*cluster.ResolveKubeConfigResponse_Context, 0)
	for name, contextConfig := range config.Contexts {
		authConfig := config.AuthInfos[contextConfig.AuthInfo]
		if authConfig == nil ||
			authConfig.ClientCertificateData == nil ||
			authConfig.ClientKeyData == nil {
			continue
		}

		clusterConfig := config.Clusters[contextConfig.Cluster]
		if clusterConfig == nil ||
			clusterConfig.Server == "" ||
			clusterConfig.CertificateAuthorityData == nil {
			continue
		}
		server := clusterConfig.Server
		if u, err := url.ParseRequestURI(server); err == nil && u.Opaque == "" && len(u.Path) > 1 {
			u.RawQuery = ""
			u.Fragment = ""
			server = u.String()
		}

		current := &cluster.ResolveKubeConfigResponse_Context{
			Name:      name,
			Namespace: contextConfig.Namespace,
			Current:   name == config.CurrentContext,
			Cluster: &cluster.ResolveKubeConfigResponse_Cluster{
				Server: server,
				Ca:     string(clusterConfig.CertificateAuthorityData),
			},
			User: &cluster.ResolveKubeConfigResponse_User{
				Cert: string(authConfig.ClientCertificateData),
				Key:  string(authConfig.ClientKeyData),
			},
		}
		res = append(res, current)
	}
	return &cluster.ResolveKubeConfigResponse{Items: res}, nil
}

func (c *ClusterBiz) getKubeRestConfig(request *cluster.Connection) *rest.Config {
	config := &rest.Config{
		Host: request.Address,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(request.Ca),
		},
	}
	if request.Type == cluster.Connection_TLS_CERT {
		config.TLSClientConfig.CertData = []byte(request.Cert)
		config.TLSClientConfig.KeyData = []byte(request.Key)
	} else if request.Type == cluster.Connection_TLS_TOKEN {
		config.BearerToken = request.Token
	}
	return config
}

func (c *ClusterBiz) TestConnection(_ context.Context, request *cluster.Connection) (*cluster.TestConnectionResponse, error) {
	config := c.getKubeRestConfig(request)
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	if version, err := client.ServerVersion(); err != nil {
		return nil, err
	} else {
		return &cluster.TestConnectionResponse{Version: version.String()}, nil
	}
}

func (c *ClusterBiz) TestOperator(ctx context.Context, request *cluster.Connection) (*cluster.TestOperatorResponse, error) {
	config := c.getKubeRestConfig(request)
	client, err := apiextensionsclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	// 检查crd是否创建
	_, err = client.ApiextensionsV1().
		CustomResourceDefinitions().
		Get(ctx, "applications.app.kubeland", metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return &cluster.TestOperatorResponse{}, nil
}
