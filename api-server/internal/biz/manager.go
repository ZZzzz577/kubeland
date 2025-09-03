package biz

import (
	"api-server/api/v1/cluster"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"api-server/internal/data/generated/application"
	"api-server/internal/kube"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sync"
	"time"
)

var (
	ErrClientNotFound = errors.New("kubernetes client not found")
)

type ClusterManager struct {
	ctx    context.Context
	cancel context.CancelFunc

	db        *data.Data
	clusterID uint64
	updateAt  time.Time
	mgr       manager.Manager
}

func NewClusterManager(conn *generated.ClusterConnection) (*ClusterManager, error) {
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

	ctrl.SetLogger(zap.New())
	mgr, err := manager.New(config, manager.Options{
		Metrics: metricsserver.Options{
			BindAddress: "0",
		},
	})
	if err != nil {
		return nil, err
	}
	if err = kube.RegisterControllers(mgr); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &ClusterManager{
		ctx:       ctx,
		cancel:    cancel,
		clusterID: conn.ClusterID,
		updateAt:  conn.UpdatedAt,
		mgr:       mgr,
	}, nil
}

func (m *ClusterManager) Start() {
	log.Info().Uint64("clusterId", m.clusterID).Msg("start cluster manager")
	go func() {
		if err := m.mgr.Start(m.ctx); err != nil {
			log.Error().Err(err).Uint64("clusterId", m.clusterID).Msg("failed to start cluster manager")
		}
	}()
}

func (m *ClusterManager) Stop() {
	log.Info().Uint64("cluster_id", m.clusterID).Msg("stopped cluster manager")
	m.cancel()
}

type ClusterManagers struct {
	db *data.Data
	// 同步间隔
	refreshInterval time.Duration
	ticker          *time.Ticker
	mu              sync.RWMutex
	managers        map[uint64]*ClusterManager
}

func NewClusterManagers(db *data.Data) (*ClusterManagers, func(), error) {
	connections, err := db.ClusterConnection.Query().
		All(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("list clusters error")
		return nil, nil, err
	}
	managers := make(map[uint64]*ClusterManager, len(connections))

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
		for _, mgr := range cm.managers {
			mgr.Stop()
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
		// 新增/更新
		oldMgr, exists := c.managers[clusterID]
		if !exists || conn.UpdatedAt.After(oldMgr.updateAt) {
			// 停止旧manager
			if exists {
				oldMgr.Stop()
				delete(c.managers, clusterID)
			}

			// 创建新manager
			newMgr, err := NewClusterManager(conn)
			if err != nil {
				log.Error().
					Err(err).
					Str("address", conn.Address).
					Msg("get client error")
				continue
			}
			// 启动新manager
			newMgr.Start()

			c.managers[clusterID] = newMgr
		}
	}
	// 删除manager
	for clusterId, mgr := range c.managers {
		if _, ok := dbClusterIds[clusterId]; !ok {
			mgr.Stop()
			delete(c.managers, clusterId)
		}
	}
	return nil
}

func (c *ClusterManagers) GetClient(ctx context.Context, appName string) (client.Client, error) {
	app, err := c.db.Application.Query().
		Select(application.FieldClusterID).
		Where(
			application.Name(appName),
		).Only(ctx)
	if err != nil {
		log.Error().Err(err).Msg("get application error")
		return nil, err
	}
	clusterId := app.ClusterID
	c.mu.RLock()
	defer c.mu.RUnlock()
	if mgr := c.managers[clusterId]; mgr == nil || mgr.mgr == nil {
		return nil, ErrClientNotFound
	} else {
		return mgr.mgr.GetClient(), nil
	}
}
