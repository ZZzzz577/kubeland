package biz

import (
	"api-server/api/v1/cluster"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"api-server/internal/kube/controller"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sync"
	"time"
)

var (
	ErrClientNotFound = errors.New("kubernetes client not found")
)

type ManagerWrapper struct {
	cancel   context.CancelFunc
	updateAt time.Time
	Manager  manager.Manager
}

type ClusterManagers struct {
	db *data.Data
	// 同步间隔
	refreshInterval time.Duration
	ticker          *time.Ticker
	mu              sync.RWMutex
	managers        map[uint64]*ManagerWrapper
}

func NewClusterManagers(db *data.Data) (*ClusterManagers, func(), error) {
	connections, err := db.ClusterConnection.Query().
		All(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("list clusters error")
		return nil, nil, err
	}
	managers := make(map[uint64]*ManagerWrapper, len(connections))

	cm := &ClusterManagers{
		db:              db,
		refreshInterval: 30 * time.Second,
		managers:        managers,
	}
	// 启动时首次同步
	if err = cm.Refresh(); err != nil {
		log.Error().Err(err).Msg("refresh client error")
		return nil, nil, err
	}

	cm.ticker = time.NewTicker(cm.refreshInterval)
	go cm.RefreshLoop()

	cleanup := func() {
		cm.ticker.Stop()
		cm.mu.Lock()
		defer cm.mu.Unlock()
		for clusterID, mgr := range cm.managers {
			log.Info().Uint64("cluster_id", clusterID).Msg("shutting down cluster manager")
			mgr.cancel()
		}
		cm.managers = nil
	}
	return cm, cleanup, nil
}

func (c *ClusterManagers) RefreshLoop() {
	for range c.ticker.C {
		if err := c.Refresh(); err != nil {
			log.Error().Err(err).Msg("refresh client error")
		}
	}
}

func (c *ClusterManagers) Refresh() error {
	ctx := context.Background()
	connections, err := c.db.ClusterConnection.Query().
		All(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("list clusters error")
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	dbClusterIds := make(map[uint64]struct{}, len(connections))
	for _, conn := range connections {
		clusterID := conn.ClusterID
		dbClusterIds[clusterID] = struct{}{}
		// 新增或更新
		currentMgr, exists := c.managers[clusterID]
		if !exists || conn.UpdatedAt.After(currentMgr.updateAt) {
			// 创建新manager
			newMgr, err := c.CreateManagerByConfig(conn)
			if err != nil {
				log.Error().
					Err(err).
					Str("address", conn.Address).
					Msg("get client error")
				continue
			}
			// 停止旧manager
			if exists {
				log.Debug().Uint64("cluster_id", clusterID).Msg("stopping old cluster manager")
				currentMgr.cancel()
			}
			// 启动新manager
			ctx, cancel := context.WithCancel(ctx)
			newMgr.cancel = cancel
			go func(mgr manager.Manager, ctx context.Context, clusterID uint64) {
				if err := mgr.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
					log.Error().
						Err(err).
						Uint64("cluster_id", clusterID).
						Msg("cluster manager stopped with error")
				}
			}(newMgr.Manager, ctx, clusterID)
			c.managers[clusterID] = newMgr
			log.Info().Uint64("cluster_id", clusterID).Msg("cluster manager updated successfully")
		}
	}
	// 删除manager
	for clusterId, mgr := range c.managers {
		if _, ok := dbClusterIds[clusterId]; !ok {
			mgr.cancel()
			delete(c.managers, clusterId)
		}
	}
	return nil
}

func (c *ClusterManagers) CreateManagerByConfig(conn *generated.ClusterConnection) (*ManagerWrapper, error) {
	// 获取config
	config := &rest.Config{
		Host: conn.Address,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(conn.Ca),
		},
	}
	switch conn.Type {
	case uint8(cluster.Connection_TLS_CERT):
		config.CertData = []byte(conn.Cert)
		config.KeyData = []byte(conn.Key)
	case uint8(cluster.Connection_TLS_TOKEN):
		config.BearerToken = conn.Token
	}
	// 创建manager
	mgr, err := manager.New(config, manager.Options{
		Logger: zap.New(),
	})
	if err != nil {
		return nil, err
	}
	// 创建controllers
	if err = controller.NewControllers(conn.ClusterID, mgr, c.db); err != nil {
		return nil, err
	}
	return &ManagerWrapper{
		updateAt: conn.UpdatedAt,
		Manager:  mgr,
	}, nil
}

func (c *ClusterManagers) GetClient(clusterId uint64) (client.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if mgr := c.managers[clusterId]; mgr == nil || mgr.Manager == nil || mgr.Manager.GetClient() == nil {
		return nil, ErrClientNotFound
	} else {
		return mgr.Manager.GetClient(), nil
	}
}
