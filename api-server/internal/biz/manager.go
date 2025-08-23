package biz

import (
	"api-server/api/v1/cluster"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"api-server/internal/kube/clientset/versioned"
	"api-server/internal/kube/informers/externalversions"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
	"sync"
	"time"
)

var (
	ErrClientNotFound = errors.New("kubernetes client not found")
)

const (
	defaultRsyncInterval = 10 * time.Hour
)

type ClusterManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	clusterID       uint64
	updateAt        time.Time
	clientset       *versioned.Clientset
	informerFactory externalversions.SharedInformerFactory

	applicationController *ApplicationController
}

func NewManagerWrapper(conn *generated.ClusterConnection) (*ClusterManager, error) {
	ctx, cancel := context.WithCancel(context.Background())

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
	clientSet, err := versioned.NewForConfig(config)
	if err != nil {
		cancel()
		return nil, err
	}

	informerFactory := externalversions.NewSharedInformerFactory(clientSet, defaultRsyncInterval)

	applicationController := NewApplicationController(conn.ClusterID, informerFactory)

	return &ClusterManager{
		ctx:    ctx,
		cancel: cancel,

		clusterID:       conn.ClusterID,
		updateAt:        conn.UpdatedAt,
		clientset:       clientSet,
		informerFactory: informerFactory,

		applicationController: applicationController,
	}, nil
}

func (m *ClusterManager) Start() {
	log.Info().Uint64("clusterId", m.clusterID).Msg("start cluster manager")
	m.informerFactory.Start(m.ctx.Done())
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		log.Info().Uint64("clusterId", m.clusterID).Msg("start application controller")
		if err := m.applicationController.Run(m.ctx, 2); err != nil {
			log.Error().
				Err(err).
				Uint64("cluster_id", m.clusterID).
				Msg("application controller start with error")
		}
	}()
}

func (m *ClusterManager) Stop() {
	m.cancel()
	m.wg.Wait()
	log.Info().Uint64("cluster_id", m.clusterID).Msg("cluster manager stopped successfully")
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
		// 新增或更新
		oldManager, exists := c.managers[clusterID]
		if !exists || conn.UpdatedAt.After(oldManager.updateAt) {
			// 创建新manager
			newManager, err := NewManagerWrapper(conn)
			if err != nil {
				log.Error().
					Err(err).
					Str("address", conn.Address).
					Msg("get client error")
				continue
			}
			// 停止旧manager
			if exists {
				oldManager.Stop()
				delete(c.managers, clusterID)
			}
			// 启动新manager
			newManager.Start()
			c.managers[clusterID] = newManager
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

func (c *ClusterManagers) GetClient(clusterId uint64) (*versioned.Clientset, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if mgr := c.managers[clusterId]; mgr == nil || mgr.clientset == nil {
		return nil, ErrClientNotFound
	} else {
		return mgr.clientset, nil
	}
}
