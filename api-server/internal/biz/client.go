package biz

import (
	"api-server/api/v1/cluster"
	"api-server/internal/data"
	"api-server/internal/data/generated"
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sync"
	"time"
)

var (
	ErrClientNotFound = errors.New("kubernetes client not found")
)

type ClientWrapper struct {
	updateAt time.Time
	*kubernetes.Clientset
	*dynamic.DynamicClient
}

type ClientManager struct {
	db *data.Data
	// 同步间隔
	refreshInterval time.Duration
	ticker          *time.Ticker
	mu              sync.RWMutex
	clients         map[uint64]*ClientWrapper
}

func NewClientManager(db *data.Data) (*ClientManager, func(), error) {
	connections, err := db.ClusterConnection.Query().
		All(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("list clusters error")
		return nil, nil, err
	}
	clients := make(map[uint64]*ClientWrapper, len(connections))

	cm := &ClientManager{
		db:              db,
		refreshInterval: 30 * time.Second,
		clients:         clients,
	}
	// 启动时首次同步
	if err = cm.Refresh(); err != nil {
		log.Error().Err(err).Msg("refresh client error")
		return nil, nil, err
	}

	cm.ticker = time.NewTicker(cm.refreshInterval)
	go cm.RefreshLoop()
	return cm, func() {
		cm.ticker.Stop()
	}, nil
}

func (c *ClientManager) RefreshLoop() {
	for range c.ticker.C {
		if err := c.Refresh(); err != nil {
			log.Error().Err(err).Msg("refresh client error")
		}
	}
}

func (c *ClientManager) Refresh() error {
	connections, err := c.db.ClusterConnection.Query().
		All(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("list clusters error")
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	// 新增或更新client
	dbConnIds := make(map[uint64]struct{}, len(connections))
	for _, conn := range connections {
		id := conn.ID
		dbConnIds[id] = struct{}{}
		if c.clients[id] == nil || conn.UpdatedAt.After(c.clients[id].updateAt) {
			client, err := c.GetClientByConfig(conn)
			if err != nil {
				log.Error().
					Err(err).
					Str("address", conn.Address).
					Msg("get client error")
				continue
			}
			c.clients[id] = client
		}
	}
	// 删除client
	for id := range c.clients {
		if _, ok := dbConnIds[id]; !ok {
			delete(c.clients, id)
		}
	}
	return nil
}

func (c *ClientManager) GetClientByConfig(conn *generated.ClusterConnection) (*ClientWrapper, error) {
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
	// 获取普通资源的client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	// 获取crd资源的client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &ClientWrapper{
		updateAt:      conn.UpdatedAt,
		Clientset:     client,
		DynamicClient: dynamicClient,
	}, nil
}

func (c *ClientManager) Get(clusterId uint64) (*ClientWrapper, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if client := c.clients[clusterId]; client == nil {
		return nil, ErrClientNotFound
	} else {
		return client, nil
	}
}
